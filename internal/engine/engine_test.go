package engine

import (
	"testing"

	"github.com/favxlaw/envcontract"
)

func TestCheckMissing(t *testing.T) {
	tests := []struct {
		name      string
		contracts []envcontract.FieldContract
		env       map[string]string
		wantLen   int
		wantError []bool // IsError for each expected finding, in order
	}{
		{
			name: "present key produces no finding",
			contracts: []envcontract.FieldContract{
				{EnvKey: "HOST", Required: true, Kind: "string"},
			},
			env:     map[string]string{"HOST": "localhost"},
			wantLen: 0,
		},
		{
			name: "empty string counts as present",
			contracts: []envcontract.FieldContract{
				{EnvKey: "HOST", Required: true, Kind: "string"},
			},
			env:     map[string]string{"HOST": ""},
			wantLen: 0,
		},
		{
			name: "missing required with no default is error",
			contracts: []envcontract.FieldContract{
				{EnvKey: "DB_URL", Required: true, Kind: "string"},
			},
			env:       map[string]string{},
			wantLen:   1,
			wantError: []bool{true},
		},
		{
			name: "missing optional with no default is warning",
			contracts: []envcontract.FieldContract{
				{EnvKey: "DEBUG", Required: false, Kind: "bool"},
			},
			env:       map[string]string{},
			wantLen:   1,
			wantError: []bool{false},
		},
		{
			name: "missing key with default produces no finding",
			contracts: []envcontract.FieldContract{
				{EnvKey: "PORT", Required: true, HasDefault: true, Default: "8080", Kind: "int"},
			},
			env:     map[string]string{},
			wantLen: 0,
		},
		{
			name: "multiple contracts mixed",
			contracts: []envcontract.FieldContract{
				{EnvKey: "A", Required: true, Kind: "string"},
				{EnvKey: "B", Required: false, Kind: "string"},
				{EnvKey: "C", Required: true, HasDefault: true, Default: "x", Kind: "string"},
			},
			env:       map[string]string{},
			wantLen:   2, // A is error, B is warning, C has default so no finding
			wantError: []bool{true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkMissing(tt.contracts, tt.env)
			if len(got) != tt.wantLen {
				t.Fatalf("expected %d findings, got %d: %+v", tt.wantLen, len(got), got)
			}
			for i, wantErr := range tt.wantError {
				if got[i].IsError != wantErr {
					t.Errorf("finding[%d].IsError = %v, want %v", i, got[i].IsError, wantErr)
				}
				if got[i].Kind != envcontract.KindMissing {
					t.Errorf("finding[%d].Kind = %v, want KindMissing", i, got[i].Kind)
				}
			}
		})
	}
}

func TestCheckTypes(t *testing.T) {
	tests := []struct {
		name    string
		contracts []envcontract.FieldContract
		env     map[string]string
		wantLen int
	}{
		{
			name: "valid string always passes",
			contracts: []envcontract.FieldContract{
				{EnvKey: "HOST", Kind: "string"},
			},
			env:     map[string]string{"HOST": "anything"},
			wantLen: 0,
		},
		{
			name: "valid int passes",
			contracts: []envcontract.FieldContract{
				{EnvKey: "PORT", Kind: "int"},
			},
			env:     map[string]string{"PORT": "8080"},
			wantLen: 0,
		},
		{
			name: "invalid int is error",
			contracts: []envcontract.FieldContract{
				{EnvKey: "PORT", Kind: "int"},
			},
			env:     map[string]string{"PORT": "abc"},
			wantLen: 1,
		},
		{
			name: "valid int64 passes",
			contracts: []envcontract.FieldContract{
				{EnvKey: "BIG", Kind: "int64"},
			},
			env:     map[string]string{"BIG": "9999999999"},
			wantLen: 0,
		},
		{
			name: "invalid int64 is error",
			contracts: []envcontract.FieldContract{
				{EnvKey: "BIG", Kind: "int64"},
			},
			env:     map[string]string{"BIG": "not_a_number"},
			wantLen: 1,
		},
		{
			name: "valid float64 passes",
			contracts: []envcontract.FieldContract{
				{EnvKey: "RATE", Kind: "float64"},
			},
			env:     map[string]string{"RATE": "3.14"},
			wantLen: 0,
		},
		{
			name: "invalid float64 is error",
			contracts: []envcontract.FieldContract{
				{EnvKey: "RATE", Kind: "float64"},
			},
			env:     map[string]string{"RATE": "abc"},
			wantLen: 1,
		},
		{
			name: "valid bool passes (true)",
			contracts: []envcontract.FieldContract{
				{EnvKey: "DEBUG", Kind: "bool"},
			},
			env:     map[string]string{"DEBUG": "true"},
			wantLen: 0,
		},
		{
			name: "valid bool passes (0)",
			contracts: []envcontract.FieldContract{
				{EnvKey: "DEBUG", Kind: "bool"},
			},
			env:     map[string]string{"DEBUG": "0"},
			wantLen: 0,
		},
		{
			name: "invalid bool is error",
			contracts: []envcontract.FieldContract{
				{EnvKey: "DEBUG", Kind: "bool"},
			},
			env:     map[string]string{"DEBUG": "yes"},
			wantLen: 1,
		},
		{
			name: "missing key is skipped (not a type error)",
			contracts: []envcontract.FieldContract{
				{EnvKey: "MISSING", Kind: "int"},
			},
			env:     map[string]string{},
			wantLen: 0,
		},
		{
			name: "empty string for int is type error",
			contracts: []envcontract.FieldContract{
				{EnvKey: "PORT", Kind: "int"},
			},
			env:     map[string]string{"PORT": ""},
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkTypes(tt.contracts, tt.env)
			if len(got) != tt.wantLen {
				t.Fatalf("expected %d findings, got %d: %+v", tt.wantLen, len(got), got)
			}
			for _, f := range got {
				if f.Kind != envcontract.KindTypeMismatch {
					t.Errorf("expected KindTypeMismatch, got %v", f.Kind)
				}
				if !f.IsError {
					t.Error("type mismatch findings should always be errors")
				}
			}
		})
	}
}

