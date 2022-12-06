package newsserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	stats "i2pgit.org/idk/newsgo/server/stats"
)

type NewsServer struct {
	NewsDir string
	stats.Stats   NewsStats
}

var serveTest http.Handler = &NewsServer{}

func (n *NewsServer) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	path := rq.URL.Path
	file := filepath.Join(n.NewsDir, path)
	if err := fileCheck(file); err != nil {
		rw.WriteHeader(404)
		return
	}
	if err := ServeFile(file, rw); err != nil {
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
		return "", fmt.Errorf("fileType: Attempted to serve invalid file type")
	}
}

func ServeFile(file string, rw http.ResponseWriter) error {
	if err := fileCheck(file); err != nil {
		return fmt.Errorf("ServeFile: %s", err)
	}
	ftype, err := fileType(file)
	if err != nil {
		return fmt.Errorf("ServeFile: %s", err)
	}
	if ftype == "application/x-i2p-su3-news" {

	}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("ServeFile: %s", err)
	}
	rw.Header().Add("Content-Type", ftype)
	rw.Write(bytes)
	return nil
}
