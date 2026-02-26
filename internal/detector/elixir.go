package detector

import (
	"context"
	"io/fs"
	"regexp"
	"strings"

	"github.com/narvanalabs/flkr/pkg/flkr"
)

var mixVersionRe = regexp.MustCompile(`version:\s*"([^"]+)"`)

// extractMixVersion extracts the version from a mix.exs project definition.
func extractMixVersion(content string) string {
	m := mixVersionRe.FindStringSubmatch(content)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}

// ElixirDetector detects Elixir applications.
type ElixirDetector struct{}

func (d *ElixirDetector) Name() string  { return "elixir" }
func (d *ElixirDetector) Priority() int { return 60 }

func (d *ElixirDetector) Detect(ctx context.Context, root fs.FS) (*flkr.AppProfile, bool, error) {
	if !fileExists(root, "mix.exs") {
		return nil, false, nil
	}

	profile := &flkr.AppProfile{
		Language:       flkr.LangElixir,
		PackageManager: flkr.PkgMix,
		Confidence:     0.8,
		DetectedBy:     d.Name(),
		BuildCommand:   "mix do deps.get, compile",
		StartCommand:   "mix phx.server",
		Port:           4000,
	}

	if fileExists(root, "mix.lock") {
		profile.HasLockfile = true
		profile.LockfileType = "mix"
	}

	// Read .elixir-version if it exists.
	if ver := readFileString(root, ".elixir-version"); ver != "" {
		profile.Version = strings.TrimSpace(ver)
	}

	// Extract project version from mix.exs.
	mixExs := readFileString(root, "mix.exs")
	if v := extractMixVersion(mixExs); v != "" {
		profile.AppVersion = v
	}

	// Detect Phoenix from mix.exs deps.
	if strings.Contains(mixExs, ":phoenix") {
		profile.Framework = flkr.FrameworkPhoenix
		profile.Confidence = 0.9
		profile.SystemDeps = []string{"inotify-tools"}
	}

	return profile, true, nil
}
