package source

import (
	"errors"
	"testing"
)

func TestMockSourceLoad(t *testing.T) {
	src := MockSource{
		Values: map[string]string{
			"HOST": "localhost",
			"PORT": "8080",
		},
		Warnings: []LoadWarning{
			{Line: 2, Message: "example warning"},
		},
	}

	got, err := src.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertMapEqual(t, got.Values, src.Values)

	if len(got.Warnings) != 1 {
		t.Fatalf("expected one warning, got %+v", got.Warnings)
	}

	got.Values["HOST"] = "changed"

	if src.Values["HOST"] != "localhost" {
		t.Fatal("mock source values should not be mutated by caller")
	}
}

func TestMockSourceError(t *testing.T) {
	wantErr := errors.New("boom")

	_, err := MockSource{Err: wantErr}.Load()
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}
