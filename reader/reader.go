package reader

import (
	"html/template"
	"time"
)

type Feed interface {
	Resource() string
	Title() string
	Url() string
	Entries() []Entry
	Error() error
	Update()
	HasRecentItems() bool
}

type Entry interface {
	Title() string
	Url() string
	Content() template.HTML
	Published() time.Time
}
