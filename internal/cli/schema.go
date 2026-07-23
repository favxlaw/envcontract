package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

type schemaFlags struct {
	structName string
	format     string
}

var schemaArgs schemaFlags

var schemaCmd = &cobra.Command{
	Use:   "schema [flags] <file.go>",
	Short: "Export field contracts from a Go config struct",
	Long: `Parse a Go source file with env-tagged struct fields
and export the contract as JSON or YAML.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		structs, err := parseGoFile(args[0])
		if err != nil {
			return fmt.Errorf("parse struct: %w", err)
		}

		s := selectStruct(structs, schemaArgs.structName)
		if s == nil {
			if schemaArgs.structName != "" {
				return fmt.Errorf("struct %q not found in %s", schemaArgs.structName, args[0])
			}
			return fmt.Errorf("multiple structs found; use --struct to specify which one\n\n%s", structNames(structs))
		}

		contracts := toFieldContracts(s.Fields)

		switch schemaArgs.format {
		case "json":
			data, err := json.MarshalIndent(contracts, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal schema: %w", err)
			}
			fmt.Println(string(data))
		default:
			fmt.Printf("struct: %s\n", s.Name)
			fmt.Printf("fields: %d\n", len(s.Fields))
			for _, f := range s.Fields {
				def := ""
				if f.HasDefault {
					def = fmt.Sprintf(" default=%s", f.Default)
				}
				req := ""
				if f.Required {
					req = " required"
				}
				fmt.Printf("  %s -> %s (%s%s%s)\n", f.EnvKey, f.Name, f.Kind, req, def)
			}
		}

		return nil
	},
}

func init() {
	schemaCmd.Flags().StringVarP(&schemaArgs.structName, "struct", "s", "", "struct name to export")
	schemaCmd.Flags().StringVarP(&schemaArgs.format, "format", "o", "text", "output format: text or json")
}
