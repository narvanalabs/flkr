package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/narvanalabs/flkr/internal/detector"
	"github.com/spf13/cobra"
)

var detectCmd = &cobra.Command{
	Use:   "detect [path]",
	Short: "Detect the application stack in a repository",
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

		if jsonOutput {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(profile)
		}

		fmt.Printf("Language:        %s\n", profile.Language)
		if profile.Version != "" {
			fmt.Printf("Version:         %s\n", profile.Version)
		}
		fmt.Printf("Package Manager: %s\n", profile.PackageManager)
		if profile.Framework != "" {
			fmt.Printf("Framework:       %s\n", profile.Framework)
		}
		if profile.BuildCommand != "" {
			fmt.Printf("Build Command:   %s\n", profile.BuildCommand)
		}
		if profile.StartCommand != "" {
			fmt.Printf("Start Command:   %s\n", profile.StartCommand)
		}
		if profile.OutputDir != "" {
			fmt.Printf("Output Dir:      %s\n", profile.OutputDir)
		}
		if profile.Port != 0 {
			fmt.Printf("Port:            %d\n", profile.Port)
		}
		if len(profile.EnvVars) > 0 {
			fmt.Printf("Env Vars:        %v\n", profile.EnvVars)
		}
		fmt.Printf("Confidence:      %.0f%%\n", profile.Confidence*100)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)
}
