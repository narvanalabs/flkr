package detector

import (
	"context"
	"io/fs"
	"strings"

	"github.com/narvanalabs/flkr/pkg/flkr"
)

// GoDetector detects Go applications.
type GoDetector struct{}

func (d *GoDetector) Name() string  { return "go" }
func (d *GoDetector) Priority() int { return 30 }

func (d *GoDetector) Detect(ctx context.Context, root fs.FS) (*flkr.AppProfile, bool, error) {
	if !fileExists(root, "go.mod") {
		return nil, false, nil
	}

	profile := &flkr.AppProfile{
		Language:       flkr.LangGo,
		PackageManager: flkr.PkgGoMod,
		Confidence:     0.8,
		DetectedBy:     d.Name(),
		BuildCommand:   "go build -o app .",
		StartCommand:   "./app",
		Port:           8080,
	}

	if fileExists(root, "go.sum") {
		profile.HasLockfile = true
		profile.LockfileType = "gomod"
	}

	// Parse go.mod for Go version and framework detection.
	content := readFileString(root, "go.mod")
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "go ") {
			profile.Version = strings.TrimPrefix(line, "go ")
		}
	}

	// Detect Gin framework.
	if strings.Contains(content, "github.com/gin-gonic/gin") {
		profile.Framework = flkr.FrameworkGin
		profile.Confidence = 0.9
	}

	return profile, true, nil
}
