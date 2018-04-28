package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/gunnihinn/evil-feed-reader/config"
	"github.com/gunnihinn/evil-feed-reader/core"
)

var logger *log.Logger

type Content struct {
	feeds []config.Feed
	Days  []DateEntries
}

type DateEntries struct {
	Date    time.Time
	Entries []core.Entry
}

func (content *Content) Refresh() {
	logger.Println("Refreshing content")
	start := time.Now()

	feeds, errors := core.ScatterGather(content.feeds, core.HTTPFetcher)
	for _, err := range errors {
		log.Printf("Error: %s\n", err)
	}

	entries := make([]core.Entry, 0)
	for _, feed := range feeds {
		entries = append(entries, feed...)
	}

	end := time.Now()
	logger.Printf("Getting feeds took %d ms\n", (end.UnixNano()-start.UnixNano())/1000000)

	start = time.Now()
	content.Days = gatherEntries(entries)
	end = time.Now()
	logger.Printf("Gathering entries took %d us\n", (end.UnixNano()-start.UnixNano())/1000)
}

func main() {
	logger = log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("Starting")

	var port = flag.Int("port", 8080, "HTTP port")
	var configFile = flag.String("config", "config.yaml", "Reader config file")
	flag.Parse()

	content := &Content{}
	err := content.Load(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT)

	reload := make(chan os.Signal, 1)
	signal.Notify(reload, syscall.SIGHUP)
	go func(filename string, sig chan os.Signal) {
		for range sig {
			if err := content.Load(filename); err != nil {
				log.Fatal(err)
			}
		}
	}(*configFile, reload)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("index.html").Parse(HTML)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}

		content.Refresh()

		if err := t.Execute(w, content); err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
	})

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
			logger.Fatalf("Server error: %s\n", err)
		}
	}()

	<-stop
	logger.Print("Shutting down")
}

func gatherEntries(entries []core.Entry) []DateEntries {
	sort.Sort(sort.Reverse(core.SortedEntries(entries)))

	days := make([]DateEntries, 0)

	var date string
	var bucket DateEntries
	for _, entry := range entries {
		if date == "" {
			date = getDate(entry.Published)
			bucket = DateEntries{
				Date:    entry.Published,
				Entries: make([]core.Entry, 0),
			}
		}

		d := getDate(entry.Published)
		if date != d {
			date = d

			if len(bucket.Entries) > 0 {
				days = append(days, bucket)
			}

			bucket = DateEntries{
				Date:    entry.Published,
				Entries: make([]core.Entry, 0),
			}
		} else {
			bucket.Entries = append(bucket.Entries, entry)
		}
	}

	return days
}

func getDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func (c *Content) Load(filename string) error {
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}

	cfg, err := config.Parse(fh)
	if err != nil {
		logger.Fatal(err)
	}

	c.feeds = cfg.Feeds

	return nil
}
