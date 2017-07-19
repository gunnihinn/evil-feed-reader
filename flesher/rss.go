package flesher

import (
	"encoding/xml"
	"html/template"
	"strings"

	"github.com/gunnihinn/evil-feed-reader/provider"
	"github.com/gunnihinn/evil-feed-reader/reader"
)

type rssFeedResult struct {
	title   string
	url     string
	entries []rssEntryResult
	err     error
}

func (f feed) Title() string {
	return f.title
}

func (f feed) Url() string {
	return f.url
}

func (f feed) Entries() []rssEntryResult {
	return f.entries
}

func (f rssFeedResult) Error() error {
	return f.err
}

type rssEntryResult struct {
	title   string
	url     string        // optional
	content template.HTML // optional
}

func (e entry) Title() string {
	return e.title
}

func (e entry) Url() string {
	return e.url
}

func (e entry) Content() template.HTML {
	return e.content
}

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
}

func parseRssFeed(blob []byte) rssFeedResult {
	f := rssFeed{}
	err := xml.Unmarshal(blob, &f)

	if err != nil {
		f.err = err
		return f
	}

	rf, err := parseFeed(blob)
	if err == nil {
		if f.title == "" {
			f.title = rf.Title
		}

		if f.url == "" {
			f.url = rf.Link
		}

		f.entries = make([]rssEntryResult, len(rf.Items))
		for i, item := range rf.Items {
			entry := entry{
				title: item.Title,
				url:   item.Link,
			}
			if item.Description != "" {
				entry.content = item.Description
			} else if item.Content != "" {
				entry.content = item.Content
			}
			f.entries[i] = entry
		}
	}

	f.err = err

	return f, err
}
