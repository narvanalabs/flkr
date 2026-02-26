package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/narvanalabs/flkr/pkg/flkr"
)

// buildReviewForm creates a huh form for reviewing/editing the detected profile.
func buildReviewForm(profile *flkr.AppProfile) *huh.Form {
	portStr := strconv.Itoa(profile.Port)

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Language").
				Options(
					huh.NewOption("Node.js", "node"),
					huh.NewOption("Python", "python"),
					huh.NewOption("Go", "go"),
					huh.NewOption("Rust", "rust"),
					huh.NewOption("Ruby", "ruby"),
					huh.NewOption("Elixir", "elixir"),
					huh.NewOption("PHP", "php"),
					huh.NewOption("Java", "java"),
				).
				Value((*string)(&profile.Language)),

			huh.NewInput().
				Title("Version").
				Value(&profile.Version),

			huh.NewInput().
				Title("Build Command").
				Value(&profile.BuildCommand),

			huh.NewInput().
				Title("Start Command").
				Value(&profile.StartCommand),

			huh.NewInput().
				Title("Port").
				Value(&portStr).
				Validate(func(s string) error {
					if s == "" {
						return nil
					}
					n, err := strconv.Atoi(s)
					if err != nil || n < 0 || n > 65535 {
						return fmt.Errorf("must be a valid port (0-65535)")
					}
					return nil
				}),
		),
	).WithShowHelp(true)
}

// applyFormValues updates the profile with values from the form.
// Port needs special handling since huh operates on strings.
func applyFormValues(profile *flkr.AppProfile, portStr string) {
	if n, err := strconv.Atoi(portStr); err == nil {
		profile.Port = n
	}
}
