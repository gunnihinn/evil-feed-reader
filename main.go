package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var port = flag.Int("port", 8080, "HTTP port")
	var configFile = flag.String("config", "feeds.json", "Reader config file")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	fh, err := os.Open(*configFile)
	if err != nil {
		logger.Fatal(err)
	}

	type FeedConfig struct {
		URL      string
		Nickname string
		Prefix   string
	}
	var feedConfigs []FeedConfig
	decoder := json.NewDecoder(fh)
	err = decoder.Decode(&feedConfigs)
	if err != nil {
		logger.Fatal(err)
	}

	handler := NewHandler()
	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", *port),
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

		if err := t.Execute(w, nil); err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
	})
	handler.Handle("/static/", http.StripPrefix("/static/", http.FileServer(assetFS())))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatalf("Server error: %s\n", err)
		}
	}()

	<-stop
	logger.Print("Shutting down")
}
