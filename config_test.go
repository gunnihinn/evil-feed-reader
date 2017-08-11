package main

import (
	"testing"
)

func TestConfigYaml(t *testing.T) {
	// Careful not to put \t here
	input := `
# comment
-
    url: a
-
    url: b
# comment
-
    url: c
`
	expected := []Config{
		Config{Resource: "a"},
		Config{Resource: "b"},
		Config{Resource: "c"},
	}

	got, err := parseConfigYaml([]byte(input))
	if err != nil {
		t.Errorf("Parse error: %s\n", err)
		return
	}

	if len(expected) != len(got) {
		t.Errorf("Got %v, expected %v\n", got, expected)
		return
	}

	for i, e := range expected {
		g := got[i]
		if e != g {
			t.Errorf("Got %v, expected %v\n", got, expected)
			return
		}
	}
}
