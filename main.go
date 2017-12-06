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

type Config struct {
	URL      string
	Nickname string
	Prefix   string
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT)

	server := setupServer(*port)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatalf("Server error: %s\n", err)
		}
	}()

	<-stop
	logger.Print("Shutting down")
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

func setupServer(port int) *http.Server {
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

		if err := t.Execute(w, nil); err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
	})
	handler.Handle("/static/", http.StripPrefix("/static/", http.FileServer(assetFS())))

	return server
}
