package flesher

import (
	"bytes"
	"html"
	"html/template"

	"github.com/mmcdole/gofeed"
)

// New detects feed type and returns an appropriate parser.
func New(blob []byte) Parser {
	return parseFeed
}

func parseFeed(blob []byte) (FeedResult, error) {
	feed, err := gofeed.NewParser().Parse(bytes.NewReader(blob))
	if err != nil {
		return feedResult{}, err
	}

	result := feedResult{
		title: feed.Title,
		url:   feed.Link,
	}

	entries := make([]EntryResult, 0)
	for _, item := range feed.Items {
		var content string
		if item.Description != "" {
			content = item.Description
		} else {
			content = item.Content
		}

		var published string
		if item.Updated != "" {
			published = item.Updated
		} else {
			published = item.Published
		}

		entries = append(entries, entryResult{
			title:     html.UnescapeString(item.Title),
			url:       item.Link,
			content:   template.HTML(content),
			published: published,
		})
	}

	result.entries = entries

	return result, nil
}

// Parser parses a feed.
type Parser func([]byte) (FeedResult, error)
