package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gunnihinn/evil-feed-reader/parser"
	"github.com/gunnihinn/evil-feed-reader/reader"
)

type Context struct {
	Sidebar []Navitem
	Feed    reader.Feed
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
	var configFile = flag.String("config", "evil.yaml", "Reader config file")
	var stateFile = flag.String("state", ".evil-state.json", "Internal state file")
	var logFile = flag.String("log", "", "Reader log file; use STDOUT if absent")
	flag.Parse()

	logger := Logger(*logFile)

	urls, err := parseConfig(*configFile)
	if err != nil {
		logger.Printf("Couldn't parse config file: %s\n", err)
	}

	feeds := make([]reader.Feed, 0)
	for _, url := range urls {
		feeds = append(feeds, reader.New(parser.HTTP, url))
	}

	state, err := parseState(*stateFile)
	if err != nil {
		logger.Printf("Couldn't parse state file: %s\n", err)
		state = make(map[string]reader.FeedState)
	}

	for _, feed := range feeds {
		s, ok := state[feed.Resource()]
		if ok {
			feed.SetState(s)
		}
	}

	go func() {
		for {
			for _, feed := range feeds {
				go func(f reader.Feed) {
					messages, err := f.Update()
					if err != nil {
						logger.Printf("Problems parsing feed '%s':\n%s", f.Resource(), err)
					} else if len(f.Entries()) == 0 {
						logger.Printf("Got no entries from '%s'\n", f.Resource())
					}

					for _, msg := range messages {
						logger.Print(msg)
					}
				}(feed)
			}
			time.Sleep(15 * time.Minute)
		}
	}()

	addr := fmt.Sprintf(":%d", *port)
	logger.Printf("Listening on %s\n", addr)

	handler := NewHandler()
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	for i, feed := range feeds {
		handler.HandleFunc(fmt.Sprintf("/%d", i), createHandler(feeds, feed))
	}
	handler.HandleFunc("/", createHandler(feeds, nil))
	handler.Handle("/static/", http.StripPrefix("/static/", http.FileServer(assetFS())))
	handler.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		contents, err := readLog(*logFile)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
		fmt.Fprintf(w, "%s", contents)
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatalf("Server error: %s\n", err)
		}
	}()

	<-stop
	logger.Print("Shutting down")

	if err := writeState(*stateFile, feeds); err != nil {
		logger.Fatalf("Couldn't write reader state to disk: %s\n", err)
	}
}

func Logger(filename string) *log.Logger {
	file := os.Stdout
	if filename != "" {
		fh, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		if err == nil {
			file = fh
		} else {
			log.Printf("Can't log to '%s': %s\n", filename, err)
			log.Print("Logging to standard output")
		}
	}

	return log.New(file, "", log.LstdFlags)
}
