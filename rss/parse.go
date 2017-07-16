package rss

import (
	"encoding/xml"
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
}

func parseFeed(feed []byte) (rssFeed, error) {
	f := rssFeed{}
	err := xml.Unmarshal(feed, &f)

	return f, err
}
