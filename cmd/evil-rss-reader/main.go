package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gunnihinn/evil-rss-reader/reader"
	"github.com/gunnihinn/evil-rss-reader/rss"
)

func parseConfig(filename string) ([]string, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return []string{}, err
	}
	defer fh.Close()

	urls := make([]string, 0)
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		if line == "" || strings.Index(line, "#") == 0 {
			continue
		}
		urls = append(urls, line)
	}

	return urls, scanner.Err()
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
	var configFile = flag.String("feeds", "feeds.cfg", "Feeds config  file")
	flag.Parse()

	urls, err := parseConfig(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	feeds := make([]reader.Feed, len(urls))
	for i, url := range urls {
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
