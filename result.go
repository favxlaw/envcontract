package envcontract

// FindingKind represents the category of a validation finding.
type FindingKind int

const (
	// KindMissing indicates a required env var is not present and has no default.
	KindMissing FindingKind = iota + 1
	// KindTypeMismatch indicates the env var value cannot be parsed as the expected type.
	KindTypeMismatch
	// KindUnused indicates an env var exists but no struct field references it.
	KindUnused
)

// Finding represents a single validation issue discovered by the engine.
type Finding struct {
	Kind    FindingKind
	EnvKey  string
	Message string
	IsError bool
}
