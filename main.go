package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var URLS = []string{
	"https://blog.regehr.org/feed",
	"http://ithare.com/feed/",
}

type Feed struct {
	Url   string
	Title string

	// Maybe add: type (rss/atom)

	// Generated at runtime
	Entries []Entry
	Error   error
}

func (f *Feed) Update() {
	rf, err := FetchRssFeed(f.Url)

	if err == nil {
		if f.Title == "" {
			f.Title = rf.Title
		}

		f.Entries = make([]Entry, len(rf.Items))
		for i, item := range rf.Items {
			entry := Entry{
				Title: item.Title,
				Url:   item.Link,
			}
			if item.Description != "" {
				entry.Content = item.Description
			} else if item.Content != "" {
				entry.Content = item.Content
			}
			f.Entries[i] = entry
		}
	}

	f.Error = err
}

type Entry struct {
	Title   string
	Url     string        // optional
	Content template.HTML // optional
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

func FetchRssFeed(url string) (RssFeed, error) {
	log.Printf("Fetching RSS feed from '%s'\n", url)

	blob, err := bytesFromUrl(url)
	if err != nil {
		log.Printf("%s\n", err)
		return RssFeed{}, err
	}

	feed, err := ParseRssFeed(blob)
	if err != nil {
		log.Printf("%s\n", err)
		return RssFeed{}, err
	}

	return feed, nil
}

func main() {
	var port = flag.Int("port", 8080, "HTTP port")
	flag.Parse()

	feeds := make([]*Feed, len(URLS))
	for i, url := range URLS {
		feeds[i] = &Feed{
			Url: url,
		}
	}

	for _, feed := range feeds {
		feed.Update()
		log.Printf("Feed '%s' has %d entries\n", feed.Title, len(feed.Entries))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("index.html")
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
		t.Execute(w, feeds)
	})
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
