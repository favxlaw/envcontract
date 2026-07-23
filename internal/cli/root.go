package cli

import (
	"github.com/spf13/cobra"
)

const version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "envcontract",
	Short: "Detect environment variable drift in Go config structs",
	Long: `envcontract checks your Go config structs against environment
variables to catch missing, renamed, unused, or mistyped env vars
before they cause incidents.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(schemaCmd)
	rootCmd.AddCommand(initCmd)
}
