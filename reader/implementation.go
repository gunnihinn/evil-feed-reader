package reader

import (
	"crypto/sha1"
	"encoding/hex"
	"html/template"
	"io"
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

// Don't show more than this many entries per feed
const entryLimit = 10

type feed struct {
	resource string
	hash     string
	seen     bool

	provider provider.Provider
	parser   flesher.Parser

	title   string
	url     string
	entries []Entry
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

func (f *feed) SetState(state FeedState) {
	f.hash = state.Hash
	f.seen = state.Seen
}

func (f feed) Entries() []Entry {
	if len(f.entries) < entryLimit {
		return f.entries
	}

	return f.entries[0:entryLimit]
}

func (f *feed) Update() error {
	blob, err := f.provider.Get(f.resource)
	if err != nil {
		return err
	}

	if f.parser == nil {
		f.parser = flesher.New(blob)
	}

	feedResult, err := f.parser(blob)
	if err != nil {
		return err
	}

	if f.title == "" {
		f.title = feedResult.Title()
	}

	if f.url == "" {
		f.url = feedResult.Url()
	}

	f.entries = make([]Entry, 0, len(feedResult.Items()))
	for _, itemResult := range feedResult.Items() {
		f.entries = append(f.entries, entry{
			title:     itemResult.Title(),
			url:       itemResult.Url(),
			content:   itemResult.Content(),
			published: itemResult.Published(),
		})
	}

	if hash := f.calculateHash(); f.hash != hash {
		if f.hash == "" {
			f.seen = true
		} else {
			f.seen = false
		}
		f.hash = hash
	}

	return nil
}

func (f feed) Seen() bool { return f.seen }

func (f *feed) SetSeen(seen bool) { f.seen = seen }

func (f feed) Hash() string { return f.hash }

func (f feed) calculateHash() string {
	h := sha1.New()
	for _, entry := range f.entries {
		io.WriteString(h, entry.Title())
		io.WriteString(h, entry.Published())
	}

	return hex.EncodeToString(h.Sum(nil))
}

type entry struct {
	title     string
	url       string        // optional
	content   template.HTML // optional
	published string
}

func (e entry) Title() string {
	return e.title
}

func (e entry) Url() string {
	return e.url
}

func (e entry) Content() template.HTML {
	if len(strings.Split(string(e.content), " ")) > 300 {
		return "<p>This post was too long to comfortably fit onto the page.</p>"
	}

	return e.content
}

func (e entry) Published() string {
	return e.published
}