package newsstats

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type NewsStats struct {
	DownloadLangs map[string]int
	StateFile     string
}

func (n *NewsStats) Increment(rq http.Request) {
	q := rq.URL.Query()
	lang := q.Get("lang")
	if lang != "" {
		n.DownloadLangs[lang]++
	} else {
		n.DownloadLangs["en_US"]++
	}
}

func (n *NewsStats) Save() error {
	bytes, err := json.Marshal(n.DownloadLangs)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(n.StateFile, bytes, 0644); err != nil {
		return err
	}
	return nil
}

func (n *NewsStats) Load() {
	bytes, err := ioutil.ReadFile(n.StateFile)
	if err != nil {
		n.DownloadLangs = make(map[string]int)
	}
	if err := json.Unmarshal(bytes, &n.DownloadLangs); err != nil {
		n.DownloadLangs = make(map[string]int)
	}
}
