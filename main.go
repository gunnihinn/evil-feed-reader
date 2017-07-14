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

	Url   string // The URL we get the feed from; may or may not equal Link
	Error error

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

func FetchFeed(url string) Feed {
	blob, err := bytesFromUrl(url)
	if err != nil {
		return Feed{
			Error: err,
		}
	}

	feed, err := ParseFeed(blob)
	if err != nil {
		return Feed{
			Error: err,
		}
	}
	feed.Url = url

	return feed
}

func main() {
	var port = flag.Int("port", 8080, "HTTP port")
	flag.Parse()

	urls := []string{
		"https://blog.regehr.org/feed",
		"http://ithare.com/feed/",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("index.html")
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}

		feeds := make([]Feed, 0)
		for _, url := range urls {
			feeds = append(feeds, FetchFeed(url))
		}

		t.Execute(w, feeds)
	})
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
