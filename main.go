package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/eyedeekay/onramp"
	"github.com/google/uuid"
	builder "i2pgit.org/idk/newsgo/builder"
	server "i2pgit.org/idk/newsgo/server"
	signer "i2pgit.org/idk/newsgo/signer"
)

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	privPem, err := ioutil.ReadFile(path)
	if nil != err {
		return nil, err
	}

	privDer, _ := pem.Decode(privPem)
	privKey, err := x509.ParsePKCS1PrivateKey(privDer.Bytes)
	if nil != err {
		return nil, err
	}

	return privKey, nil
}

func Sign(xmlfeed string) error {
	sk, err := loadPrivateKey(*signingkey)
	if err != nil {
		return err
	}
	signer := signer.NewsSigner{
		SignerID:   *signerid,
		SigningKey: sk,
	}
	return signer.CreateSu3(xmlfeed)
}

var (
	serve     = flag.String("command", "help", "command to run(may be `serve`,`build`,`sign`, or `help`(default)")
	dir       = flag.String("newsdir", "build", "directory to serve news from")
	statsfile = flag.String("statsfile", "build/stats.json", "file to store stats in")
	host      = flag.String("host", "127.0.0.1", "host to serve on")
	port      = flag.String("port", "9696", "port to serve on")
	i2p       = flag.Bool("i2p", isSamAround(), "automatically co-host on an I2P service using SAMv3")
	tcp       = flag.Bool("http", true, "host on an HTTP service at host:port")
	//newsfile    = flag.String("newsfile", "data/entries.html", "entries to pass to news generator. If passed a directory, all `entries.html` files in the directory will be processed")
	newsfile    = flag.String("newsfile", "data", "entries to pass to news generator. If passed a directory, all `entries.html` files in the directory will be processed")
	bloclist    = flag.String("blocklist", "data/blocklist.xml", "block list file to pass to news generator")
	releasejson = flag.String("releasejson", "data/releases.json", "json file describing an update to pass to news generator")
	title       = flag.String("feedtitle", "I2P News", "title to use for the RSS feed to pass to news generator")
	subtitle    = flag.String("feedsubtitle", "News feed, and router updates", "subtitle to use for the RSS feed to pass to news generator")
	site        = flag.String("feedsite", "http://i2p-projekt.i2p", "site for the RSS feed to pass to news generator")
	mainurl     = flag.String("feedmain", DefaultFeedURL(), "Primary newsfeed for updates to pass to news generator")
	backupurl   = flag.String("feedbackup", "http://dn3tvalnjz432qkqsvpfdqrwpqkw3ye4n4i2uyfr4jexvo3sp5ka.b32.i2p/news/news.atom.xml", "Backup newsfeed for updates to pass to news generator")
	urn         = flag.String("feeduid", uuid.New().String(), "UUID to use for the RSS feed to pass to news generator")
	builddir    = flag.String("builddir", "build", "Build directory to output feeds to.")
	signerid    = flag.String("signerid", "null@example.i2p", "ID to use when signing the news")
	signingkey  = flag.String("signingkey", "signing_key.pem", "Path to a signing key")
)

func validatePort(s *string) {
	_, err := strconv.Atoi(*s)
	if err != nil {
		log.Println("Port is invalid")
		os.Exit(1)
	}
}

func validateCommand(s *string) string {
	switch *s {
	case "serve":
		return "serve"
	case "build":
		return "build"
	case "sign":
		return "sign"
	default:
		return "help"
	}
}

func Help() {
	fmt.Println("newsgo")
	fmt.Println("======")
	fmt.Println("")
	fmt.Println("I2P News Server Tool/Library. A whole lot faster than the python one. Otherwise compatible.")
	fmt.Println("")
	fmt.Println("Usage")
	fmt.Println("-----")
	fmt.Println("")
	fmt.Println("./newsgo -command $command -newsdir $news_directory -statsfile $news_stats_file")
	fmt.Println("")
	fmt.Println("### Commands")
	fmt.Println("")
	fmt.Println(" - serve: Serve newsfeeds from a directory")
	fmt.Println(" - build: Build newsfeeds from XML")
	fmt.Println(" - sign: Sign newsfeeds with local keys")
	fmt.Println("")
	fmt.Println("### Options")
	fmt.Println("")
	fmt.Println("Use these options to configure the software")
	fmt.Println("")
	fmt.Println("#### Server Options(use with `serve`)")
	fmt.Println("")
	fmt.Println(" - `-newsdir`: directory to serve newsfeed from")
	fmt.Println(" - `-statsfile`: file to store the stats in, in json format")
	fmt.Println(" - `-host`: host to serve news files on")
	fmt.Println(" - `-port`: port to serve news files on")
	fmt.Println(" - `-http`: serve news on host:port using HTTP")
	fmt.Println(" - `-i2p`: serve news files directly to I2P using SAMv3")
	fmt.Println("")
	fmt.Println("#### Builder Options(use with `build`)")
	fmt.Println("")
	fmt.Println(" - `-newsfile`: entries to pass to news generator. If passed a directory, all `entries.html` files in the directory will be processed")
	fmt.Println(" - `-blockfile`: block list file to pass to news generator")
	fmt.Println(" - `-releasejson`: json file describing an update to pass to news generator")
	fmt.Println(" - `-feedtitle`: title to use for the RSS feed to pass to news generator")
	fmt.Println(" - `-feedsubtitle`: subtitle to use for the RSS feed to pass to news generator")
	fmt.Println(" - `-feedsite`: site for the RSS feed to pass to news generator")
	fmt.Println(" - `-feedmain`: Primary newsfeed for updates to pass to news generator")
	fmt.Println(" - `-feedbackup`: Backup newsfeed for updates to pass to news generator")
	fmt.Println(" - `-feeduri`: UUID to use for the RSS feed to pass to news generator")
	fmt.Println(" - `-builddir`: directory to output XML files in")
	fmt.Println("")
	fmt.Println("#### Signer Options(use with `sign`)")
	fmt.Println("")
	fmt.Println(" - `-signerid`: ID of the news signer")
	fmt.Println(" - `-signingkey`: path to the signing key")
}

