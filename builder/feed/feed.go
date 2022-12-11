package newsfeed

import (
	"fmt"
	"io/ioutil"

	"github.com/anaskhan96/soup"
)

type Feed struct {
	HeaderTitle         string
	ArticlesSet         []string
	EntriesHTMLPath     string
	BaseEntriesHTMLPath string
	doc                 soup.Root
}

func (f *Feed) LoadHTML() error {
	data, err := ioutil.ReadFile(f.EntriesHTMLPath)
	if err != nil {
		return fmt.Errorf("LoadHTML: error", err)
	}
	f.doc = soup.HTMLParse(string(data))
	f.HeaderTitle = f.doc.Find("header").FullText()
	articles := f.doc.FindAll("article")
	for _, article := range articles {
		f.ArticlesSet = append(f.ArticlesSet, article.HTML())
	}
	if f.BaseEntriesHTMLPath != "" {
		data, err := ioutil.ReadFile(f.BaseEntriesHTMLPath)
		if err != nil {
			return fmt.Errorf("LoadHTML: error", err)
		}
		f.doc = soup.HTMLParse(string(data))
		f.HeaderTitle = f.doc.Find("header").FullText()
		articles := f.doc.FindAll("article")
		for _, article := range articles {
			f.ArticlesSet = append(f.ArticlesSet, article.HTML())
		}
	}
	return nil
}

func (f *Feed) Length() int {
	return len(f.ArticlesSet)
}

func (f *Feed) Article(index int) *Article {
	html := soup.HTMLParse(f.ArticlesSet[index])
	articleData := html.Find("article").Attrs()
	articleSummary := html.Find("details").Find("summary").FullText()
	return &Article{
		UID:           articleData["id"],
		Title:         articleData["title"],
		Link:          articleData["href"],
		Author:        articleData["author"],
		PublishedDate: articleData["published"],
		UpdatedDate:   articleData["updated"],
		Summary:       articleSummary,
		content:       html.HTML(),
	}
}

type Article struct {
	UID           string
	Title         string
	Link          string
	Author        string
	PublishedDate string
	UpdatedDate   string
	Summary       string
	// TODO: you have to collect this from the HTML itself and you have to take away the article and summary parts
	content string
}

func (a *Article) Content() string {
	str := ""
	doc := soup.HTMLParse(string(a.content))
	articleBody := doc.FindAll("")
	for _, v := range articleBody[5:] {
		str += v.HTML()
	}
	return str
}

func (a *Article) Entry() string {
	return fmt.Sprintf(
		"<entry>\n\t<id>%s</id>\n\t<title>%s</title>\n\t<updated>%s</updated>\n\t<author><name>%s</name></author>\n\t<link href=\"%s\" rel=\"alternate\"/>\n\t<published>%s</published>\n\t<summary>%s</summary>\n\t<content type=\"xhtml\">\n\t\t<div xmlns=\"http://www.w3.org/1999/xhtml\">\n\t\t%s\n\t\t</div>\n\t</content>\n</entry>",
		a.UID,
		a.Title,
		a.UpdatedDate,
		a.Author,
		a.Link,
		a.PublishedDate,
		a.Summary,
		a.Content(),
	)
}
