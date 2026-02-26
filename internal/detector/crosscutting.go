package detector

import (
	"context"
	"io/fs"
	"strings"

	"github.com/narvanalabs/flkr/pkg/flkr"
)

// CrosscuttingDetector enriches a profile with data from .env.example,
// Procfile, and Makefile. It doesn't match on its own — it's used as
// a post-processing step.
type CrosscuttingDetector struct{}

func (d *CrosscuttingDetector) Name() string  { return "crosscutting" }
func (d *CrosscuttingDetector) Priority() int { return 100 }

func (d *CrosscuttingDetector) Detect(ctx context.Context, root fs.FS) (*flkr.AppProfile, bool, error) {
	profile := &flkr.AppProfile{}
	matched := false

	// Parse .env.example for env var keys.
	if envKeys := parseEnvExample(root); len(envKeys) > 0 {
		profile.EnvVars = envKeys
		matched = true
	}

	// Parse Procfile for start command.
	if cmd := parseProcfile(root); cmd != "" {
		profile.StartCommand = cmd
		matched = true
	}

	return profile, matched, nil
}

// parseEnvExample extracts variable names from .env.example.
func parseEnvExample(root fs.FS) []string {
	content := readFileString(root, ".env.example")
	if content == "" {
		return nil
	}
	var keys []string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if k, _, ok := strings.Cut(line, "="); ok {
			k = strings.TrimSpace(k)
			if k != "" {
				keys = append(keys, k)
			}
		}
	}
	return keys
}

// parseProcfile extracts the web process command from a Procfile.
func parseProcfile(root fs.FS) string {
	content := readFileString(root, "Procfile")
	if content == "" {
		return ""
	}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "web:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "web:"))
		}
	}
	return ""
}
