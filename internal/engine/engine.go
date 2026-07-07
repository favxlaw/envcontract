package engine

import (
	"fmt"
	"strconv"
	"time"

	"github.com/favxlaw/envcontract/internal/contract"
)

// Options controls optional engine behavior.
type Options struct {
	CheckUnused bool
}

// Run executes all checks and returns the combined findings.
func Run(contracts []contract.FieldContract, env map[string]string, opts Options) []contract.Finding {
	var findings []contract.Finding

	findings = append(findings, checkDefaults(contracts)...)
	findings = append(findings, checkMissing(contracts, env)...)
	findings = append(findings, checkTypes(contracts, env)...)

	if opts.CheckUnused {
		findings = append(findings, checkUnused(contracts, env)...)
	}

	return findings
}

func checkDefaults(contracts []contract.FieldContract) []contract.Finding {
	var findings []contract.Finding

	for _, c := range contracts {
		if !c.HasDefault {
			continue
		}

		if !canParse(c.Kind, c.Default) {
			findings = append(findings, contract.Finding{
				Kind:    contract.KindInvalidDefault,
				EnvKey:  c.EnvKey,
				Message: fmt.Sprintf("default for env var %s: expected %s, got %q", c.EnvKey, c.Kind, c.Default),
				IsError: true,
			})
		}
	}

	return findings
}

func checkMissing(contracts []contract.FieldContract, env map[string]string) []contract.Finding {
	var findings []contract.Finding

	for _, c := range contracts {
		if _, exists := env[c.EnvKey]; exists {
			continue
		}

		if c.HasDefault {
			continue
		}

		finding := contract.Finding{
			Kind:   contract.KindMissing,
			EnvKey: c.EnvKey,
		}

		if c.Required {
			finding.IsError = true
			finding.Message = fmt.Sprintf("required env var %s is missing and has no default", c.EnvKey)
		} else {
			finding.IsError = false
			finding.Message = fmt.Sprintf("optional env var %s is not set", c.EnvKey)
		}

		findings = append(findings, finding)
	}

	return findings
}

func checkTypes(contracts []contract.FieldContract, env map[string]string) []contract.Finding {
	var findings []contract.Finding

	for _, c := range contracts {
		value, exists := env[c.EnvKey]
		if !exists {
			continue
		}

		if canParse(c.Kind, value) {
			continue
		}

		findings = append(findings, contract.Finding{
			Kind:    contract.KindTypeMismatch,
			EnvKey:  c.EnvKey,
			Message: fmt.Sprintf("env var %s: expected %s, got %q", c.EnvKey, c.Kind, value),
			IsError: true,
		})
	}

	return findings
}

func checkUnused(contracts []contract.FieldContract, env map[string]string) []contract.Finding {
	expected := make(map[string]bool, len(contracts))

	for _, c := range contracts {
		expected[c.EnvKey] = true
	}

	var findings []contract.Finding

	for key := range env {
		if expected[key] {
			continue
		}

		findings = append(findings, contract.Finding{
			Kind:    contract.KindUnused,
			EnvKey:  key,
			Message: fmt.Sprintf("env var %s is set but not referenced by any config field", key),
			IsError: false,
		})
	}

	return findings
}

func canParse(kind, value string) bool {
	switch kind {
	case "string":
		return true
	case "int":
		_, err := strconv.Atoi(value)
		return err == nil
	case "int64":
		_, err := strconv.ParseInt(value, 10, 64)
		return err == nil
	case "float64":
		_, err := strconv.ParseFloat(value, 64)
		return err == nil
	case "bool":
		_, err := strconv.ParseBool(value)
		return err == nil
	case "duration":
		_, err := time.ParseDuration(value)
		return err == nil
	default:
		return false
	}
}
