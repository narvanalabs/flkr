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

	if fileExists(root, "vendor") {
		profile.HasVendor = true
	}

	// Parse go.mod for module name, Go version, and framework detection.
	content := readFileString(root, "go.mod")
	moduleName := ""
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			moduleName = strings.TrimPrefix(line, "module ")
		}
		if strings.HasPrefix(line, "go ") {
			profile.Version = strings.TrimPrefix(line, "go ")
		}
	}

	// Default binary name from module path.
	binName := "app"
	if moduleName != "" {
		parts := strings.Split(moduleName, "/")
		binName = parts[len(parts)-1]
	}

	// Find where package main lives.
	mainPkg := findMainPackage(root, binName)
	profile.BuildCommand = "go build -o " + mainPkg.binName + " " + mainPkg.pkgPath
	profile.StartCommand = "./" + mainPkg.binName

	// Detect Gin framework.
	if strings.Contains(content, "github.com/gin-gonic/gin") {
		profile.Framework = flkr.FrameworkGin
		profile.Confidence = 0.9
	}

	return profile, true, nil
}

type mainPackageInfo struct {
	binName string // binary name (e.g. "flkr", "server")
	pkgPath string // Go package path (e.g. ".", "./cmd/server")
}

// findMainPackage locates the main package in a Go project.
// Checks root first, then cmd/<name>/ directories.
func findMainPackage(root fs.FS, moduleBinName string) mainPackageInfo {
	// Check root directory for package main.
	if dirHasMainPackage(root, ".") {
		return mainPackageInfo{binName: moduleBinName, pkgPath: "."}
	}

	// Check cmd/ subdirectories.
	cmdEntries, err := fs.ReadDir(root, "cmd")
	if err == nil {
		// Prefer a subdirectory matching the module name.
		for _, e := range cmdEntries {
			if e.IsDir() && e.Name() == moduleBinName {
				if dirHasMainPackage(root, "cmd/"+e.Name()) {
					return mainPackageInfo{binName: moduleBinName, pkgPath: "./cmd/" + moduleBinName}
				}
			}
		}
		// Otherwise take the first cmd/ subdirectory with package main.
		for _, e := range cmdEntries {
			if e.IsDir() && dirHasMainPackage(root, "cmd/"+e.Name()) {
				return mainPackageInfo{binName: e.Name(), pkgPath: "./cmd/" + e.Name()}
			}
		}
	}

	// Fallback: assume root.
	return mainPackageInfo{binName: moduleBinName, pkgPath: "."}
}

// dirHasMainPackage checks if a directory contains a Go file with "package main".
func dirHasMainPackage(root fs.FS, dir string) bool {
	entries, err := fs.ReadDir(root, dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		// Skip test files.
		if strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		path := e.Name()
		if dir != "." {
			path = dir + "/" + e.Name()
		}
		content := readFileString(root, path)
		// Check the first non-empty, non-comment line for package declaration.
		for _, line := range strings.Split(content, "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "//") {
				continue
			}
			if line == "package main" {
				return true
			}
			break
		}
	}
	return false
}
