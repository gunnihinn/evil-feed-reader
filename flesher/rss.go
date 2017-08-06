package flesher

import (
	"bytes"
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"html/template"
)

type rssFeed struct {
	XMLName xml.Name `xml:"rss"`

	// required
	Title       string `xml:"channel>title"`
	Link        string `xml:"channel>link"`
	Description string `xml:"channel>description"`

	// optional
	Items []rssItem `xml:"channel>item"`
}

type rssItem struct {
	/*
	 * "All elements of an item are optional, however at least one of title or
	 * description must be present."
	 * http://cyber.harvard.edu/rss/rss.html
	 */
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	Description template.HTML `xml:"description"` // lol security
	Content     template.HTML `xml:"encoded"`
	PubDate     string        `xml:"pubDate"`
}

func parseRssFeed(blob []byte) FeedResult {
	f := rssFeed{}

	d := xml.NewDecoder(bytes.NewReader(blob))
	d.CharsetReader = charset.NewReaderLabel
	if err := d.Decode(&f); err != nil {
		return feedResult{
			err: err,
		}
	}

	result := feedResult{
		title:   f.Title,
		url:     f.Link,
		entries: make([]EntryResult, len(f.Items)),
	}

	for i, item := range f.Items {
		entry := entryResult{
			title:     item.Title,
			url:       item.Link,
			published: item.PubDate,
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
