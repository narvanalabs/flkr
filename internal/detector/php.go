package detector

import (
	"context"
	"io/fs"

	"github.com/narvanalabs/flkr/internal/parser"
	"github.com/narvanalabs/flkr/pkg/flkr"
)

// PHPDetector detects PHP applications.
type PHPDetector struct{}

func (d *PHPDetector) Name() string  { return "php" }
func (d *PHPDetector) Priority() int { return 70 }

func (d *PHPDetector) Detect(ctx context.Context, root fs.FS) (*flkr.AppProfile, bool, error) {
	if !fileExists(root, "composer.json") {
		return nil, false, nil
	}

	profile := &flkr.AppProfile{
		Language:       flkr.LangPHP,
		PackageManager: flkr.PkgComposer,
		Confidence:     0.7,
		DetectedBy:     d.Name(),
		Port:           8000,
	}

	if fileExists(root, "composer.lock") {
		profile.HasLockfile = true
		profile.LockfileType = "composer"
	}

	// Parse composer.json for framework detection.
	comp, err := parser.ParseComposerJSON(root, "composer.json")
	if err == nil {
		// Extract project version.
		if comp.Version != "" {
			profile.AppVersion = comp.Version
		}

		// Detect PHP version.
		if v, ok := comp.Require["php"]; ok {
			profile.Version = cleanVersion(v)
		}

		// Detect Laravel.
		if comp.HasRequire("laravel/framework") {
			profile.Framework = flkr.FrameworkLaravel
			profile.Confidence = 0.9
			profile.BuildCommand = "composer install --no-dev --optimize-autoloader"
			profile.StartCommand = "php artisan serve --host=0.0.0.0 --port=8000"
			profile.OutputDir = "public"
		}
	}

	return profile, true, nil
}
