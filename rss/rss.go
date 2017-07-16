package rss

import (
	"encoding/xml"
	"github.com/gunnihinn/evil-rss-reader/reader"
	"html/template"
	"io/ioutil"
	"net/http"
)

func New(url string) *Feed {
	return &Feed{
		url: url,
	}
}

type Feed struct {
	url   string
	title string

	// Maybe add: type (rss/atom)

	// Generated at runtime
	entries []reader.Entry
	err     error
}

func (f Feed) Title() string {
	return f.title
}

func (f Feed) Url() string {
	return f.url
}

func (f Feed) Entries() []reader.Entry {
	return f.entries
}

func (f Feed) Error() error {
	return f.err
}

func (f *Feed) Update() {
	rf, err := fetchRssFeed(f.Url())

	if err == nil {
		if f.title == "" {
			f.title = rf.Title
		}

		f.entries = make([]reader.Entry, len(rf.Items))
		for i, item := range rf.Items {
			entry := Entry{
				title: item.Title,
				url:   item.Link,
			}
			if item.Description != "" {
				entry.content = item.Description
			} else if item.Content != "" {
				entry.content = item.Content
			}
			f.entries[i] = entry
		}
	}

	f.err = err
}

type Entry struct {
	title   string
	url     string        // optional
	content template.HTML // optional
}

func (e Entry) Title() string {
	return e.title
}

func (e Entry) Url() string {
	return e.url
}

func (e Entry) Content() template.HTML {
	return e.content
}

type RssFeed struct {
	XMLName xml.Name `xml:"rss"`

	// required
	Title       string `xml:"channel>title"`
	Link        string `xml:"channel>link"`
	Description string `xml:"channel>description"`

	// optional
	Items []RssItem `xml:"channel>item"`
}

type RssItem struct {
	/*
	 * "All elements of an item are optional, however at least one of title or
	 * description must be present."
	 * http://cyber.harvard.edu/rss/rss.html
	 */
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	Description template.HTML `xml:"description"` // lol security
	Content     template.HTML `xml:"encoded"`
}

func ParseRssFeed(feed []byte) (RssFeed, error) {
	f := RssFeed{}
	err := xml.Unmarshal(feed, &f)

	return f, err
}

func bytesFromUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	blob, err := ioutil.ReadAll(resp.Body)

	return blob, err
}

func fetchRssFeed(url string) (RssFeed, error) {
	blob, err := bytesFromUrl(url)
	if err != nil {
		return RssFeed{}, err
	}

	feed, err := ParseRssFeed(blob)
	if err != nil {
		return RssFeed{}, err
	}

	return feed, nil
}
