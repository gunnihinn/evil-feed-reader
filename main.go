package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/gunnihinn/evil-feed-reader/config"
)

var logger *log.Logger

type Content struct {
	configs []config.Feed
	Days    []DateEntries
}

type DateEntries struct {
	Date    time.Time
	Entries []Entry
}

func SetupContent(configs []config.Feed) *Content {
	return &Content{
		configs: configs,
	}
}

func (content *Content) Refresh() {
	logger.Println("Refreshing content")
	start := time.Now()

	var client = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	https := make(chan HTTP)
	feeds := make(chan Feed)
	entries := make([]Entry, 0)

	for _, cfg := range content.configs {
		go func(c config.Feed) {
			response, err := client.Get(c.URL)

			https <- HTTP{
				config:   c,
				response: response,
				err:      err,
			}
		}(cfg)

		go func(https chan HTTP) {
			feed := parseEntries(<-https)
			feeds <- feed
		}(https)
	}

	i := 0
	for feed := range feeds {
		i++

		if feed.err != nil {
			log.Printf("Error: %s\n", feed.err)
			continue
		}

		entries = append(entries, feed.entries...)

		if i == len(content.configs) {
			close(https)
			close(feeds)
		}
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
	var configFile = flag.String("config", "feeds.json", "Reader config file")
	flag.Parse()

	fh, err := os.Open(*configFile)
	if err != nil {
		logger.Fatal(err)
	}

	configs, err := config.Parse(fh)
	if err != nil {
		logger.Fatal(err)
	}

	content := SetupContent(configs)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT)

	server := setupServer(*port, content)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatalf("Server error: %s\n", err)
		}
	}()

	<-stop
	logger.Print("Shutting down")
}

func gatherEntries(entries []Entry) []DateEntries {
	sort.Sort(sort.Reverse(sortedEntries(entries)))

	days := make([]DateEntries, 0)

	var date string
	var bucket DateEntries
	for _, entry := range entries {
		if date == "" {
			date = getDate(entry.Published)
			bucket = DateEntries{
				Date:    entry.Published,
				Entries: make([]Entry, 0),
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
				Entries: make([]Entry, 0),
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

func setupServer(port int, content *Content) *http.Server {
	handler := NewHandler()
	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: handler,
	}

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := Asset("static/index.html")
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}

		t := template.New("static/index.html")
		_, err = t.Parse(string(data))
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
	handler.Handle("/static/", http.StripPrefix("/static/", http.FileServer(assetFS())))

	return server
}
