package config

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	type setup struct {
		input    string
		expected Config
	}

	tests := []setup{
		setup{
			input: `feeds:
- url: URL
  nickname: NICKNAME
  prefix: PREFIX
`,
			expected: Config{
				Feeds: []Feed{
					Feed{
						URL:      "URL",
						Nickname: "NICKNAME",
						Prefix:   "PREFIX",
					},
				},
			},
		},
		setup{
			input: `feeds:
- url: URL
`,
			expected: Config{
				[]Feed{
					Feed{
						URL: "URL",
					},
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

func equal(a, b Config) bool {
	if len(a.Feeds) != len(b.Feeds) {
		return false
	}

	for i := 0; i < len(a.Feeds); i++ {
		if a.Feeds[i] != b.Feeds[i] {
			return false
		}
	}

	return true
}
