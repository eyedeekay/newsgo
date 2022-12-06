package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	server "i2pgit.org/idk/newsgo/server"
)

var (
	serve     = flag.String("command", "help", "command to run(may be `serve`,`build`,`sign`")
	dir       = flag.String("newsdir", "build", "directory to serve news from")
	statsfile = flag.String("statsfile", "build/stats.json", "file to store stats in")
	host      = flag.String("host", "127.0.0.1", "host to serve on")
	port      = flag.String("port", "9696", "port to serve on")
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
	fmt.Println("I2P News Server Tool/Library")
	fmt.Println("")
	fmt.Println("Usage")
	fmt.Println("-----")
	fmt.Println("")
	fmt.Println("./newsgo -command $command -newsdir $news_directory -statsfile $news_stats_file")
	fmt.Println("")
	fmt.Println("### Commands")
	fmt.Println("")
	fmt.Println(" - serve: Serve newsfeeds from a directory")
	fmt.Println(" - build: Build newsfeeds from XML(Not Implemented Yet)")
	fmt.Println(" - sign: Sign newsfeeds with local keys(Not Implemented Yet)")
	fmt.Println("")
	fmt.Println("### Options")
	fmt.Println("")
	fmt.Println("Use these options to configure the software")
	fmt.Println("")
	fmt.Println("#### Server Options(use with `serve`")
	fmt.Println("")
	fmt.Println(" - `-newsdir`: directory to serve newsfeed from")
	fmt.Println(" - `-statsfile`: file to store the stats in, in json format")
	fmt.Println(" - `-host`: host to serve news files on")
	fmt.Println(" - `-port`: port to serve news files on")
	fmt.Println(" - `-i2p`: serve news files directly to I2P using SAMv3")
	fmt.Println("")
	fmt.Println("#### Builder Options(use with `build`")
	fmt.Println("")
	fmt.Println("Not implemented yet")
	fmt.Println("")
	fmt.Println("#### Signer Options(use with `sign`")
	fmt.Println("")
	fmt.Println("Not implemented yet")
}

func Serve(host, port string) error {
	ln, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return err
	}
	s := server.Serve(*dir, *statsfile)
	return http.Serve(ln, s)
}

func main() {
	flag.Parse()
	command := validateCommand(serve)
	validatePort(port)
	switch command {
	case "serve":
		if err := Serve(*host, *port); err != nil {
			panic(err)
		}
	case "build":
	case "sign":
	case "help":
		Help()
	default:
		Help()
	}
}
