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

	// required
	Title       string `xml:"channel>title"`
	Url         string `xml:"channel>link"`
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
	Url         string        `xml:"link"`
	Content     template.HTML `xml:"encoded"`
}

func ParseFeed(feed []byte) (Feed, error) {
	f := Feed{}
	err := xml.Unmarshal(feed, &f)

	return f, err
}

func main() {
	var port = flag.Int("port", 8080, "HTTP port")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("index.html")
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}

		test, err := ioutil.ReadFile("example-feed.rss")
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}

		feed, err := ParseFeed(test)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}

		t.Execute(w, feed)
	})
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
