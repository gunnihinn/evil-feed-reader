package atom

import (
	"encoding/xml"
	"html/template"
)

type atomFeed struct {
	XMLName xml.Name `xml:"feed"`

	Title       string     `xml:"title"`
	Links       []atomLink `xml:"link"`
	Description string     `xml:"subtitle"`
	Items       []atomItem `xml:"entry"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Type string `xml:"type,attr"`
	Rel  string `xml:"rel,attr"`
}

type atomItem struct {
	Title       string        `xml:"title"`
	Links       []atomLink    `xml:"link"`
	Description template.HTML `xml:"summary"`
	Content     template.HTML `xml:"content"`
}

func parseFeed(feed []byte) (atomFeed, error) {
	f := atomFeed{}
	err := xml.Unmarshal(feed, &f)

	return f, err
}
