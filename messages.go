package main

import (
	//"html/template"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"time"

	"github.com/mmcdole/gofeed"
)

type HTTP struct {
	config   Config
	response *http.Response
	err      error
}

type Entry struct {
	Feed      string
	Title     string
	Link      *url.URL
	Published time.Time
}

func parseEntries(message HTTP) ([]Entry, error) {
	if message.err != nil {
		return nil, message.err
	}

	parser := gofeed.NewParser()
	feed, err := parser.Parse(message.response.Body)
	if err != nil {
		return nil, err
	}

	es := make([]Entry, 0)
	for _, item := range feed.Items {
		var published time.Time
		if item.UpdatedParsed != nil {
			published = *item.UpdatedParsed
		} else if item.PublishedParsed != nil {
			published = *item.PublishedParsed
		} else {
			panic(fmt.Sprintf("No date in %s: %s\n", feed.Title, item.Title))
		}

		if published.Unix() < 20*365*24*60*60 {
			panic(fmt.Sprintf("No date in %s: %s\n", feed.Title, item.Title))
		}

		var u string
		if item.Link != "" {
			u = item.Link
		} else if len(item.Enclosures) > 0 {
			for _, e := range item.Enclosures {
				if e.URL != "" {
					u = e.URL
					break
				}
			}
		}

		base, err := url.Parse(message.config.Prefix)
		if err != nil {
			panic(fmt.Sprintf("%s\n", err))
		}
		ext, err := url.Parse(u)
		if err != nil {
			panic(fmt.Sprintf("%s\n", err))
		}

		var feedTitle string
		if message.config.Nickname != "" {
			feedTitle = message.config.Nickname
		} else {
			feedTitle = html.UnescapeString(feed.Title)
		}

		es = append(es, Entry{
			Feed:      feedTitle,
			Title:     html.UnescapeString(item.Title),
			Link:      base.ResolveReference(ext),
			Published: published,
		})
	}

	return es, nil
}

type sortedEntries []Entry

func (p sortedEntries) Len() int {
	return len(p)
}

func (p sortedEntries) Less(i, j int) bool {
	return p[i].Published.Before(p[j].Published)
}

func (p sortedEntries) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
