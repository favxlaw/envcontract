package parser

import (
	"testing"

	"github.com/favxlaw/envcontract"
)

func TestParseStruct(t *testing.T) {
	type Simple struct {
		Host string `env:"HOST"`
	}

	type AllTypes struct {
		S  string  `env:"S"`
		I  int     `env:"I"`
		I6 int64   `env:"I6"`
		F  float64 `env:"F"`
		B  bool    `env:"B"`
	}

	type PointerScalars struct {
		Host *string `env:"HOST"`
		Port *int    `env:"PORT,required"`
	}

	type Inner struct {
		Port int `env:"PORT"`
	}

	type WithNested struct {
		Inner Inner
	}

	type WithPtrNested struct {
		Inner *Inner
	}

	type WithUntagged struct {
		Host    string `env:"HOST"`
		Ignored string
	}

	type WithRequired struct {
		Host string `env:"HOST,required"`
	}

	type WithDefault struct {
		Host string `env:"HOST,default=localhost"`
	}

	type WithSkipTag struct {
		Host   string `env:"HOST"`
		Secret string `env:"-"`
	}

	tests := []struct {
		name      string
		input     any
		wantErr   bool
		wantLen   int
		wantFirst *envcontract.FieldContract
	}{
		{
			name:    "nil input returns error",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "non-pointer input returns error",
			input:   Simple{},
			wantErr: true,
		},
		{
			name:    "nil pointer returns error",
			input:   (*Simple)(nil),
			wantErr: true,
		},
		{
			name:    "pointer to non-struct returns error",
			input:   func() any { s := "hello"; return &s }(),
			wantErr: true,
		},
		{
			name:    "valid struct with one tagged field",
			input:   &Simple{},
			wantLen: 1,
			wantFirst: &envcontract.FieldContract{
				Name:   "Host",
				EnvKey: "HOST",
				Kind:   "string",
			},
		},
		{
			name:    "untagged field is skipped",
			input:   &WithUntagged{},
			wantLen: 1,
			wantFirst: &envcontract.FieldContract{
				Name:   "Host",
				EnvKey: "HOST",
				Kind:   "string",
			},
		},
		{
			name:    "skip tag is skipped",
			input:   &WithSkipTag{},
			wantLen: 1,
			wantFirst: &envcontract.FieldContract{
				Name:   "Host",
				EnvKey: "HOST",
				Kind:   "string",
			},
		},
		{
			name:    "required tag sets Required true",
			input:   &WithRequired{},
			wantLen: 1,
			wantFirst: &envcontract.FieldContract{
				Name:     "Host",
				EnvKey:   "HOST",
				Required: true,
				Kind:     "string",
			},
		},
		{
			name:    "default tag sets HasDefault and Default",
			input:   &WithDefault{},
			wantLen: 1,
			wantFirst: &envcontract.FieldContract{
				Name:       "Host",
				EnvKey:     "HOST",
				HasDefault: true,
				Default:    "localhost",
				Kind:       "string",
			},
		},
		{
			name:    "all supported types",
			input:   &AllTypes{},
			wantLen: 5,
		},
		{
			name:    "nil pointer scalar fields are still parsed",
			input:   &PointerScalars{},
			wantLen: 2,
			wantFirst: &envcontract.FieldContract{
				Name:   "Host",
				EnvKey: "HOST",
				Kind:   "string",
			},
		},
		{
			name:    "nested struct fields are flattened",
			input:   &WithNested{},
			wantLen: 1,
			wantFirst: &envcontract.FieldContract{
				Name:   "Port",
				EnvKey: "PORT",
				Kind:   "int",
			},
		},
		{
			name:    "nil pointer to nested struct is still parsed",
			input:   &WithPtrNested{},
			wantLen: 1,
			wantFirst: &envcontract.FieldContract{
				Name:   "Port",
				EnvKey: "PORT",
				Kind:   "int",
			},
		},
		{
			name:    "pointer to nested struct is parsed",
			input:   &WithPtrNested{Inner: &Inner{}},
			wantLen: 1,
			wantFirst: &envcontract.FieldContract{
				Name:   "Port",
				EnvKey: "PORT",
				Kind:   "int",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStruct(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != tt.wantLen {
				t.Fatalf("expected %d contracts, got %d: %+v", tt.wantLen, len(got), got)
			}

			if tt.wantFirst != nil && got[0] != *tt.wantFirst {
				t.Errorf("first contract mismatch\ngot:  %+v\nwant: %+v", got[0], *tt.wantFirst)
			}
		})
	}
}

func TestParseStructUnknownTagOption(t *testing.T) {
	type Config struct {
		Host string `env:"HOST,banana"`
	}

	_, err := ParseStruct(&Config{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
