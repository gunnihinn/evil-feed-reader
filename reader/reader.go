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
	Error() error
	Update()
}

type Entry interface {
	Title() string
	Url() string
	Content() template.HTML
	Published() string
}
