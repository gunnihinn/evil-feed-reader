package reader

import (
	"html/template"
)

type Feed interface {
	Resource() string
	Seen() bool
	SetSeen(bool)
	Hash() string
	SetState(FeedState)
	Nickname() string

	Title() string
	Url() string
	Entries() []Entry
	Update() ([]string, error)
}

type Entry interface {
	Title() string
	Url() string
	Content() template.HTML
	Published() string
}
