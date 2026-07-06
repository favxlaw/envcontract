package source

// MockSource is an in-memory Source useful for tests.
type MockSource struct {
	Values   map[string]string
	Warnings []LoadWarning
	Err      error
}

// Load returns the configured mock values, warnings, and error.
func (m MockSource) Load() (LoadResult, error) {
	if m.Err != nil {
		return LoadResult{}, m.Err
	}

	values := make(map[string]string, len(m.Values))
	for key, value := range m.Values {
		values[key] = value
	}

	warnings := make([]LoadWarning, len(m.Warnings))
	copy(warnings, m.Warnings)

	return LoadResult{
		Values:   values,
		Warnings: warnings,
	}, nil
}
