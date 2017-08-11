package reader

import (
	"encoding/json"
)

type FeedState struct {
	Title    string
	URL      string
	Resource string
	Hash     string
	Seen     bool
}

func Marshal(feeds []Feed) ([]byte, error) {
	data := make([]FeedState, 0)

	for _, feed := range feeds {
		data = append(data, FeedState{
			Title:    feed.Title(),
			URL:      feed.Url(),
			Resource: feed.Resource(),
			Hash:     feed.Hash(),
			Seen:     feed.Seen(),
		})
	}

	return json.Marshal(data)
}

func Unmarshal(blob []byte) (map[string]FeedState, error) {
	data := make([]FeedState, 0)
	err := json.Unmarshal(blob, &data)

	if err != nil {
		return nil, err
	}

	feeds := make(map[string]FeedState)
	for _, fs := range data {
		feeds[fs.Resource] = fs
	}

	return feeds, nil
}
