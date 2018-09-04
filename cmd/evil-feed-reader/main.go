package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/gunnihinn/evil-feed-reader/config"
	"github.com/gunnihinn/evil-feed-reader/core"

	log "github.com/sirupsen/logrus"
)

type Content struct {
	feeds []config.Feed
	Days  []DateEntries
}

type DateEntries struct {
	Date    time.Time
	Entries []core.Entry
}

func (content *Content) Refresh() {
	log.Info("Refreshing content")
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	feeds, errors := core.ScatterGather(ctx, content.feeds, core.HTTPFetcher)
	for _, err := range errors {
		log.Error(err)
	}

	entries := make([]core.Entry, 0)
	for _, feed := range feeds {
		entries = append(entries, feed...)
	}

	end := time.Now()
	log.WithFields(log.Fields{
		"took_ms": (end.UnixNano() - start.UnixNano()) / 1000000,
	}).Info("Got feeds")

	start = time.Now()
	content.Days = gatherEntries(entries)
	end = time.Now()
	log.WithFields(log.Fields{
		"took_us": (end.UnixNano() - start.UnixNano()) / 1000,
	}).Info("Gathered feeds")
}

func main() {
	log.Info("Starting")

	var port = flag.Int("port", 8080, "HTTP port")
	var configFileOption = flag.String("config", "", "Reader config file")
	flag.Parse()

	content := &Content{}
	configFile := findConfigFile(configFileOption)
	err := content.Load(configFile)
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
	}(configFile, reload)

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      nil,
		ReadTimeout:  200 * time.Millisecond,
		WriteTimeout: 3 * time.Second,
	}

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
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Info("Shutting down")
	if err := server.Close(); err != nil {
		log.Fatal(err)
	}
}

func gatherEntries(entries []core.Entry) []DateEntries {
	sort.Sort(sort.Reverse(core.SortedEntries(entries)))

	days := make([]DateEntries, 0)

	var date string
	var bucket DateEntries
	for _, entry := range entries {
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
		}

		bucket.Entries = append(bucket.Entries, entry)
	}

	if len(bucket.Entries) > 0 {
		days = append(days, bucket)
	}

	return days
}

func getDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func (c *Content) Load(filename string) error {
	log.WithFields(log.Fields{
		"filename": filename,
	}).Info("Loading config")

	fh, err := os.Open(filename)
	if err != nil {
		return err
	}

	cfg, err := config.Parse(fh)
	if err != nil {
		log.Fatal(err)
	}

	c.feeds = cfg.Feeds

	return nil
}

func findConfigFile(option *string) string {
	if option != nil && *option != "" {
		return *option
	}

	prefixes := make([]string, 0)

	// Check current working directory first
	prefixes = append(prefixes, "")

	// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html#variables
	if p := os.Getenv("XDG_CONFIG_HOME"); p != "" {
		prefixes = append(prefixes, p)
	} else {
		prefixes = append(prefixes, path.Join(os.Getenv("HOME"), ".config"))
	}

	if ps := os.Getenv("XDG_CONFIG_DIRS"); ps != "" {
		for _, p := range strings.Split(ps, ":") {
			prefixes = append(prefixes, p)
		}
	} else {
		prefixes = append(prefixes, "/etc/xdg")
	}

	for _, dir := range prefixes {
		fn := path.Join(dir, "evil-feed-reader.yaml")
		if _, err := os.Stat(fn); err == nil {
			return fn
		}
	}

	return "config.yaml"
}
