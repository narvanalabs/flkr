package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/narvanalabs/flkr/internal/detector"
	"github.com/narvanalabs/flkr/internal/generator"
	"github.com/spf13/cobra"
)

var (
	dryRun          bool
	templateVersion string
	outputPath      string
)

var generateCmd = &cobra.Command{
	Use:   "generate [path]",
	Short: "Generate a flake.nix for the detected application stack",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		reg := detector.NewRegistry()
		profile, err := reg.DetectFromPath(context.Background(), path)
		if err != nil {
			return err
		}
		if profile == nil {
			fmt.Fprintln(os.Stderr, "no application stack detected")
			os.Exit(1)
		}

		out := outputPath
		if out == "" {
			out = filepath.Join(path, "flake.nix")
		}

		gen := &generator.DefaultGenerator{}
		result, err := gen.Generate(profile, generator.Options{
			OutputPath:      out,
			TemplateVersion: templateVersion,
			DryRun:          dryRun,
		})
		if err != nil {
			return err
		}

		if dryRun {
			fmt.Print(result.FlakeContent)
		} else {
			fmt.Printf("wrote %s\n", result.OutputPath)
		}
		return nil
	},
}

func init() {
	generateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "print flake.nix to stdout instead of writing")
	generateCmd.Flags().StringVar(&templateVersion, "template-version", "", "pin flkr-templates to a specific revision")
	generateCmd.Flags().StringVarP(&outputPath, "output", "o", "", "output file path (default: <path>/flake.nix)")
	rootCmd.AddCommand(generateCmd)
}
