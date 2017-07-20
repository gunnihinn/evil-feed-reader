package flesher

import (
	"testing"
)

func TestDetectRss(t *testing.T) {
	var RSSs = []string{
		`<rss></rss>`,
		`<?xml version="1.0" ?><rss></rss>`,
		`<?xml version="1.0" encoding="ISO-8859-1" ?><rss></rss>`,
	}

	for _, rss := range RSSs {
		if !isRss([]byte(rss)) {
			t.Errorf("'%s' is RSS", rss)
		}
	}
}

func TestDetectAtom(t *testing.T) {
	var atoms = []string{
		`<feed></feed>`,
		`<?xml version="1.0" ?><feed></feed>`,
		`<?xml version="1.0" encoding="ISO-8859-1" ?><feed></feed>`,
	}

	for _, atom := range atoms {
		if !isAtom([]byte(atom)) {
			t.Errorf("'%s' is Atom", atom)
		}
	}
}
