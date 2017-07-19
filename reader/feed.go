package reader

import (
	"html/template"
	"strings"

	"github.com/gunnihinn/evil-feed-reader/flesher"
	"github.com/gunnihinn/evil-feed-reader/provider"
)

func New(provider provider.Provider, resource string) Feed {
	return &feed{
		resource: resource,
		provider: provider,
	}
}

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
