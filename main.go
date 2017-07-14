package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Feed struct {
	XMLName xml.Name `xml:"rss"`

	Url string // The URL we get the feed from; may or may not equal Link

	// required
	Title       string `xml:"channel>title"`
	Link        string `xml:"channel>link"`
	Description string `xml:"channel>description"`

	// optional
	Items []Item `xml:"channel>item"`
}

type Item struct {
	/*
	 * "All elements of an item are optional, however at least one of title or
	 * description must be present."
	 * http://cyber.harvard.edu/rss/rss.html
	 */
	Title       string        `xml:"title"`
	Description template.HTML `xml:"description"` // lol security
	Link        string        `xml:"link"`
	Content     template.HTML `xml:"encoded"`
}

func ParseFeed(feed []byte) (Feed, error) {
	f := Feed{}
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

func FetchFeed(url string) (Feed, error) {
	blob, err := bytesFromUrl(url)
	if err != nil {
		return Feed{}, err
	}

	feed, err := ParseFeed(blob)
	if err != nil {
		return feed, err
	}
	feed.Url = url

	return feed, nil
}

func main() {
	var port = flag.Int("port", 8080, "HTTP port")
	flag.Parse()

	url := "https://blog.regehr.org/feed"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("index.html")
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}

		feed, err := FetchFeed(url)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}

		t.Execute(w, feed)
	})
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
