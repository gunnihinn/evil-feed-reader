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

	Title() string
	Url() string
	Entries() []Entry
	Update() error
}

type Entry interface {
	Title() string
	Url() string
	Content() template.HTML
	Published() string
}
