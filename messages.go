package main

import (
	"fmt"
	"html"
	"net/http"
	"net/url"
	"time"

	"github.com/gunnihinn/evil-feed-reader/config"

	"github.com/mmcdole/gofeed"
)

type HTTP struct {
	config   config.Feed
	response *http.Response
	err      error
}

type Feed struct {
	config  config.Feed
	entries []Entry
	err     error
}

type Entry struct {
	Feed      string
	Title     string
	Link      *url.URL
	Published time.Time
}

func parseEntries(message HTTP) Feed {
	if message.err != nil {
		return Feed{
			config: message.config,
			err:    message.err,
		}
	}

	parser := gofeed.NewParser()
	rawFeed, err := parser.Parse(message.response.Body)
	if err != nil {
		return Feed{
			config: message.config,
			err:    err,
		}
	}

	es := make([]Entry, 0)
	for _, item := range rawFeed.Items {
		var published time.Time
		if item.UpdatedParsed != nil {
			published = *item.UpdatedParsed
		} else if item.PublishedParsed != nil {
			published = *item.PublishedParsed
		} else {
			panic(fmt.Sprintf("No date in %s: %s\n", rawFeed.Title, item.Title))
		}

		if published.Unix() < 20*365*24*60*60 {
			panic(fmt.Sprintf("No date in %s: %s\n", rawFeed.Title, item.Title))
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
			feedTitle = html.UnescapeString(rawFeed.Title)
		}

		es = append(es, Entry{
			Feed:      feedTitle,
			Title:     html.UnescapeString(item.Title),
			Link:      base.ResolveReference(ext),
			Published: published,
		})
	}

	return Feed{
		config:  message.config,
		entries: es,
	}
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
