package engine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/favxlaw/envcontract"
)

// Options controls optional engine behaviour.
type Options struct {
	// CheckUnused enables reporting env vars that no struct field references.
	CheckUnused bool
}

// Run executes all checks and returns the combined findings.
func Run(contracts []envcontract.FieldContract, env map[string]string, opts Options) []envcontract.Finding {
	var findings []envcontract.Finding
	findings = append(findings, checkMissing(contracts, env)...)
	findings = append(findings, checkTypes(contracts, env)...)
	if opts.CheckUnused {
		findings = append(findings, checkUnused(contracts, env)...)
	}
	return findings
}

// checkMissing reports env vars that the struct expects but the environment
// does not provide. The rules:
//   - If the env key exists in the map (even as ""), the field is present → no finding.
//   - If the field has a default (HasDefault), a missing key is fine → no finding.
//   - If the field is required with no default → error finding.
//   - If the field is optional (required=false) with no default → warning finding.
func checkMissing(contracts []envcontract.FieldContract, env map[string]string) []envcontract.Finding {
	var findings []envcontract.Finding
	for _, c := range contracts {
		if _, exists := env[c.EnvKey]; exists {
			continue
		}
		if c.HasDefault {
			continue
		}
		f := envcontract.Finding{
			Kind:   envcontract.KindMissing,
			EnvKey: c.EnvKey,
		}
		if c.Required {
			f.IsError = true
			f.Message = fmt.Sprintf("required env var %s is missing and has no default", c.EnvKey)
		} else {
			f.IsError = false
			f.Message = fmt.Sprintf("optional env var %s is not set", c.EnvKey)
		}
		findings = append(findings, f)
	}
	return findings
}

// checkTypes verifies that every present env value can be parsed as the
// expected Go type. Only values that actually exist in env are checked —
// missing keys are handled by checkMissing.
func checkTypes(contracts []envcontract.FieldContract, env map[string]string) []envcontract.Finding {
	var findings []envcontract.Finding
	for _, c := range contracts {
		val, exists := env[c.EnvKey]
		if !exists {
			continue
		}
		if !canParse(c.Kind, val) {
			findings = append(findings, envcontract.Finding{
				Kind:    envcontract.KindTypeMismatch,
				EnvKey:  c.EnvKey,
				Message: fmt.Sprintf("env var %s: expected %s, got %q", c.EnvKey, c.Kind, val),
				IsError: true,
			})
		}
	}
	return findings
}

// checkUnused reports env keys that exist in the map but are not referenced
// by any struct field. This is disabled by default.
func checkUnused(contracts []envcontract.FieldContract, env map[string]string) []envcontract.Finding {
	expected := make(map[string]bool, len(contracts))
	for _, c := range contracts {
		expected[c.EnvKey] = true
	}
	var findings []envcontract.Finding
	for key := range env {
		if !expected[key] {
			findings = append(findings, envcontract.Finding{
				Kind:    envcontract.KindUnused,
				EnvKey:  key,
				Message: fmt.Sprintf("env var %s is set but not referenced by any config field", key),
				IsError: false,
			})
		}
	}
	return findings
}

// canParse checks whether val can be parsed as the given kind string.
func canParse(kind, val string) bool {
	switch kind {
	case "string":
		return true
	case "int":
		_, err := strconv.Atoi(val)
		return err == nil
	case "int64":
		_, err := strconv.ParseInt(val, 10, 64)
		return err == nil
	case "float64":
		_, err := strconv.ParseFloat(val, 64)
		return err == nil
	case "bool":
		lower := strings.ToLower(val)
		return lower == "true" || lower == "false" || lower == "1" || lower == "0"
	}
	return false
}
