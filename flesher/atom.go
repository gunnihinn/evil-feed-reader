package flesher

import (
	"encoding/xml"
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
}

func parseAtomFeed(blob []byte) FeedResult {
	f := atomFeed{}
	if err := xml.Unmarshal(blob, &f); err != nil {
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
		result.entries[i] = entry
	}

	return result
}
