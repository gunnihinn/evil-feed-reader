package main

import (
	"html/template"
	"net/http"
	"net/url"
	"time"
)

type HTTP struct {
	config   Config
	response *http.Response
	err      error
}

type Entry struct {
	feed      template.HTML
	title     template.HTML
	link      url.URL
	published time.Time
}

func parseEntries(message HTTP) ([]Entry, error) {
	if message.err != nil {
		return nil, message.err
	}

	es := make([]Entry, 0)

	return es, nil
}
