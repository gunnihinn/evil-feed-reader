package main

import (
	"crypto/tls"
	"encoding/json"
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
)

type Config struct {
	URL      string
	Nickname string
	Prefix   string
}

type Content struct {
	Days []DateEntries
}

type DateEntries struct {
	Date    string
	Entries []Entry
}

func main() {
	var port = flag.Int("port", 8080, "HTTP port")
	var configFile = flag.String("config", "feeds.json", "Reader config file")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	configs, err := parseConfigFile(*configFile)
	if err != nil {
		logger.Fatal(err)
	}

	var client = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	logger.Println("START")

	https := make(chan HTTP)
	entries := make([]Entry, 0)

	for _, config := range configs {
		go func(cfg Config) {
			response, err := client.Get(cfg.URL)

			https <- HTTP{
				config:   cfg,
				response: response,
				err:      err,
			}
		}(config)
	}

	// TODO: Use wait groups
	i := 0
	for msg := range https {
		es, err := parseEntries(msg)
		if err != nil {
			log.Fatal(err)
		}

		entries = append(entries, es...)

		i++
		if i == len(configs) {
			break
		}
	}
	logger.Println("END")

	content := gatherEntries(entries)

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

func gatherEntries(entries []Entry) Content {
	sort.Sort(sort.Reverse(sortedEntries(entries)))

	content := Content{
		Days: make([]DateEntries, 0),
	}

	var date string
	var bucket DateEntries
	for _, entry := range entries {
		if date == "" {
			date = getDate(entry.Published)
			bucket = DateEntries{
				Date:    date,
				Entries: make([]Entry, 0),
			}
		}

		d := getDate(entry.Published)
		if date != d {
			date = d

			if len(bucket.Entries) > 0 {
				content.Days = append(content.Days, bucket)
			}

			bucket = DateEntries{
				Date:    date,
				Entries: make([]Entry, 0),
			}
		} else {
			bucket.Entries = append(bucket.Entries, entry)
		}
	}

	return content
}

func getDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func parseConfigFile(filename string) ([]Config, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var configs []Config
	decoder := json.NewDecoder(fh)
	err = decoder.Decode(&configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func setupServer(port int, content Content) *http.Server {
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

		if err := t.Execute(w, content); err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
	})
	handler.Handle("/static/", http.StripPrefix("/static/", http.FileServer(assetFS())))

	return server
}
