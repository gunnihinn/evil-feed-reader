package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
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
			HasRecentItems: !f.Seen(),
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

		if active != nil {
			active.SetSeen(true)
		}

		if err := t.Execute(w, Prepare(feeds, active)); err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
	}
}

func main() {
	var port = flag.Int("port", 8080, "HTTP port")
	var configFile = flag.String("feeds", "feeds.cfg", "Feeds config file")
	var stateFile = flag.String("state", ".evil-state.json", "Internal state file")
	flag.Parse()

	urls, err := parseConfig(*configFile)
	if err != nil {
		log.Printf("Couldn't parse config file: %s\n", err)
	}

	provider := provider.HTTP()
	feeds := make([]reader.Feed, len(urls))
	for i, url := range urls {
		feeds[i] = reader.New(provider, url)
	}

	state, err := parseState(*stateFile)
	if err != nil {
		log.Printf("Couldn't parse state file: %s\n", err)
	} else {
		for _, feed := range feeds {
			s, ok := state[feed.Resource()]
			if ok {
				feed.SetState(s)
			}
		}
	}

	go func() {
		for {
			for _, feed := range feeds {
				go func(f reader.Feed) {
					if err := f.Update(); err != nil {
						log.Printf("Problems parsing feed '%s':\n", f.Resource())
						log.Printf("%s\n", err)
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

	// TODO: Replace with signal handler
	go func(feeds []reader.Feed) {
		for {
			time.Sleep(1 * time.Minute) // LOL scheduling

			log.Printf("Writing state to disk\n")

			blob, err := reader.Marshal(feeds)
			if err != nil {
				log.Printf("Error creating state JSON: %s\n", err)
			} else {
				err = ioutil.WriteFile(*stateFile, blob, 0644)
				if err != nil {
					log.Printf("Error writing state JSON to disk: %s\n", err)
				}
			}

			time.Sleep(15 * time.Minute)
		}
	}(feeds)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Listening on %s\n", addr)
	handler := NewHandler()
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	for i, feed := range feeds {
		handler.mux.HandleFunc(fmt.Sprintf("/%d", i), createHandler(feeds, feed))
	}
	handler.mux.HandleFunc("/", createHandler(feeds, nil))

	log.Fatal(server.ListenAndServe())
}
