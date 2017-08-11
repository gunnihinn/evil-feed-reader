package parser

import (
	"html/template"
)

type FeedResult interface {
	Title() string
	Url() string
	Items() []EntryResult
}

type EntryResult interface {
	Title() string
	Url() string
	Content() template.HTML
	Published() string
}

type feedResult struct {
	title   string
	url     string
	entries []EntryResult
}

func (f feedResult) Title() string {
	return f.title
}

func (f feedResult) Url() string {
	return f.url
}

func (f feedResult) Items() []EntryResult {
	return f.entries
}

type entryResult struct {
	title     string
	url       string        // optional
	content   template.HTML // optional
	published string
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

func (e entryResult) Published() string {
	return e.published
}
