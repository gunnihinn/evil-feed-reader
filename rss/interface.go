package rss

import (
	"html/template"
	"strings"

	"github.com/gunnihinn/evil-rss-reader/provider"
	"github.com/gunnihinn/evil-rss-reader/reader"
)

func New(provider provider.Provider, resource string) reader.Feed {
	return &feed{
		resource: resource,
		provider: provider,
	}
}

type feed struct {
	resource string
	provider provider.Provider

	// Generated at runtime
	title   string
	url     string
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
	blob, err := f.provider.Get(f.resource)
	if err != nil {
		f.err = err
		return
	}

	rf, err := parseFeed(blob)
	if err == nil {
		if f.title == "" {
			f.title = rf.Title
		}

		if f.url == "" {
			f.url = rf.Link
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
	if len(strings.Split(string(e.content), " ")) > 300 {
		return ""
	}

	return e.content
}
