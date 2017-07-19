package flesher

import (
	"encoding/xml"
	"fmt"
	"html/template"
)

// detect feed type, return appropriate parser
func New(blob []byte) Parser {
	if isRss(blob) {
		return parseRssFeed
	}

	if isAtom(blob) {
		return parseAtomFeed
	}

	return parseNothing
}

func isAtom(blob []byte) bool {
	type atom struct {
		XMLName xml.Name `xml:"feed"`
	}

	if err := xml.Unmarshal(blob, &atom{}); err != nil {
		return false
	}

	return true
}

func isRss(blob []byte) bool {
	type rss struct {
		XMLName xml.Name `xml:"rss"`
	}

	if err := xml.Unmarshal(blob, &rss{}); err != nil {
		return false
	}

	return true
}

func parseNothing(blob []byte) FeedResult {
	return feedResult{
		err: fmt.Errorf("Unable to determine feed type"),
	}
}

type FeedResult interface {
	Title() string
	Url() string
	Items() []EntryResult
	Error() error
}

type EntryResult interface {
	Title() string
	Url() string
	Content() template.HTML
}

type Parser func([]byte) FeedResult
