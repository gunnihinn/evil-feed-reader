package rss

import (
	"github.com/gunnihinn/evil-rss-reader/reader"
	"html/template"
)

func New(url string) *Feed {
	return &Feed{
		url: url,
	}
}

type Feed struct {
	url   string
	title string

	// Generated at runtime
	entries []reader.Entry
	err     error
}

func (f Feed) Title() string {
	return f.title
}

func (f Feed) Url() string {
	return f.url
}

func (f Feed) Entries() []reader.Entry {
	return f.entries
}

func (f Feed) Error() error {
	return f.err
}

func (f *Feed) Update() {
	rf, err := fetchFeed(f.Url())

	if err == nil {
		if f.title == "" {
			f.title = rf.Title
		}

		f.entries = make([]reader.Entry, len(rf.Items))
		for i, item := range rf.Items {
			entry := Entry{
				title: item.Title,
				url:   item.Link,
			}
			if item.Description != "" {
				entry.content = item.Description
			} else if item.Content != "" {
				entry.content = item.Content
			}
			f.entries[i] = entry
		}
	}

	f.err = err
}

type Entry struct {
	title   string
	url     string        // optional
	content template.HTML // optional
}

func (e Entry) Title() string {
	return e.title
}

func (e Entry) Url() string {
	return e.url
}

func (e Entry) Content() template.HTML {
	return e.content
}
