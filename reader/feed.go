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
			title:     itemResult.Title(),
			url:       itemResult.Url(),
			content:   itemResult.Content(),
			published: itemResult.Published(),
		}

		f.entries[i] = entry
	}

	h := sha1.New()
	for _, entry := range f.entries {
		io.WriteString(h, entry.Title())
		io.WriteString(h, entry.Published())
	}
	hash := hex.EncodeToString(h.Sum(nil))
	if f.hash != hash {
		f.hash = hash
		f.seen = false
	}
}

func (f feed) Seen() bool { return f.seen }

func (f *feed) SetSeen(seen bool) { f.seen = seen }

func (f feed) Hash() string { return f.hash }

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
