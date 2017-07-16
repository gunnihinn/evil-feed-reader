package rss

import (
	"github.com/gunnihinn/evil-rss-reader/reader"
	"html/template"
)

func New(url string) reader.Feed {
	return &feed{
		url: url,
	}
}

type feed struct {
	url   string
	title string

	// Generated at runtime
	entries []reader.Entry
	err     error
}

func (f feed) Title() string {
	return f.title
}

func (f feed) Url() string {
	return f.url
}

func (f feed) Entries() []reader.Entry {
	return f.entries
}

func (f feed) Error() error {
	return f.err
}

func (f *feed) Update() {
	rf, err := fetchFeed(f.Url())

	if err == nil {
		if f.title == "" {
			f.title = rf.Title
		}

		f.entries = make([]reader.Entry, len(rf.Items))
		for i, item := range rf.Items {
			entry := entry{
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

type entry struct {
	title   string
	url     string        // optional
	content template.HTML // optional
}

func (e entry) Title() string {
	return e.title
}

func (e entry) Url() string {
	return e.url
}

func (e entry) Content() template.HTML {
	return e.content
}
