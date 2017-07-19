package flesher

import (
	"html/template"
)

type feedResult struct {
	title   string
	url     string
	entries []entryResult
	err     error
}

func (f feedResult) Title() string {
	return f.title
}

func (f feedResult) Url() string {
	return f.url
}

func (f feedResult) Entries() []entryResult {
	return f.entries
}

func (f feedResult) Error() error {
	return f.err
}

type entryResult struct {
	title   string
	url     string        // optional
	content template.HTML // optional
}

func (e entryResult) Title() string {
	return e.title
}

func (e entryResult) Url() string {
	return e.url
}

func (e entryResult) Content() template.HTML {
	return e.content
}
