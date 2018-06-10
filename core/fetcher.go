package core

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"github.com/gunnihinn/evil-feed-reader/config"
)

// Fetcher fetches a Feed based on a config.
type Fetcher func(config.Feed) ([]Entry, error)

// ScatterGather uses a given Fetcher to fetch feeds concurrently.
func ScatterGather(configs []config.Feed, goodboy Fetcher) ([][]Entry, []error) {
	fs := make(chan []Entry)
	es := make(chan error)
	defer close(fs)
	defer close(es)

	var wg sync.WaitGroup
	for _, cfg := range configs {
		wg.Add(1)
		go func(c config.Feed) {
			f, err := goodboy(c)
			if err != nil {
				es <- err
			} else {
				fs <- f
			}
		}(cfg)
	}

	feeds := make([][]Entry, 0, len(configs))
	go func() {
		for feed := range fs {
			feeds = append(feeds, feed)
			wg.Done()
		}
	}()

	errors := make([]error, 0, len(configs))
	go func() {
		for err := range es {
			errors = append(errors, err)
			wg.Done()
		}
	}()

	wg.Wait()

	return feeds, errors
}

var client = http.Client{
	Transport: &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		ResponseHeaderTimeout: 2 * time.Second,
	},
}

// HTTPFetcher fetches feeds over HTTP.
func HTTPFetcher(f config.Feed) ([]Entry, error) {
	response, err := client.Get(f.URL)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return Parse(f, response.Body)
}
