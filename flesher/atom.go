package flesher

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/araddon/dateparse"
	"golang.org/x/net/html/charset"
	"html/template"
)

type atomFeed struct {
	XMLName xml.Name `xml:"feed"`

	Title       string     `xml:"title"`
	Links       []atomLink `xml:"link"`
	Description string     `xml:"subtitle"`
	Items       []atomItem `xml:"entry"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Type string `xml:"type,attr"`
	Rel  string `xml:"rel,attr"`
}

type atomItem struct {
	Title       string        `xml:"title"`
	Links       []atomLink    `xml:"link"`
	Description template.HTML `xml:"summary"`
	Content     template.HTML `xml:"content"`
	Published   string        `xml:"published"`
	Updated     string        `xml:"updated"`
}

func parseAtomFeed(blob []byte) FeedResult {
	f := atomFeed{}

	d := xml.NewDecoder(bytes.NewReader(blob))
	d.CharsetReader = charset.NewReaderLabel
	if err := d.Decode(&f); err != nil {
		return feedResult{
			err: err,
		}
	}

	result := feedResult{
		title:   f.Title,
		entries: make([]EntryResult, len(f.Items)),
	}

	for _, link := range f.Links {
		if link.Rel == "alternate" {
			result.url = link.Href
			break
		}
	}

	for i, item := range f.Items {
		entry := entryResult{
			title: item.Title,
		}

		for _, link := range item.Links {
			if link.Rel == "alternate" {
				entry.url = link.Href
				break
			}
		}

		if item.Description != "" {
			entry.content = item.Description
		} else if item.Content != "" {
			entry.content = item.Content
		}

		var published string
		if item.Updated != "" {
			published = item.Updated
		} else {
			published = item.Published
		}

		t, err := dateparse.ParseAny(published)
		if err == nil {
			entry.published = t
		} else {
			result.err = fmt.Errorf("%s\nAtom date parsing error: %s\n", result.err, err)
		}

		result.entries[i] = entry
	}

	return result
}
