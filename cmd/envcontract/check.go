package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/favxlaw/envcontract"
)

func cmdCheck(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	fs.SetOutput(stderr)

	schemaPath := fs.String("schema", "", "path to schema JSON produced by ExportSchema (required)")
	envPath := fs.String("env", "", "path to a .env file to validate")
	useSystem := fs.Bool("system", false, "also validate against the system process environment")
	checkUnused := fs.Bool("unused", false, "warn about env vars not referenced by the schema")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	if *schemaPath == "" {
		fmt.Fprintln(stderr, "check: -schema is required")
		fs.Usage()
		return 2
	}

	data, err := os.ReadFile(*schemaPath)
	if err != nil {
		fmt.Fprintf(stderr, "check: read schema: %v\n", err)
		return 1
	}

	contracts, err := envcontract.ParseSchema(data)
	if err != nil {
		fmt.Fprintf(stderr, "check: %v\n", err)
		return 1
	}

	var opts []envcontract.Option
	if *envPath != "" {
		opts = append(opts, envcontract.WithFile(*envPath))
	}
	if *useSystem || *envPath == "" {
		opts = append(opts, envcontract.WithSystemEnv())
	}
	if *checkUnused {
		opts = append(opts, envcontract.WithUnusedCheck())
	}

	result, err := envcontract.CheckSchema(contracts, opts...)
	if err != nil {
		fmt.Fprintf(stderr, "check: %v\n", err)
		return 1
	}

	if err := result.Print(stdout); err != nil {
		fmt.Fprintf(stderr, "check: %v\n", err)
		return 1
	}

	if result.HasErrors() {
		return 1
	}

	return 0
}
