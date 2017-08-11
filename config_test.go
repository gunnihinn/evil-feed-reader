package main

import "testing"

func TestConfigYaml(t *testing.T) {
	input := `
# comment
- a
- b
# comment
- c
`
	expected := []string{
		"a",
		"b",
		"c",
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
