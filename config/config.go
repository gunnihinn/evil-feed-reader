package config

import (
	"encoding/json"
	"io"
)

type Feed struct {
	URL      string
	Nickname string
	Prefix   string
}

func Parse(r io.Reader) ([]Feed, error) {
	var configs []Feed
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&configs)

	return configs, err
}
