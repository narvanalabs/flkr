package tui

import (
	"github.com/charmbracelet/huh"
	"github.com/narvanalabs/flkr/internal/generator"
	"github.com/narvanalabs/flkr/pkg/flkr"
)

// generatePreview renders the flake.nix content for preview.
func generatePreview(profile *flkr.AppProfile, templateVersion string) (string, error) {
	gen := &generator.DefaultGenerator{}
	result, err := gen.Generate(profile, generator.Options{
		DryRun:          true,
		TemplateVersion: templateVersion,
	})
	if err != nil {
		return "", err
	}
	return result.FlakeContent, nil
}

// buildConfirmForm creates a confirmation form.
func buildConfirmForm(confirmed *bool) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Write flake.nix?").
				Affirmative("Yes").
				Negative("No").
				Value(confirmed),
		),
	)
}
