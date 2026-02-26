package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/narvanalabs/flkr/internal/tui"
	"github.com/spf13/cobra"
)

var initTemplateVersion string

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Interactive wizard to detect and generate a flake.nix",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		model := tui.New(path, initTemplateVersion)
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&initTemplateVersion, "template-version", "", "pin flkr-templates to a specific revision")
	rootCmd.AddCommand(initCmd)
}
