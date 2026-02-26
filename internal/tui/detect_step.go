package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/narvanalabs/flkr/internal/detector"
	"github.com/narvanalabs/flkr/pkg/flkr"
)

// detectResultMsg carries the detection result back to the TUI.
type detectResultMsg struct {
	profile *flkr.AppProfile
	err     error
}

// runDetection starts detection in the background and returns a message.
func runDetection(path string) tea.Cmd {
	return func() tea.Msg {
		reg := detector.NewRegistry()
		profile, err := reg.DetectFromPath(context.Background(), path)
		return detectResultMsg{profile: profile, err: err}
	}
}

// renderDetecting shows a spinner while detection is running.
func renderDetecting(s spinner.Model) string {
	return fmt.Sprintf("\n %s Scanning repository...\n", s.View())
}

// renderDetectResult shows the detection results.
func renderDetectResult(profile *flkr.AppProfile) string {
	s := "\n" + successStyle.Render("  Detection complete!") + "\n\n"
	s += formatField("Language", string(profile.Language))
	if profile.Version != "" {
		s += formatField("Version", profile.Version)
	}
	s += formatField("Package Manager", string(profile.PackageManager))
	if profile.Framework != "" {
		s += formatField("Framework", string(profile.Framework))
	}
	if profile.BuildCommand != "" {
		s += formatField("Build Command", profile.BuildCommand)
	}
	if profile.StartCommand != "" {
		s += formatField("Start Command", profile.StartCommand)
	}
	if profile.Port != 0 {
		s += formatField("Port", fmt.Sprintf("%d", profile.Port))
	}
	s += formatField("Confidence", fmt.Sprintf("%.0f%%", profile.Confidence*100))
	s += "\n"
	return s
}

func formatField(key, value string) string {
	return "  " + keyStyle.Render(key) + valueStyle.Render(value) + "\n"
}
