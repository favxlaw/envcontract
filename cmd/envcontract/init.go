package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/favxlaw/envcontract"
	"github.com/favxlaw/envcontract/generator"
)

func cmdInit(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(stderr)

	schemaPath := fs.String("schema", "", "path to schema JSON produced by ExportSchema (required)")
	outPath := fs.String("out", ".env.example", "path to write the generated file")
	force := fs.Bool("force", false, "overwrite the output file if it already exists")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	if *schemaPath == "" {
		fmt.Fprintln(stderr, "init: -schema is required")
		fs.Usage()
		return 2
	}

	data, err := os.ReadFile(*schemaPath)
	if err != nil {
		fmt.Fprintf(stderr, "init: read schema: %v\n", err)
		return 1
	}

	contracts, err := envcontract.ParseSchema(data)
	if err != nil {
		fmt.Fprintf(stderr, "init: %v\n", err)
		return 1
	}

	if !*force {
		if _, err := os.Stat(*outPath); err == nil {
			fmt.Fprintf(stderr, "init: %s already exists (use -force to overwrite)\n", *outPath)
			return 1
		} else if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(stderr, "init: %v\n", err)
			return 1
		}
	}

	content := generator.EnvExample(contracts)

	if err := os.WriteFile(*outPath, []byte(content), 0o644); err != nil {
		fmt.Fprintf(stderr, "init: write %s: %v\n", *outPath, err)
		return 1
	}

	fmt.Fprintf(stdout, "wrote %s\n", *outPath)
	return 0
}
