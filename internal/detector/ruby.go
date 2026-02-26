package detector

import (
	"context"
	"io/fs"
	"strings"

	"github.com/narvanalabs/flkr/pkg/flkr"
)

// RubyDetector detects Ruby applications.
type RubyDetector struct{}

func (d *RubyDetector) Name() string  { return "ruby" }
func (d *RubyDetector) Priority() int { return 50 }

func (d *RubyDetector) Detect(ctx context.Context, root fs.FS) (*flkr.AppProfile, bool, error) {
	if !fileExists(root, "Gemfile") {
		return nil, false, nil
	}

	profile := &flkr.AppProfile{
		Language:       flkr.LangRuby,
		PackageManager: flkr.PkgBundler,
		Confidence:     0.7,
		DetectedBy:     d.Name(),
		Port:           3000,
	}

	if fileExists(root, "Gemfile.lock") {
		profile.HasLockfile = true
		profile.LockfileType = "bundler"
	}

	// Read .ruby-version if it exists.
	if ver := readFileString(root, ".ruby-version"); ver != "" {
		profile.Version = strings.TrimSpace(ver)
	}

	// Detect Rails.
	gemfile := readFileString(root, "Gemfile")
	if strings.Contains(gemfile, "'rails'") || strings.Contains(gemfile, "\"rails\"") {
		profile.Framework = flkr.FrameworkRails
		profile.Confidence = 0.9
		profile.BuildCommand = "bundle exec rake assets:precompile"
		profile.StartCommand = "bundle exec rails server -b 0.0.0.0"
	}

	// Also check for config/routes.rb as a Rails indicator.
	if fileExists(root, "config/routes.rb") {
		profile.Framework = flkr.FrameworkRails
		profile.Confidence = 0.9
	}

	return profile, true, nil
}