func TestCheckUnused(t *testing.T) {
	tests := []struct {
		name      string
		contracts []envcontract.FieldContract
		env       map[string]string
		wantLen   int
		wantKeys  []string
	}{
		{
			name: "no unused vars",
			contracts: []envcontract.FieldContract{
				{EnvKey: "HOST", Kind: "string"},
			},
			env:     map[string]string{"HOST": "localhost"},
			wantLen: 0,
		},
		{
			name: "one unused var",
			contracts: []envcontract.FieldContract{
				{EnvKey: "HOST", Kind: "string"},
			},
			env:      map[string]string{"HOST": "localhost", "STALE_KEY": "leftover"},
			wantLen:  1,
			wantKeys: []string{"STALE_KEY"},
		},
		{
			name:      "all unused when no contracts",
			contracts: []envcontract.FieldContract{},
			env:       map[string]string{"A": "1", "B": "2"},
			wantLen:   2,
		},
		{
			name: "empty env produces no findings",
			contracts: []envcontract.FieldContract{
				{EnvKey: "HOST", Kind: "string"},
			},
			env:     map[string]string{},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkUnused(tt.contracts, tt.env)
			if len(got) != tt.wantLen {
				t.Fatalf("expected %d findings, got %d: %+v", tt.wantLen, len(got), got)
			}
			for _, f := range got {
				if f.Kind != envcontract.KindUnused {
					t.Errorf("expected KindUnused, got %v", f.Kind)
				}
				if f.IsError {
					t.Error("unused findings should be warnings, not errors")
				}
			}
			for i, wantKey := range tt.wantKeys {
				if got[i].EnvKey != wantKey {
					t.Errorf("finding[%d].EnvKey = %q, want %q", i, got[i].EnvKey, wantKey)
				}
			}
		})
	}
}

func TestRunIntegration(t *testing.T) {
	tests := []struct {
		name      string
		contracts []envcontract.FieldContract
		env       map[string]string
		opts      Options
		wantLen   int
	}{
		{
			name: "all good — no findings",
			contracts: []envcontract.FieldContract{
				{EnvKey: "HOST", Required: true, Kind: "string"},
				{EnvKey: "PORT", Required: true, Kind: "int"},
			},
			env:     map[string]string{"HOST": "localhost", "PORT": "8080"},
			wantLen: 0,
		},
		{
			name: "missing + type error combined",
			contracts: []envcontract.FieldContract{
				{EnvKey: "HOST", Required: true, Kind: "string"},
				{EnvKey: "PORT", Required: true, Kind: "int"},
			},
			env:     map[string]string{"PORT": "abc"},
			wantLen: 2, // HOST missing + PORT type mismatch
		},
		{
			name: "unused check disabled by default",
			contracts: []envcontract.FieldContract{
				{EnvKey: "HOST", Required: true, Kind: "string"},
			},
			env:     map[string]string{"HOST": "localhost", "EXTRA": "ignored"},
			opts:    Options{CheckUnused: false},
			wantLen: 0,
		},
		{
			name: "unused check enabled",
			contracts: []envcontract.FieldContract{
				{EnvKey: "HOST", Required: true, Kind: "string"},
			},
			env:     map[string]string{"HOST": "localhost", "EXTRA": "found"},
			opts:    Options{CheckUnused: true},
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Run(tt.contracts, tt.env, tt.opts)
			if len(got) != tt.wantLen {
				t.Fatalf("expected %d findings, got %d: %+v", tt.wantLen, len(got), got)
			}
		})
	}
}
