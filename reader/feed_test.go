package reader

import (
	"testing"
	"time"
)

func TestRecentNew(t *testing.T) {
	e := entry{
		published: time.Now(),
	}
	f := feed{
		entries: []Entry{e},
	}

	if !f.HasRecentItems() {
		t.Errorf("Feed should have recent items")
	}
}

func TestRecentOld(t *testing.T) {
	e := entry{
		published: time.Unix(0, 0),
	}
	f := feed{
		entries: []Entry{e},
	}

	if f.HasRecentItems() {
		t.Errorf("Feed should not have recent items")
	}
}
