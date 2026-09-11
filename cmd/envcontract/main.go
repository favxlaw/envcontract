package main

import (
	"fmt"
	"io"
	"os"
)

const version = "0.1.0-dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: envcontract <command> [flags]")
		return 1
	}

	switch args[0] {
	case "version":
		fmt.Fprintln(stdout, "envcontract "+version)
		return 0
	default:
		fmt.Fprintf(stderr, "envcontract: unknown command %q\n", args[0])
		return 1
	}
}
