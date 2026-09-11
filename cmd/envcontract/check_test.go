package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunCheckValid(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"check", "-schema", "testdata/schema.json", "-env", "testdata/valid.env"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d (stderr: %s)", code, stderr.String())
	}
}

func TestRunCheckMissingRequired(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"check", "-schema", "testdata/schema.json", "-env", "testdata/missing_required.env"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}

	if !strings.Contains(stdout.String(), "PORT") {
		t.Fatalf("expected output to mention PORT, got %q", stdout.String())
	}
}

func TestRunCheckTypeMismatch(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"check", "-schema", "testdata/schema.json", "-env", "testdata/type_mismatch.env"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestRunCheckRequiresSchemaFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"check", "-env", "testdata/valid.env"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}

	if !strings.Contains(stderr.String(), "-schema") {
		t.Fatalf("expected error to mention -schema, got %q", stderr.String())
	}
}

func TestRunCheckUnknownSchemaFile(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"check", "-schema", "testdata/does-not-exist.json"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}
