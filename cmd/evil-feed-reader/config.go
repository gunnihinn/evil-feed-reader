package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gunnihinn/evil-feed-reader/reader"
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

func parseState(filename string) (map[string]reader.FeedState, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	blob, err := ioutil.ReadAll(fh)
	if err != nil {
		return nil, err
	}

	return reader.Unmarshal(blob)
}
