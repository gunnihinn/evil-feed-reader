package flesher

// detect feed type, return appropriate parser
func New([]byte) Parser {}

type FeedResult interface {
	Title() string
	Url() string
	Items() []EntryResult
}

type EntryResult interface {
	Title() string
	Url() string
	Content() string
}

type Parser func([]byte) (Flesh, error)
