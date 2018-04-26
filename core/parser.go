package core

import (
	"fmt"
	"html"
	"io"
	"net/url"
	"time"

	"github.com/gunnihinn/evil-feed-reader/config"

	"github.com/mmcdole/gofeed"
)

type Entry struct {
	Feed      string
	Title     string
	Link      *url.URL
	Published time.Time
}

func Parse(cfg config.Feed, raw io.Reader) ([]Entry, error) {
	parser := gofeed.NewParser()
	rawFeed, err := parser.Parse(raw)
	if err != nil {
		return nil, err
	}

	es := make([]Entry, 0)
	for _, item := range rawFeed.Items {
		var published time.Time
		if item.UpdatedParsed != nil {
			published = *item.UpdatedParsed
		} else if item.PublishedParsed != nil {
			published = *item.PublishedParsed
		} else {
			// TODO: Don't panic
			panic(fmt.Sprintf("No date in %s: %s\n", rawFeed.Title, item.Title))
		}

		if published.Unix() < 20*365*24*60*60 {
			// TODO: Don't panic
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

		base, err := url.Parse(cfg.Prefix)
		if err != nil {
			// TODO: Don't panic
			panic(fmt.Sprintf("%s\n", err))
		}
		ext, err := url.Parse(u)
		if err != nil {
			// TODO: Don't panic
			panic(fmt.Sprintf("%s\n", err))
		}

		var feedTitle string
		if cfg.Nickname != "" {
			feedTitle = cfg.Nickname
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

	return es, nil
}

type SortedEntries []Entry

func (p SortedEntries) Len() int {
	return len(p)
}

func (p SortedEntries) Less(i, j int) bool {
	return p[i].Published.Before(p[j].Published)
}

func (p SortedEntries) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
