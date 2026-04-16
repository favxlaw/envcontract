package source

import (
	"bufio"
	"os"
	"strings"
)

type FileSource struct {
	Path string
}

func (f *FileSource) Load() (map[string]string, error) {
	file, err := os.Open(f.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		result[key] = value
	}

	return result, scanner.Err()
}
