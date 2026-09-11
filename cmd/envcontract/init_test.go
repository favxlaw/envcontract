package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInitWritesFile(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), ".env.example")
	var stdout, stderr bytes.Buffer

	code := run([]string{"init", "-schema", "testdata/schema.json", "-out", outPath}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d (stderr: %s)", code, stderr.String())
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("expected output file to exist: %v", err)
	}

	if !strings.Contains(string(data), "PORT is required") {
		t.Fatalf("expected generated file to mention PORT, got %q", data)
	}
}

func TestRunInitRequiresSchemaFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"init"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
}

func TestRunInitRefusesToOverwrite(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), ".env.example")
	if err := os.WriteFile(outPath, []byte("existing"), 0o644); err != nil {
		t.Fatalf("seed existing file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"init", "-schema", "testdata/schema.json", "-out", outPath}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}
	if string(data) != "existing" {
		t.Fatalf("expected existing file to be untouched, got %q", data)
	}
}

func TestRunInitForceOverwrites(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), ".env.example")
	if err := os.WriteFile(outPath, []byte("existing"), 0o644); err != nil {
		t.Fatalf("seed existing file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"init", "-schema", "testdata/schema.json", "-out", outPath, "-force"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d (stderr: %s)", code, stderr.String())
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}
	if strings.Contains(string(data), "existing") {
		t.Fatalf("expected file to be overwritten, got %q", data)
	}
}
