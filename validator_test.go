package envcontract

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

type testConfig struct {
	Host    string        `env:"HOST,default=localhost"`
	Port    int           `env:"PORT,required"`
	Debug   bool          `env:"DEBUG"`
	Rate    float64       `env:"RATE"`
	Timeout time.Duration `env:"TIMEOUT,default=30s"`
}

func TestValidateWithFile(t *testing.T) {
	result, err := Validate(&testConfig{}, WithFile("testdata/valid.env"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.HasErrors() {
		t.Fatalf("expected no errors, got %+v", result.Errors())
	}
}

func TestValidateReportsMissingRequired(t *testing.T) {
	result, err := Validate(&testConfig{}, WithFile("testdata/missing_required.env"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.HasErrors() {
		t.Fatal("expected validation errors")
	}

	errors := result.Errors()
	if len(errors) != 1 {
		t.Fatalf("expected one error, got %+v", errors)
	}

	if errors[0].Kind != KindMissing {
		t.Fatalf("expected KindMissing, got %s", errors[0].Kind)
	}
}

func TestValidateReportsTypeMismatch(t *testing.T) {
	result, err := Validate(&testConfig{}, WithFile("testdata/type_mismatch.env"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.HasErrors() {
		t.Fatal("expected validation errors")
	}

	if len(result.Errors()) == 0 {
		t.Fatal("expected at least one error")
	}
}

func TestValidateReportsUnusedWhenEnabled(t *testing.T) {
	src := mockSource{
		values: map[string]string{
			"HOST":    "localhost",
			"PORT":    "8080",
			"DEBUG":   "false",
			"RATE":    "1.5",
			"TIMEOUT": "30s",
			"EXTRA":   "unused",
		},
	}

	result, err := Validate(&testConfig{}, WithSource(src), WithUnusedCheck())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	warnings := result.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("expected one warning, got %+v", warnings)
	}

	if warnings[0].Kind != KindUnused {
		t.Fatalf("expected KindUnused, got %s", warnings[0].Kind)
	}
}

func TestResultPrint(t *testing.T) {
	src := mockSource{
		values: map[string]string{
			"PORT": "abc",
		},
	}

	result, err := Validate(&testConfig{}, WithSource(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	if err := result.Print(&buf); err != nil {
		t.Fatalf("print result: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "error: PORT") {
		t.Fatalf("expected output to contain PORT error, got %q", output)
	}
}

func TestResultJSON(t *testing.T) {
	src := mockSource{
		values: map[string]string{
			"PORT": "abc",
		},
	}

	result, err := Validate(&testConfig{}, WithSource(src))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := result.JSON()
	if err != nil {
		t.Fatalf("json result: %v", err)
	}

	output := string(data)
	if !strings.Contains(output, `"has_errors": true`) {
		t.Fatalf("expected JSON to report errors, got %s", output)
	}
}

type mockSource struct {
	values map[string]string
}

func (m mockSource) Load() (LoadResult, error) {
	values := make(map[string]string, len(m.values))
	for key, value := range m.values {
		values[key] = value
	}

	return LoadResult{
		Values: values,
	}, nil
}
