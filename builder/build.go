package newsbuilder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/google/uuid"
	"github.com/yosssi/gohtml"
	newsfeed "i2pgit.org/idk/newsgo/builder/feed"
)

type NewsBuilder struct {
	Feed         newsfeed.Feed
	ReleasesJson string
	BlocklistXML string
	URNID        string
	TITLE        string
	SITEURL      string
	MAINFEED     string
	BACKUPFEED   string
	SUBTITLE     string
}

func (nb *NewsBuilder) JSONtoXML() (string, error) {
	content, err := ioutil.ReadFile(nb.ReleasesJson)
	if err != nil {
		return "", err
	}
	// Now let's unmarshall the data into `payload`
	var payload []map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return "", err
	}
	str := ""
	/*
			<i2p:release date="2022-11-21" minVersion="0.9.9" minJavaVersion="1.8">
		    <i2p:version>2.0.0</i2p:version>
		    <i2p:update type="su3">
		      <i2p:torrent href="magnet:?xt=urn:btih:a50f8479a39896f00431d7b500447fe303d2b6b5&amp;dn=i2pupdate-2.0.0.su3&amp;tr=http://tracker2.postman.i2p/announce.php"/>
		      <i2p:url href="http://stats.i2p/i2p/2.0.0/i2pupdate.su3"/>
		      <i2p:url href="http://mgp6yzdxeoqds3wucnbhfrdgpjjyqbiqjdwcfezpul3or7bzm4ga.b32.i2p/releases/2.0.0/i2pupdate.su3"/>
		    </i2p:update>
		  </i2p:release>
	*/
	releasedate := payload[0]["date"]
	version := payload[0]["version"]
	minVersion := payload[0]["minVersion"]
	minJavaVersion := payload[0]["minJavaVersion"]
	updates := payload[0]["updates"].(map[string]interface{})["su3"].(map[string]interface{})
	magnet := updates["torrent"]
	urls := updates["url"].([]interface{})
	str += "<i2p:release date=" + releasedate.(string) + " minVersion=" + minVersion.(string) + " minJavaVersion=" + minJavaVersion.(string) + ">\n"
	str += "<i2p:version>" + version.(string) + "</i2p:version>"
	str += "<i2p:update type=\"su3\">"
	str += "<i2p:torrent href=\"" + magnet.(string) + "\"/>"
	for _, u := range urls {
		str += "<i2p:url href=\"" + u.(string) + "\"/>"
	}
	str += "</i2p:update>"
	str += "</i2p:release>"
	return str, nil
}

func (nb *NewsBuilder) Build() (string, error) {
	if err := nb.Feed.LoadHTML(); err != nil {
		return "", fmt.Errorf("Build: error %s", err.Error())
	}
	current_time := time.Now()
	str := "<?xml version='1.0' encoding='UTF-8'?>"
	str += "<feed xmlns:i2p=\"http://geti2p.net/en/docs/spec/updates\" xmlns=\"http://www.w3.org/2005/Atom\" xml:lang=\"en\">"
	str += "<id>" + "urn:uuid:" + nb.URNID + "</id>"
	str += "<title>" + nb.TITLE + "</title>"
	milli := current_time.Nanosecond() / 1000
	t := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.%02d+00:00\n", current_time.Year(), current_time.Month(), current_time.Day(), current_time.Hour(), current_time.Minute(), current_time.Second(), milli)
	str += "<updated>" + t + "</updated>"
	str += "<link href=\"" + nb.SITEURL + "\"/>"
	str += "<link href=\"" + nb.MAINFEED + "\" rel=\"self\"/>"
	if nb.BACKUPFEED != "" {
		str += "<link href=\"" + nb.BACKUPFEED + "\" rel=\"alternate\"/>"
	}
	str += "<generator uri=\"http://idk.i2p/newsgo\" version=\"0.1.0\">newsgo</generator>"
	str += "<subtitle>" + nb.SUBTITLE + "</subtitle>"
	blocklistBytes, err := ioutil.ReadFile(nb.BlocklistXML)
	if err != nil {
		return "", err
	}
	str += string(blocklistBytes)
	jsonxml, err := nb.JSONtoXML()
	if err != nil {
		return "", err
	}
	str += jsonxml
	for index, _ := range nb.Feed.ArticlesSet {
		art := nb.Feed.Article(index)
		str += art.Entry()
	}
	str += "</feed>"
	return gohtml.Format(str), nil
}

func Builder(newsFile, releasesJson, blocklistXML string) *NewsBuilder {
	nb := &NewsBuilder{
		Feed: newsfeed.Feed{
			EntriesHTMLPath: newsFile,
		},
		ReleasesJson: releasesJson,
		BlocklistXML: blocklistXML,
		URNID:        uuid.New().String(),
		TITLE:        "I2P News",
		SITEURL:      "http://i2p-projekt.i2p",
		MAINFEED:     "http://tc73n4kivdroccekirco7rhgxdg5f3cjvbaapabupeyzrqwv5guq.b32.i2p/news.atom.xml",
		BACKUPFEED:   "http://dn3tvalnjz432qkqsvpfdqrwpqkw3ye4n4i2uyfr4jexvo3sp5ka.b32.i2p/news/news.atom.xml",
		SUBTITLE:     "News feed, and router updates",
	}
	return nb
}
