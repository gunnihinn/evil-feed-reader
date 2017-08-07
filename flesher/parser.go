package flesher

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"html/template"
)

// New detects feed type and returns an appropriate parser.
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

func parseNothing(blob []byte) (FeedResult, error) {
	return feedResult{}, fmt.Errorf("Unable to determine feed type")
}

// Parser parses a feed.
type Parser func([]byte) (FeedResult, error)
