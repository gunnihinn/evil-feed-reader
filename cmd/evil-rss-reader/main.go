package main

import (
	"flag"
	"fmt"
	"github.com/gunnihinn/evil-rss-reader/reader"
	"github.com/gunnihinn/evil-rss-reader/rss"
	"html/template"
	"log"
	"net/http"
	"time"
)

var URLS = []string{
	"https://blog.regehr.org/feed",
	"http://ithare.com/feed/",
}

func main() {
	var port = flag.Int("port", 8080, "HTTP port")
	flag.Parse()

	feeds := make([]reader.Feed, len(URLS))
	for i, url := range URLS {
		feeds[i] = rss.New(url)
	}

	go func() {
		for {
			for _, feed := range feeds {
				go func(f reader.Feed) {
					f.Update()
					log.Printf("Feed '%s' has %d entries\n", f.Title(), len(f.Entries()))
				}(feed)
			}
			time.Sleep(15 * time.Minute)
		}
	}()

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
