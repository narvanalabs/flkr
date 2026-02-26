package detector

import (
	"context"
	"io/fs"
	"strings"

	"github.com/narvanalabs/flkr/internal/parser"
	"github.com/narvanalabs/flkr/pkg/flkr"
)

// NodeDetector detects Node.js applications.
type NodeDetector struct{}

func (d *NodeDetector) Name() string  { return "node" }
func (d *NodeDetector) Priority() int { return 10 }

func (d *NodeDetector) Detect(ctx context.Context, root fs.FS) (*flkr.AppProfile, bool, error) {
	if !fileExists(root, "package.json") {
		return nil, false, nil
	}

	pkg, err := parser.ParsePackageJSON(root, "package.json")
	if err != nil {
		return nil, false, err
	}

	profile := &flkr.AppProfile{
		Language:   flkr.LangNode,
		Confidence: 0.7,
		DetectedBy: d.Name(),
	}

	// Detect Node version from engines field.
	if pkg.Engines.Node != "" {
		profile.Version = cleanVersion(pkg.Engines.Node)
	}

	// Detect package manager from lockfiles.
	switch {
	case fileExists(root, "pnpm-lock.yaml"):
		profile.PackageManager = flkr.PkgPNPM
		profile.HasLockfile = true
		profile.LockfileType = "pnpm"
	case fileExists(root, "yarn.lock"):
		profile.PackageManager = flkr.PkgYarn
		profile.HasLockfile = true
		profile.LockfileType = "yarn"
	default:
		profile.PackageManager = flkr.PkgNPM
		if fileExists(root, "package-lock.json") {
			profile.HasLockfile = true
			profile.LockfileType = "npm"
		}
	}

	// Detect framework.
	d.detectFramework(pkg, profile)

	// Detect build/start commands from scripts.
	if cmd, ok := pkg.Scripts["build"]; ok {
		profile.BuildCommand = cmd
	}
	if cmd, ok := pkg.Scripts["start"]; ok {
		profile.StartCommand = cmd
	}

	// Default port.
	if profile.Port == 0 {
		profile.Port = 3000
	}

	return profile, true, nil
}

func (d *NodeDetector) detectFramework(pkg *parser.PackageJSON, profile *flkr.AppProfile) {
	switch {
	case pkg.HasDep("next"):
		profile.Framework = flkr.FrameworkNextJS
		profile.OutputDir = ".next"
		profile.Confidence = 0.9
	case pkg.HasDep("nuxt"):
		profile.Framework = flkr.FrameworkNuxt
		profile.OutputDir = ".output"
		profile.Confidence = 0.9
	case pkg.HasDep("@remix-run/node") || pkg.HasDep("@remix-run/react"):
		profile.Framework = flkr.FrameworkRemix
		profile.OutputDir = "build"
		profile.Confidence = 0.85
	case pkg.HasDep("vite"):
		profile.Framework = flkr.FrameworkVite
		profile.OutputDir = "dist"
		profile.Confidence = 0.8
	}
}

// cleanVersion strips common version prefixes/ranges to extract a bare version.
func cleanVersion(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimLeft(v, ">=^~<>!v")
	if i := strings.IndexAny(v, " |&"); i != -1 {
		v = v[:i]
	}
	return v
}
