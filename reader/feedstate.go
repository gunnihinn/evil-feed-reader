package reader

import (
	"encoding/json"
)

type FeedState struct {
	URI  string
	Hash string
	Seen bool
}

func Marshal(feeds []Feed) ([]byte, error) {
	data := make([]FeedState, 0)

	for _, feed := range feeds {
		data = append(data, FeedState{
			URI:  feed.Resource(),
			Hash: feed.Hash(),
			Seen: feed.Seen(),
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
		feeds[fs.URI] = fs
	}

	return feeds, nil
}
