package rss

import (
	"io/ioutil"
	"net/http"
)

func bytesFromUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	blob, err := ioutil.ReadAll(resp.Body)

	return blob, err
}

func fetchFeed(url string) (rssFeed, error) {
	blob, err := bytesFromUrl(url)
	if err != nil {
		return rssFeed{}, err
	}

	feed, err := parseFeed(blob)
	if err != nil {
		return rssFeed{}, err
	}

	return feed, nil
}
