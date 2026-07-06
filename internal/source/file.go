package source

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type FileSource struct {
	Path string
}

func (f FileSource) Load() (result LoadResult, err error) {
	file, err := os.Open(f.Path)
	if err != nil {
		return LoadResult{}, fmt.Errorf("open env file %q: %w", f.Path, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close env file %q: %w", f.Path, closeErr)
		}
	}()

	result = LoadResult{
		Values: make(map[string]string),
	}

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++

		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			result.Warnings = append(result.Warnings, LoadWarning{
				Line:    lineNumber,
				Message: "malformed line skipped: missing =",
			})
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		if key == "" {
			result.Warnings = append(result.Warnings, LoadWarning{
				Line:    lineNumber,
				Message: "malformed line skipped: empty key",
			})
			continue
		}

		result.Values[key] = value
	}

	if err := scanner.Err(); err != nil {
		return LoadResult{}, fmt.Errorf("read env file %q: %w", f.Path, err)
	}

	return result, nil
}
