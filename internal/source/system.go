package source

import (
	"os"
	"strings"
)

type SystemSource struct{}

// If duplicate keys somehow appear, the last value wins.
func (SystemSource) Load() (LoadResult, error) {
	result := LoadResult{
		Values: make(map[string]string),
	}

	for _, entry := range os.Environ() {
		key, value, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}

		result.Values[key] = value
	}

	return result, nil
}
