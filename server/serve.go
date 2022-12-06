package newsserver

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/golang-commonmark/markdown"
	stats "i2pgit.org/idk/newsgo/server/stats"
)

type NewsServer struct {
	NewsDir string
	Stats   stats.NewsStats
}

var serveTest http.Handler = &NewsServer{}

func (n *NewsServer) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	path := rq.URL.Path
	file := filepath.Join(n.NewsDir, path)
	if err := fileCheck(file); err != nil {
		log.Println("ServeHTTP:", err.Error())
		rw.WriteHeader(404)
		return
	}
	if err := ServeFile(file, rw); err != nil {
		log.Println("ServeHTTP:", err.Error())
		rw.WriteHeader(404)
	}
}

func fileCheck(file string) error {
	if _, err := os.Stat(file); err != nil {
		return fmt.Errorf("fileCheck: %s", err)
	}
	return nil
}

func fileType(file string) (string, error) {
	base := filepath.Base(file)
	if base == "" {
		return "", fmt.Errorf("fileType: Invalid file path passed to type determinator")
	}
	extension := filepath.Ext(base)
	switch extension {
	case ".su3":
		return "application/x-i2p-su3-news", nil
	case ".html":
		return "text/html", nil
	case ".xml":
		return "application/rss+xml", nil
	case ".atom.xml":
		return "application/rss+xml", nil
	case ".svg":
		return "image/svg+xml", nil
	default:
		return "text/html", nil
	}
}

func openDirectory(wd string) string {
	//wd = strings.Replace(wd, leader, "", 1)
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Fatal(err)
	}
	var readme string
	log.Println("Navigating directory:", wd)
	nwd := strings.Join(strings.Split(wd, "/")[1:], "/")
	readme += fmt.Sprintf("%s\n", filepath.Base(wd))
	readme += fmt.Sprintf("%s\n", head(len(filepath.Base(wd))))
	readme += fmt.Sprintf("%s\n", "")
	readme += fmt.Sprintf("%s\n", "**Directory Listing:**")
	readme += fmt.Sprintf("%s\n", "")
	for _, file := range files {
		if !file.IsDir() {
			fmt.Println(file.Name(), file.IsDir())
			xname := filepath.Join(wd, file.Name())
			bytes, err := ioutil.ReadFile(xname)
			if err != nil {
				log.Println("Listing error:", err)
			}
			sum := fmt.Sprintf("%x", sha256.Sum256(bytes))
			readme += fmt.Sprintf(" - [%s](%s/%s) : `%d` : `%s` - `%s`\n", file.Name(), filepath.Base(nwd), filepath.Base(file.Name()), file.Size(), file.Mode(), sum)
		} else {
			fmt.Println(file.Name(), file.IsDir())
			readme += fmt.Sprintf(" - [%s](%s/) : `%d` : `%s`\n", file.Name(), file.Name(), file.Size(), file.Mode())
		}
	}
	return readme
}

func hTML(mdtxt string) []byte {
	md := markdown.New(markdown.XHTMLOutput(true))
	return []byte(md.RenderToString([]byte(mdtxt)))
}

func head(num int) string {
	var r string
	for i := 0; i < num; i++ {
		r += "="
	}
	return r
}

func ServeFile(file string, rw http.ResponseWriter) error {
	//if err := fileCheck(file); err != nil {
	//	return fmt.Errorf("ServeFile: %s", err)
	//}
	ftype, err := fileType(file)
	if err != nil {
		return fmt.Errorf("ServeFile: %s", err)
	}
	if ftype == "application/x-i2p-su3-news" {
		// Log stats here
	}
	rw.Header().Add("Content-Type", ftype)
	f, _ := os.Stat(file)
	if f.IsDir() {
		bytes := hTML(openDirectory(file))
		rw.Write(bytes)
		return nil
	}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("ServeFile: %s", err)
	}

	log.Println("ServeFile: ", file, ftype)
	rw.Write(bytes)
	return nil
}

func Serve(newsDir, newsStats string) *NewsServer {
	s := &NewsServer{
		NewsDir: newsDir,
		Stats: stats.NewsStats{
			StateFile: newsStats,
		},
	}
	return s
}
