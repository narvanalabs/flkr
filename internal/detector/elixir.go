package detector

import (
	"context"
	"io/fs"
	"strings"

	"github.com/narvanalabs/flkr/pkg/flkr"
)

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

	// Detect Phoenix from mix.exs deps.
	mixExs := readFileString(root, "mix.exs")
	if strings.Contains(mixExs, ":phoenix") {
		profile.Framework = flkr.FrameworkPhoenix
		profile.Confidence = 0.9
		profile.SystemDeps = []string{"inotify-tools"}
	}

	return profile, true, nil
}
