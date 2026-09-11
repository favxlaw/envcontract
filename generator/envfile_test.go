package generator_test

import (
	"testing"

	"github.com/favxlaw/envcontract"
	"github.com/favxlaw/envcontract/generator"
)

func TestEnvExample(t *testing.T) {
	contracts := []envcontract.FieldContract{
		{EnvKey: "PORT", Required: true, Kind: "int"},
		{EnvKey: "HOST", HasDefault: true, Default: "localhost", Kind: "string"},
		{EnvKey: "DEBUG", Kind: "bool"},
	}

	want := "# DEBUG is optional\n" +
		"DEBUG=\n" +
		"\n" +
		"# HOST is optional (default: localhost)\n" +
		"HOST=localhost\n" +
		"\n" +
		"# PORT is required\n" +
		"PORT=\n"

	got := generator.EnvExample(contracts)
	if got != want {
		t.Fatalf("unexpected output:\n got: %q\nwant: %q", got, want)
	}
}

func TestEnvExampleEmpty(t *testing.T) {
	if got := generator.EnvExample(nil); got != "" {
		t.Fatalf("expected empty output, got %q", got)
	}
}
