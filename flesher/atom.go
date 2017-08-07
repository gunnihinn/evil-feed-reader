package flesher

import (
	"bytes"
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"html"
	"html/template"
	"strings"
)

type atomFeed struct {
	XMLName xml.Name `xml:"feed"`

	Title       atomTitle  `xml:"title"`
	Links       []atomLink `xml:"link"`
	Description string     `xml:"subtitle"`
	Items       []atomItem `xml:"entry"`
}

type atomTitle struct {
	Value string `xml:",chardata"`
	Type  string `xml:"type,attr"`
}

func (a atomTitle) String() string {
	var title string
	if a.Type == "html" {
		title = html.UnescapeString(a.Value)
	} else {
		title = a.Value
	}

	return strings.TrimSpace(title)
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

type atomItem struct {
	Title       atomTitle     `xml:"title"`
	Links       []atomLink    `xml:"link"`
	Description template.HTML `xml:"summary"`
	Content     template.HTML `xml:"content"`
	Published   string        `xml:"published"`
	Updated     string        `xml:"updated"`
}

func getLink(links []atomLink) string {
	/*
		https://tools.ietf.org/html/rfc4287

		4.2.7.2.  The "rel" Attribute

		atom:link elements MAY have a "rel" attribute that indicates the link
		relation type.  If the "rel" attribute is not present, the link
		element MUST be interpreted as if the link relation type is
		"alternate".
	*/

	for _, link := range links {
		if link.Rel == "alternate" || link.Rel == "" {
			return link.Href
		}
	}

	return ""
}

func parseAtomFeed(blob []byte) (FeedResult, error) {
	f := atomFeed{}

	d := xml.NewDecoder(bytes.NewReader(blob))
	d.CharsetReader = charset.NewReaderLabel
	if err := d.Decode(&f); err != nil {
		// TODO: Wrap error
		return feedResult{}, err
	}

	result := feedResult{
		title:   f.Title.String(),
		entries: make([]EntryResult, len(f.Items)),
		url:     getLink(f.Links),
	}

	for i, item := range f.Items {
		entry := entryResult{
			title: item.Title.String(),
			url:   getLink(item.Links),
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
		entry.published = published

		result.entries[i] = entry
	}

	return result, nil
}
