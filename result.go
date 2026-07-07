package envcontract

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/favxlaw/envcontract/internal/contract"
)

type FindingKind = contract.FindingKind

const (
	KindMissing = contract.KindMissing

	// KindTypeMismatch indicates an env var value cannot be parsed as the expected type.
	KindTypeMismatch = contract.KindTypeMismatch

	// KindUnused indicates an env var exists but no struct field references it.
	KindUnused = contract.KindUnused

	// KindInvalidDefault indicates a default value cannot be parsed as the expected type.
	KindInvalidDefault = contract.KindInvalidDefault
)

type Finding = contract.Finding

type Result struct {
	findings       []Finding
	sourceWarnings []LoadWarning
}

func newResult(findings []Finding, sourceWarnings []LoadWarning) Result {
	return Result{
		findings:       append([]Finding(nil), findings...),
		sourceWarnings: append([]LoadWarning(nil), sourceWarnings...),
	}
}

// Findings returns all validation findings.
func (r Result) Findings() []Finding {
	return append([]Finding(nil), r.findings...)
}

// SourceWarnings returns non-fatal warnings produced while loading sources.
func (r Result) SourceWarnings() []LoadWarning {
	return append([]LoadWarning(nil), r.sourceWarnings...)
}

// HasErrors reports whether any validation finding is an error.
func (r Result) HasErrors() bool {
	for _, finding := range r.findings {
		if finding.IsError {
			return true
		}
	}

	return false
}

// HasWarnings reports whether there are validation warnings or source warnings.
func (r Result) HasWarnings() bool {
	if len(r.sourceWarnings) > 0 {
		return true
	}

	for _, finding := range r.findings {
		if !finding.IsError {
			return true
		}
	}

	return false
}

// Errors returns validation findings marked as errors.
func (r Result) Errors() []Finding {
	var errors []Finding

	for _, finding := range r.findings {
		if finding.IsError {
			errors = append(errors, finding)
		}
	}

	return errors
}

// Warnings returns validation findings marked as warnings.
//
// Source-loading warnings are returned separately by SourceWarnings.
func (r Result) Warnings() []Finding {
	var warnings []Finding

	for _, finding := range r.findings {
		if !finding.IsError {
			warnings = append(warnings, finding)
		}
	}

	return warnings
}

// Print writes a human-readable validation report.
func (r Result) Print(w io.Writer) error {
	if len(r.findings) == 0 && len(r.sourceWarnings) == 0 {
		_, err := fmt.Fprintln(w, "envcontract: ok")
		return err
	}

	for _, warning := range r.sourceWarnings {
		if warning.Line > 0 {
			if _, err := fmt.Fprintf(w, "warning: line %d: %s\n", warning.Line, warning.Message); err != nil {
				return err
			}
			continue
		}

		if _, err := fmt.Fprintf(w, "warning: %s\n", warning.Message); err != nil {
			return err
		}
	}

	for _, finding := range r.findings {
		level := "warning"
		if finding.IsError {
			level = "error"
		}

		if _, err := fmt.Fprintf(w, "%s: %s: %s\n", level, finding.EnvKey, finding.Message); err != nil {
			return err
		}
	}

	return nil
}

// JSON returns a structured JSON representation of the result.
func (r Result) JSON() ([]byte, error) {
	output := struct {
		HasErrors      bool          `json:"has_errors"`
		HasWarnings    bool          `json:"has_warnings"`
		Findings       []Finding     `json:"findings"`
		SourceWarnings []LoadWarning `json:"source_warnings"`
	}{
		HasErrors:      r.HasErrors(),
		HasWarnings:    r.HasWarnings(),
		Findings:       r.Findings(),
		SourceWarnings: r.SourceWarnings(),
	}

	return json.MarshalIndent(output, "", "  ")
}
