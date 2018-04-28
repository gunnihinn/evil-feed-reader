package config

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	type setup struct {
		input    string
		expected []Feed
	}

	tests := []setup{
		setup{
			input: `[
	{
		"Url": "URL",
		"Nickname": "NICKNAME",
		"Prefix": "PREFIX"
	}
]`,
			expected: []Feed{
				Feed{
					URL:      "URL",
					Nickname: "NICKNAME",
					Prefix:   "PREFIX",
				},
			},
		},
		setup{
			input: `[
	{
		"url": "URL",
		"nickname": "NICKNAME",
		"prefix": "PREFIX"
	}
]`,
			expected: []Feed{
				Feed{
					URL:      "URL",
					Nickname: "NICKNAME",
					Prefix:   "PREFIX",
				},
			},
		},
		setup{
			input: `[
	{
		"url": "URL"
	}
]`,
			expected: []Feed{
				Feed{
					URL: "URL",
				},
			},
		},
	}

	for _, test := range tests {
		got, err := Parse(strings.NewReader(test.input))
		if err != nil {
			t.Error(err)
		}

		if !equal(got, test.expected) {
			t.Error("Didn't get what we expected")
		}
	}
}

func equal(a, b []Feed) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
