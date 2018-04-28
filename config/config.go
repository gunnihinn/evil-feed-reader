package config

import (
	"gopkg.in/yaml.v2"
	"io"
)

type Config struct {
	Feeds []Feed `yaml:"feeds"`
}

type Feed struct {
	URL      string `yaml:"url"`
	Nickname string `yaml:"nickname"`
	Prefix   string `yaml:"prefix"`
}

func Parse(r io.Reader) (Config, error) {
	config := Config{}

	decoder := yaml.NewDecoder(r)
	err := decoder.Decode(&config)

	return config, err
}
