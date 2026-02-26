package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose    bool
	jsonOutput bool
)

var rootCmd = &cobra.Command{
	Use:   "flkr",
	Short: "Detect application stacks and generate Nix flakes",
	Long:  `flkr scans repositories, detects the application stack, and generates a minimal flake.nix referencing the flkr-templates registry.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output in JSON format")
}
