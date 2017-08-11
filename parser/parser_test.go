package parser

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"
)

func TestSavedFeeds(t *testing.T) {
	testInputDir := "test-inputs"

	type expect struct {
		title string
		url   string
	}

	expecteds := []expect{
		expect{
			title: "Podcasts (Radiolab)",
			url:   "http://www.radiolab.org/series/podcasts/",
		},
		expect{
			title: "Embedded in Academia",
			url:   "https://blog.regehr.org",
		},
		expect{
			title: "IT Hare on Soft.ware",
			url:   "http://ithare.com",
		},
		expect{
			title: "Daniel Lemire's blog",
			url:   "http://lemire.me/blog",
		},
		expect{
			title: "Schneier on Security",
			url:   "https://www.schneier.com/blog/",
		},
		expect{
			title: "programming is terrible",
			url:   "http://programmingisterrible.com/",
		},
		expect{
			title: "Still Drinking",
			url:   "http://www.stilldrinking.org/",
		},
		expect{
			title: "Errata Security",
			url:   "http://blog.erratasec.com/",
		},
		expect{
			title: "Accidentally Quadratic",
			url:   "https://accidentallyquadratic.tumblr.com/",
		},
		expect{
			title: "fuzzy notepad - blog",
			url:   "https://eev.ee/",
		},
		expect{
			title: "Julia Evans",
			url:   "http://jvns.ca",
		},
		expect{
			title: "High Scalability",
			url:   "http://highscalability.com/blog/",
		},
		expect{
			title: "Dan Luu",
			url:   "https://danluu.com/",
		},
		expect{
			title: "null program",
			url:   "http://nullprogram.com",
		},
		expect{
			title: "Bit Bashing",
			url:   "/",
		},
		expect{
			title: "New stuff",
			url:   "http://loup-vaillant.fr/",
		},
		expect{
			title: "flak rss",
			url:   "http://www.tedunangst.com/flak/",
		},
		expect{
			title: "Programming in the 21st Century",
			url:   "http://prog21.dadgum.com/",
		},
		expect{
			title: "Reply All",
			url:   "http://gimletmedia.com/shows/reply-all",
		},
		expect{
			title: "Welcome to Night Vale",
			url:   "http://welcometonightvale.com",
		},
	}

	files, err := ioutil.ReadDir(testInputDir)
	if err != nil {
		t.Error(err)
		return
	}

	if len(files) != len(expecteds) {
		t.Errorf("Expected %d test input files, only found %d\n", len(expecteds), len(files))
		return
	}

	p := New(File)
	for i, expected := range expecteds {
		testFile := path.Join(testInputDir, fmt.Sprintf("%d.xml", i))

		got, err := p.Parse(testFile)
		if err != nil {
			t.Errorf("Parse error in '%s': %s\n", testFile, err)
			return
		}

		if got.Title() != expected.title {
			t.Errorf("Title:\nGot:\t\t'%s'\nExpected:\t'%s'\n", got.Title(), expected.title)
			return
		}

		if got.Url() != expected.url {
			t.Errorf("Title:\nGot:\t\t'%s'\nExpected:\t'%s'\n", got.Url(), expected.url)
			return
		}
	}
}
