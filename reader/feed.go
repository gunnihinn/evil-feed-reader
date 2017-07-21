package reader

import (
	"html/template"
	"strings"
	"time"

	"github.com/gunnihinn/evil-feed-reader/flesher"
	"github.com/gunnihinn/evil-feed-reader/provider"
)

func New(provider provider.Provider, resource string) Feed {
	return &feed{
		resource: resource,
		provider: provider,
	}
}

// Don't show more than this many entries per feed
const entryLimit = 10

type feed struct {
	resource string
	provider provider.Provider
	parser   flesher.Parser

	// Generated at runtime
	title   string
	url     string
	entries []Entry
	err     error
}

func (f feed) Resource() string {
	return f.resource
}

func (f feed) Title() string {
	return f.title
}

func (f feed) Url() string {
	return f.url
}

func (f feed) Entries() []Entry {
	if len(f.entries) < entryLimit {
		return f.entries
	}

	return f.entries[0:entryLimit]
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

	if f.parser == nil {
		f.parser = flesher.New(blob)
	}

	feedResult := f.parser(blob)

	f.err = feedResult.Error()

	if f.title == "" {
		f.title = feedResult.Title()
	}

	if f.url == "" {
		f.url = feedResult.Url()
	}

	f.entries = make([]Entry, len(feedResult.Items()))
	for i, itemResult := range feedResult.Items() {
		entry := entry{
			title:   itemResult.Title(),
			url:     itemResult.Url(),
			content: itemResult.Content(),
		}

		f.entries[i] = entry
	}
}

func (f feed) HasRecentItems() bool {
	now := time.Now()
	var limit int64 = 24 * 60 * 60

	for _, entry := range f.Entries() {
		if now.Unix()-entry.Published().Unix() < limit {
			return true
		}
	}

	return false
}

type entry struct {
	title     string
	url       string        // optional
	content   template.HTML // optional
	published time.Time
}

func (e entry) Title() string {
	return e.title
}

func (e entry) Url() string {
	return e.url
}

func (e entry) Published() time.Time {
	return e.published
}

func (e entry) Content() template.HTML {
	if len(strings.Split(string(e.content), " ")) > 300 {
		return "<p>This post was too long to comfortably fit onto the page.</p>"
	}

	return e.content
}
