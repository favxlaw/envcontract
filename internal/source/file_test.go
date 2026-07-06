package source

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFileSourceLoad(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		wantValues   map[string]string
		wantWarnings int
	}{
		{
			name: "valid env file",
			content: `
# app config
HOST=localhost
PORT=8080
DEBUG=true
DATABASE_URL=postgres://user:pass@localhost:5432/app
`,
			wantValues: map[string]string{
				"HOST":         "localhost",
				"PORT":         "8080",
				"DEBUG":        "true",
				"DATABASE_URL": "postgres://user:pass@localhost:5432/app",
			},
			wantWarnings: 0,
		},
		{
			name: "malformed lines warn and skip",
			content: `
HOST=localhost
BROKEN_LINE
=missing_key
PORT=8080
`,
			wantValues: map[string]string{
				"HOST": "localhost",
				"PORT": "8080",
			},
			wantWarnings: 2,
		},
		{
			name: "empty file",
			content: `
# only comment

`,
			wantValues:   map[string]string{},
			wantWarnings: 0,
		},
		{
			name: "duplicate keys last wins",
			content: `
HOST=first
HOST=second
`,
			wantValues: map[string]string{
				"HOST": "second",
			},
			wantWarnings: 0,
		},
		{
			name: "value may contain equals signs",
			content: `
TOKEN=abc=def=ghi
`,
			wantValues: map[string]string{
				"TOKEN": "abc=def=ghi",
			},
			wantWarnings: 0,
		},
		{
			name: "spaces around keys and values are trimmed",
			content: `
HOST = localhost
PORT = 8080
`,
			wantValues: map[string]string{
				"HOST": "localhost",
				"PORT": "8080",
			},
			wantWarnings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeTempFile(t, tt.content)

			got, err := FileSource{Path: path}.Load()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertMapEqual(t, got.Values, tt.wantValues)

			if len(got.Warnings) != tt.wantWarnings {
				t.Fatalf("expected %d warnings, got %d: %+v", tt.wantWarnings, len(got.Warnings), got.Warnings)
			}
		})
	}
}

func TestFileSourceMissingFile(t *testing.T) {
	_, err := FileSource{Path: "does-not-exist.env"}.Load()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected os.ErrNotExist, got %v", err)
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	return path
}

func assertMapEqual(t *testing.T, got, want map[string]string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("expected map length %d, got %d: %+v", len(want), len(got), got)
	}

	for key, wantValue := range want {
		if gotValue, ok := got[key]; !ok {
			t.Fatalf("expected key %q to exist", key)
		} else if gotValue != wantValue {
			t.Fatalf("expected %s=%q, got %q", key, wantValue, gotValue)
		}
	}
}
