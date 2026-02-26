package detector

import (
	"context"
	"io/fs"
	"strings"

	"github.com/narvanalabs/flkr/internal/parser"
	"github.com/narvanalabs/flkr/pkg/flkr"
)

// PythonDetector detects Python applications.
type PythonDetector struct{}

func (d *PythonDetector) Name() string  { return "python" }
func (d *PythonDetector) Priority() int { return 20 }

func (d *PythonDetector) Detect(ctx context.Context, root fs.FS) (*flkr.AppProfile, bool, error) {
	hasPyproject := fileExists(root, "pyproject.toml")
	hasRequirements := fileExists(root, "requirements.txt")
	hasPipfile := fileExists(root, "Pipfile")
	hasSetupPy := fileExists(root, "setup.py")

	if !hasPyproject && !hasRequirements && !hasPipfile && !hasSetupPy {
		return nil, false, nil
	}

	profile := &flkr.AppProfile{
		Language:   flkr.LangPython,
		Confidence: 0.7,
		DetectedBy: d.Name(),
		Port:       8000,
	}

	// Detect package manager.
	switch {
	case fileExists(root, "uv.lock"):
		profile.PackageManager = flkr.PkgUV
		profile.HasLockfile = true
		profile.LockfileType = "uv"
	case fileExists(root, "poetry.lock"):
		profile.PackageManager = flkr.PkgPoetry
		profile.HasLockfile = true
		profile.LockfileType = "poetry"
	case hasPipfile:
		profile.PackageManager = flkr.PkgPipenv
		if fileExists(root, "Pipfile.lock") {
			profile.HasLockfile = true
			profile.LockfileType = "pipenv"
		}
	default:
		profile.PackageManager = flkr.PkgPip
	}

	// Parse pyproject.toml for framework detection.
	if hasPyproject {
		pyproj, err := parser.ParsePyprojectTOML(root, "pyproject.toml")
		if err == nil {
			d.detectFramework(pyproj, profile)
			if pyproj.Project.Version != "" {
				profile.AppVersion = pyproj.Project.Version
			}
			if pyproj.Project.RequiresPython != "" {
				profile.Version = cleanVersion(pyproj.Project.RequiresPython)
			}
		}
	}

	// Fallback: check requirements.txt for framework hints.
	if profile.Framework == "" && hasRequirements {
		d.detectFrameworkFromRequirements(root, profile)
	}

	// Set start commands based on framework.
	switch profile.Framework {
	case flkr.FrameworkDjango:
		profile.StartCommand = "python manage.py runserver 0.0.0.0:8000"
	case flkr.FrameworkFlask:
		profile.StartCommand = "flask run --host=0.0.0.0"
		profile.Port = 5000
	case flkr.FrameworkFastAPI:
		profile.StartCommand = "uvicorn main:app --host 0.0.0.0 --port 8000"
	}

	return profile, true, nil
}

func (d *PythonDetector) detectFramework(pyproj *parser.PyprojectTOML, profile *flkr.AppProfile) {
	switch {
	case pyproj.HasDep("django"):
		profile.Framework = flkr.FrameworkDjango
		profile.Confidence = 0.9
	case pyproj.HasDep("flask"):
		profile.Framework = flkr.FrameworkFlask
		profile.Confidence = 0.85
	case pyproj.HasDep("fastapi"):
		profile.Framework = flkr.FrameworkFastAPI
		profile.Confidence = 0.9
	}
}

func (d *PythonDetector) detectFrameworkFromRequirements(root fs.FS, profile *flkr.AppProfile) {
	content := readFileString(root, "requirements.txt")
	lines := strings.Split(strings.ToLower(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "django"):
			profile.Framework = flkr.FrameworkDjango
			profile.Confidence = 0.85
			return
		case strings.HasPrefix(line, "flask"):
			profile.Framework = flkr.FrameworkFlask
			profile.Confidence = 0.8
			return
		case strings.HasPrefix(line, "fastapi"):
			profile.Framework = flkr.FrameworkFastAPI
			profile.Confidence = 0.85
			return
		}
	}
}
