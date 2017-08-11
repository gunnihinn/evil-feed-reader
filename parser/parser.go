package parser

import (
	"html"
	"html/template"

	"github.com/mmcdole/gofeed"
)

type Parser struct {
	parser   *gofeed.Parser
	provider Provider
}

func New(p Provider) Parser {
	return Parser{
		parser:   gofeed.NewParser(),
		provider: p,
	}
}

func (p Parser) Parse(resource string) (FeedResult, error) {
	reader, err := p.provider(resource)
	if err != nil {
		return feedResult{}, err
	}
	defer reader.Close()

	feed, err := p.parser.Parse(reader)
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
