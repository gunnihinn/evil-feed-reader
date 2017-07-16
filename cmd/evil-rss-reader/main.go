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

type Context struct {
	Sidebar []Navitem
	Feed    reader.Feed
}

type Navitem struct {
	Title    string
	Resource string
}

func Prepare(feeds []reader.Feed, active reader.Feed) Context {
	ctx := Context{
		Feed:    active,
		Sidebar: make([]Navitem, len(feeds)),
	}
	for i, f := range feeds {
		ctx.Sidebar[i] = Navitem{
			Title:    f.Title(),
			Resource: fmt.Sprintf("%d", i),
		}
	}
	return ctx
}

func createHandler(feeds []reader.Feed, active reader.Feed) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("index.html")
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
		t.Execute(w, Prepare(feeds, active))
	}
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

	for i, feed := range feeds {
		http.HandleFunc(fmt.Sprintf("/%d", i), createHandler(feeds, feed))
	}
	http.HandleFunc("/", createHandler(feeds, nil))

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
