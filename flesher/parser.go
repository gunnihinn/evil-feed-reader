package flesher

import (
	"encoding/xml"
	"fmt"
)

// detect feed type, return appropriate parser
func New(blob []byte) Parser {
	type atom struct {
		XMLName xml.Name `xml:"feed"`
	}

	type rss struct {
		XMLName xml.Name `xml:"rss"`
	}

	err := xml.Unmarshal(blob, &atom{})
	if err == nil {
		return parseAtomFeed
	}

	err = xml.Unmarshal(blob, &rss{})
	if err == nil {
		return parseRssFeed
	}

	return parseNothing
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
	Content() string
}

type Parser func([]byte) FeedResult
