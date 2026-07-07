package contract

type FieldContract struct {
	Name       string `json:"name"`
	EnvKey     string `json:"env_key"`
	Required   bool   `json:"required"`
	HasDefault bool   `json:"has_default"`
	Default    string `json:"default,omitempty"`
	Kind       string `json:"kind"`
}

// FindingKind represents the category of a validation finding.
type FindingKind string

const (
	KindMissing        FindingKind = "missing"
	KindTypeMismatch   FindingKind = "type_mismatch"
	KindUnused         FindingKind = "unused"
	KindInvalidDefault FindingKind = "invalid_default"
)

// Finding represents a single validation issue discovered by envcontract.
type Finding struct {
	Kind    FindingKind `json:"kind"`
	EnvKey  string      `json:"env_key"`
	Message string      `json:"message"`
	IsError bool        `json:"is_error"`
}
