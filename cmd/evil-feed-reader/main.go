package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gunnihinn/evil-feed-reader/provider"
	"github.com/gunnihinn/evil-feed-reader/reader"
)

type Context struct {
	Sidebar []Navitem
	Feed    reader.Feed
	CSS     template.CSS
}

type Navitem struct {
	Title          string
	Resource       string
	HasRecentItems bool
}

func Prepare(feeds []reader.Feed, active reader.Feed) Context {
	ctx := Context{
		Feed:    active,
		Sidebar: make([]Navitem, len(feeds)),
	}

	data, err := Asset("cmd/evil-feed-reader/data/style.css")
	if err == nil {
		ctx.CSS = template.CSS(data)
	}

	for i, f := range feeds {
		ctx.Sidebar[i] = Navitem{
			Title:          f.Title(),
			Resource:       fmt.Sprintf("%d", i),
			HasRecentItems: f.HasRecentItems(),
		}
	}
	return ctx
}

func createHandler(feeds []reader.Feed, active reader.Feed) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := Asset("cmd/evil-feed-reader/data/index.html")
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}

		t := template.New("data/index.html")
		_, err = t.Parse(string(data))
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
		t.Execute(w, Prepare(feeds, active))
	}
}

func main() {
	var port = flag.Int("port", 8080, "HTTP port")
	var configFile = flag.String("feeds", "feeds.cfg", "Feeds config file")
	flag.Parse()

	urls, err := parseConfig(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	provider := provider.HTTP()
	feeds := make([]reader.Feed, len(urls))
	for i, url := range urls {
		feeds[i] = reader.New(provider, url)
	}

	go func() {
		for {
			for _, feed := range feeds {
				go func(f reader.Feed) {
					f.Update()

					if f.Error() != nil {
						log.Printf("Problems parsing feed '%s':\n", f.Resource())
						log.Printf("%s\n", f.Error())
					}

					if len(f.Entries()) != 0 {
						log.Printf("Feed '%s' has %d entries\n", f.Title(), len(f.Entries()))
					} else {
						log.Printf("Got no entries from '%s':\n", f.Resource())
					}
				}(feed)
			}
			time.Sleep(15 * time.Minute)
		}
	}()

	for i, feed := range feeds {
		http.HandleFunc(fmt.Sprintf("/%d", i), createHandler(feeds, feed))
	}
	http.HandleFunc("/", createHandler(feeds, nil))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
