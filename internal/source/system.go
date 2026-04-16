package source

import (
	"os"
	"strings"
)

type SystemSource struct{}

func (s *SystemSource) Load() (map[string]string, error) {
	result := make(map[string]string)

	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			continue
		}
		result[parts[0]] = parts[1]
	}

	return result, nil
}
