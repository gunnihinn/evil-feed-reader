package main

import (
	"io/ioutil"
	"os"

	"github.com/gunnihinn/evil-feed-reader/reader"
	"gopkg.in/yaml.v2"
)

func parseConfig(filename string) ([]string, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return []string{}, err
	}
	defer fh.Close()

	blob, err := ioutil.ReadAll(fh)
	if err != nil {
		return []string{}, err
	}

	return parseConfigYaml(blob)
}

func parseConfigYaml(blob []byte) ([]string, error) {
	urls := make([]string, 0)
	err := yaml.Unmarshal(blob, &urls)
	if err != nil {
		return []string{}, err
	}

	return urls, nil
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

func writeState(filename string, feeds []reader.Feed) error {
	blob, err := reader.Marshal(feeds)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(filename, blob, 0644); err != nil {
		return err
	}

	return nil
}

func readLog(filename string) (string, error) {
	if filename == "" {
		return "Logging to standard output", nil
	}

	fh, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer fh.Close()

	blob, err := ioutil.ReadAll(fh)
	if err != nil {
		return "", err
	}

	return string(blob), err
}
