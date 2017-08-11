package reader

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/gunnihinn/evil-feed-reader/parser"
)

func New(provider parser.Provider, resource string) Feed {
	return &feed{
		resource: resource,
		parser:   parser.New(provider),
	}
}

// Don't show more than this many entries per feed
const entryLimit = 10

type feed struct {
	resource string
	hash     string
	seen     bool

	parser parser.Parser

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

func (f *feed) Update() ([]string, error) {
	messages := make([]string, 0)

	feedResult, err := f.parser.Parse(f.resource)
	if err != nil {
		return messages, err
	}

	if title := feedResult.Title(); f.title == "" && title != "" {
		f.title = title
		messages = append(messages, fmt.Sprintf("Setting feed title to '%s'", title))
	} else {
		messages = append(messages, fmt.Sprintf("No title found"))
	}

	if url := feedResult.Url(); f.url == "" && url != "" {
		f.url = url
		messages = append(messages, fmt.Sprintf("Setting feed URL to '%s'", url))
	} else {
		messages = append(messages, fmt.Sprintf("No URL found"))
	}

	if f.url == "" || f.title == "" {
		msg := fmt.Sprintf("%+v", feedResult)

		if words := strings.Split(msg, " "); len(words) > 100 {
			msg = strings.Join(words[0:100], " ")
			msg += " [...]"
		}

		messages = append(messages, msg)
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
			messages = append(messages, "New feed defaults to 'seen'")
		} else {
			f.seen = false
			messages = append(messages, "New items in feed; marking feed 'not seen'")
		}
		f.hash = hash
	}

	return messages, nil
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