func Serve(host, port string, s *server.NewsServer) error {
	ln, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return err
	}
	return http.Serve(ln, s)
}

func ServeI2P(s *server.NewsServer) error {
	garlic := &onramp.Garlic{}
	defer garlic.Close()
	ln, err := garlic.Listen()
	if err != nil {
		return err
	}
	defer ln.Close()
	return http.Serve(ln, s)
}

func isSamAround() bool {
	ln, err := net.Listen("tcp", "127.0.0.1:7656")
	if err != nil {
		return true
	}
	ln.Close()
	return false
}

func DefaultFeedURL() string {
	if !isSamAround() {
		return "http://tc73n4kivdroccekirco7rhgxdg5f3cjvbaapabupeyzrqwv5guq.b32.i2p/news.atom.xml"
	}
	garlic := &onramp.Garlic{}
	defer garlic.Close()
	ln, err := garlic.Listen()
	if err != nil {
		return "http://tc73n4kivdroccekirco7rhgxdg5f3cjvbaapabupeyzrqwv5guq.b32.i2p/news.atom.xml"
	}
	defer ln.Close()
	return "http://" + ln.Addr().String() + "/news.atom.xml"
}

func Build(newsFile string) {
	news := builder.Builder(newsFile, *releasejson, *bloclist)
	news.TITLE = *title
	news.SITEURL = *site
	news.MAINFEED = *mainurl
	news.BACKUPFEED = *backupurl
	news.SUBTITLE = *subtitle
	news.URNID = *urn
	base := filepath.Join(*newsfile, "entries.html")
	if newsFile != base {
		news.Feed.BaseEntriesHTMLPath = base
	}
	if feed, err := news.Build(); err != nil {
		log.Printf("Build error: %s", err)
	} else {
		filename := strings.Replace(strings.Replace(strings.Replace(strings.Replace(newsFile, ".html", ".atom.xml", -1), "entries.", "news_", -1), "translations", "", -1), "news_atom", "news.atom", -1)
		if err := os.MkdirAll(filepath.Join(*builddir, filepath.Dir(filename)), 0755); err != nil {
			panic(err)
		}
		if err = ioutil.WriteFile(filepath.Join(*builddir, filename), []byte(feed), 0644); err != nil {
			panic(err)
		}
	}
}

func main() {
	flag.Parse()
	command := validateCommand(serve)
	validatePort(port)
	switch command {
	case "serve":
		s := server.Serve(*dir, *statsfile)
		if *tcp {
			go func() {
				if err := Serve(*host, *port, s); err != nil {
					panic(err)
				}
			}()
		}
		if *i2p {
			go func() {
				if err := ServeI2P(s); err != nil {
					panic(err)
				}
			}()
		}
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for sig := range c {
				log.Println("captured: ", sig)
				s.Stats.Save()
				os.Exit(0)
			}
		}()
		i := 0
		for {
			time.Sleep(time.Minute)
			log.Printf("Running for %d minutes.", i)
			i++
		}
	case "build":
		f, e := os.Stat(*newsfile)
		if e != nil {
			panic(e)
		}
		if f.IsDir() {
			err := filepath.Walk(*newsfile,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					ext := filepath.Ext(path)
					if ext == ".html" {
						Build(path)
					}
					return nil
				})
			if err != nil {
				log.Println(err)
			}
		} else {
			Build(*newsfile)
		}
	case "sign":
		f, e := os.Stat(*newsfile)
		if e != nil {
			panic(e)
		}
		if f.IsDir() {
			err := filepath.Walk(*newsfile,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					ext := filepath.Ext(path)
					if ext == ".html" {
						Sign(path)
					}
					return nil
				})
			if err != nil {
				log.Println(err)
			}
		} else {
			Sign(*newsfile)
		}
	case "help":
		Help()
	default:
		Help()
	}
}
