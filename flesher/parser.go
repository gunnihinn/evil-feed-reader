package flesher

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
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

	d := xml.NewDecoder(bytes.NewReader(blob))
	d.CharsetReader = charset.NewReaderLabel
	if err := d.Decode(&atom{}); err != nil {
		return false
	}

	return true
}

func isRss(blob []byte) bool {
	type rss struct {
		XMLName xml.Name `xml:"rss"`
	}

	d := xml.NewDecoder(bytes.NewReader(blob))
	d.CharsetReader = charset.NewReaderLabel
	if err := d.Decode(&rss{}); err != nil {
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
	Published() string
}

type Parser func([]byte) FeedResult
