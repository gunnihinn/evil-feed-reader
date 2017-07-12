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
	Title   string   `xml:"channel>title"`
	Items   []Item   `xml:"channel>item"`
}

type Item struct {
	Title       string `xml:"title"`
	Url         string `xml:"link"`
	Description string `xml:"description"`
	Content     string `xml:"content"`
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
