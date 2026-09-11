package envcontract

import (
	"strings"
	"testing"
)

func TestExportSchema(t *testing.T) {
	data, err := ExportSchema(&testConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(string(data), `"env_key": "PORT"`) {
		t.Fatalf("expected schema to mention PORT, got %s", data)
	}
}

func TestExportSchemaRejectsNonPointer(t *testing.T) {
	_, err := ExportSchema(testConfig{})
	if err == nil {
		t.Fatal("expected error for non-pointer input")
	}
}

func TestParseSchemaRoundTrip(t *testing.T) {
	data, err := ExportSchema(&testConfig{})
	if err != nil {
		t.Fatalf("export schema: %v", err)
	}

	contracts, err := ParseSchema(data)
	if err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	if len(contracts) == 0 {
		t.Fatal("expected at least one field contract")
	}

	var found bool
	for _, c := range contracts {
		if c.EnvKey == "PORT" && c.Required && c.Kind == "int" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected PORT contract to round-trip, got %+v", contracts)
	}
}

func TestParseSchemaRejectsInvalidJSON(t *testing.T) {
	_, err := ParseSchema([]byte("not json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
