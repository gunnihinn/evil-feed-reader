package flesher

import (
	"encoding/xml"
	"testing"
)

func TestTitle(t *testing.T) {
	input := `
<title type="html">
<![CDATA[
Foo &amp; bar
]]>
</title>
`

	expected := atomTitle{
		Value: "\n\nFoo &amp; bar\n\n",
		Type:  "html",
	}

	got := atomTitle{}
	err := xml.Unmarshal([]byte(input), &got)

	if err != nil {
		t.Errorf("XML parse error:\n%s\n", err)
	} else {
		if got != expected {
			t.Errorf("Got: '%#v'\nExpected: '%#v'\n", got, expected)
		}

		if got.String() != "Foo & bar" {
			t.Errorf("Stringification error:\nGot '%s', expected '%s'", got.String(), "Foo & bar")
		}
	}
}
