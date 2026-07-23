package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/favxlaw/envcontract"
	"github.com/favxlaw/envcontract/internal/contract"
	"github.com/favxlaw/envcontract/internal/engine"
	"github.com/favxlaw/envcontract/internal/source"
	"github.com/spf13/cobra"
)

type checkFlags struct {
	envFiles   []string
	structName string
	unused     bool
	system     bool
	format     string
}

var checkArgs checkFlags

var checkCmd = &cobra.Command{
	Use:   "check [flags] <file.go>",
	Short: "Validate env vars against a Go config struct",
	Long: `Parse a Go source file containing config structs with env tags,
load environment variables from .env files or the system,
and report missing, mismatched, or unused variables.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		structs, err := parseGoFile(args[0])
		if err != nil {
			return fmt.Errorf("parse struct: %w", err)
		}

		s := selectStruct(structs, checkArgs.structName)
		if s == nil {
			if checkArgs.structName != "" {
				return fmt.Errorf("struct %q not found in %s", checkArgs.structName, args[0])
			}
			return fmt.Errorf("multiple structs found; use --struct to specify which one\n\n%s", structNames(structs))
		}

		contracts := toFieldContracts(s.Fields)

		env, loadWarnings, err := loadSources(
			checkArgs.envFiles,
			checkArgs.system,
		)
		if err != nil {
			return err
		}

		findings := engine.Run(contracts, env, engine.Options{
			CheckUnused: checkArgs.unused,
		})

		result := envcontract.NewResult(findings, loadWarnings)

		switch checkArgs.format {
		case "json":
			data, err := result.JSON()
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		default:
			if err := result.Print(os.Stdout); err != nil {
				return err
			}
		}

		if result.HasErrors() {
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	checkCmd.Flags().StringArrayVarP(&checkArgs.envFiles, "file", "f", nil, "path to .env file (can be specified multiple times)")
	checkCmd.Flags().StringVarP(&checkArgs.structName, "struct", "s", "", "struct name to validate (required if file has multiple structs)")
	checkCmd.Flags().BoolVarP(&checkArgs.unused, "unused", "u", false, "report env vars not referenced by any struct field")
	checkCmd.Flags().BoolVarP(&checkArgs.system, "system", "e", false, "include system environment variables")
	checkCmd.Flags().StringVarP(&checkArgs.format, "format", "o", "text", "output format: text or json")
}

func selectStruct(structs []sourceStruct, name string) *sourceStruct {
	if name != "" {
		for _, s := range structs {
			if s.Name == name {
				return &s
			}
		}
		return nil
	}

	if len(structs) == 1 {
		return &structs[0]
	}

	return nil
}

func structNames(structs []sourceStruct) string {
	var names []string
	for _, s := range structs {
		names = append(names, "  - "+s.Name)
	}
	return strings.Join(names, "\n")
}

func toFieldContracts(fields []sourceField) []contract.FieldContract {
	out := make([]contract.FieldContract, len(fields))
	for i, f := range fields {
		out[i] = contract.FieldContract{
			Name:       f.Name,
			EnvKey:     f.EnvKey,
			Required:   f.Required,
			HasDefault: f.HasDefault,
			Default:    f.Default,
			Kind:       f.Kind,
		}
	}
	return out
}

func loadSources(envFiles []string, includeSystem bool) (map[string]string, []envcontract.LoadWarning, error) {
	env := make(map[string]string)
	var warnings []envcontract.LoadWarning

	for _, path := range envFiles {
		src := source.FileSource{Path: path}
		loaded, err := src.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("load %s: %w", path, err)
		}

		for key, value := range loaded.Values {
			env[key] = value
		}

		for _, w := range loaded.Warnings {
			warnings = append(warnings, envcontract.LoadWarning(w))
		}
	}

	if includeSystem {
		src := source.SystemSource{}
		loaded, err := src.Load()
		if err != nil {
			return nil, nil, fmt.Errorf("load system env: %w", err)
		}

		for key, value := range loaded.Values {
			env[key] = value
		}
	}

	return env, warnings, nil
}
