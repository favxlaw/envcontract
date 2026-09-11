// Package generator turns field contracts into scaffold files such as a
// starter .env.example.
package generator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/favxlaw/envcontract"
)

// EnvExample renders contracts as .env.example content: one commented
// annotation line per field (required/optional, and any default) followed
// by a KEY=value line, sorted by env key for stable output.
func EnvExample(contracts []envcontract.FieldContract) string {
	sorted := append([]envcontract.FieldContract(nil), contracts...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].EnvKey < sorted[j].EnvKey
	})

	var b strings.Builder

	for i, c := range sorted {
		if i > 0 {
			b.WriteString("\n")
		}

		switch {
		case c.Required:
			fmt.Fprintf(&b, "# %s is required\n", c.EnvKey)
		case c.HasDefault:
			fmt.Fprintf(&b, "# %s is optional (default: %s)\n", c.EnvKey, c.Default)
		default:
			fmt.Fprintf(&b, "# %s is optional\n", c.EnvKey)
		}

		fmt.Fprintf(&b, "%s=%s\n", c.EnvKey, c.Default)
	}

	return b.String()
}
