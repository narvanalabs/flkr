package flkr

import (
	"fmt"
	"strings"
)

// Language represents a detected programming language.
type Language string

const (
	LangNode   Language = "node"
	LangPython Language = "python"
	LangGo     Language = "go"
	LangRust   Language = "rust"
	LangRuby   Language = "ruby"
	LangElixir Language = "elixir"
	LangPHP    Language = "php"
	LangJava   Language = "java"
)

// PackageManager represents a detected package manager.
type PackageManager string

const (
	PkgNPM      PackageManager = "npm"
	PkgYarn     PackageManager = "yarn"
	PkgPNPM    PackageManager = "pnpm"
	PkgPip      PackageManager = "pip"
	PkgPoetry   PackageManager = "poetry"
	PkgPipenv   PackageManager = "pipenv"
	PkgUV       PackageManager = "uv"
	PkgGoMod    PackageManager = "gomod"
	PkgCargo    PackageManager = "cargo"
	PkgBundler  PackageManager = "bundler"
	PkgMix      PackageManager = "mix"
	PkgComposer PackageManager = "composer"
	PkgMaven    PackageManager = "maven"
	PkgGradle   PackageManager = "gradle"
)

// Framework represents a detected web framework.
type Framework string

const (
	FrameworkNone    Framework = ""
	FrameworkNextJS  Framework = "nextjs"
	FrameworkNuxt    Framework = "nuxt"
	FrameworkRemix   Framework = "remix"
	FrameworkVite    Framework = "vite"
	FrameworkDjango  Framework = "django"
	FrameworkFlask   Framework = "flask"
	FrameworkFastAPI Framework = "fastapi"
	FrameworkGin     Framework = "gin"
	FrameworkActix   Framework = "actix"
	FrameworkRails   Framework = "rails"
	FrameworkPhoenix Framework = "phoenix"
	FrameworkLaravel Framework = "laravel"
	FrameworkSpring  Framework = "spring"
)

// AppProfile represents the full detected profile of an application.
type AppProfile struct {
	Language       Language       `json:"language"`
	Version        string         `json:"version,omitempty"`
	PackageManager PackageManager `json:"packageManager"`
	Framework      Framework      `json:"framework,omitempty"`
	BuildCommand   string         `json:"buildCommand,omitempty"`
	StartCommand   string         `json:"startCommand,omitempty"`
	OutputDir      string         `json:"outputDir,omitempty"`
	Port           int            `json:"port,omitempty"`
	SystemDeps     []string       `json:"systemDeps,omitempty"`
	EnvVars        []string       `json:"envVars,omitempty"`
	HasLockfile    bool           `json:"hasLockfile"`
	LockfileType   string         `json:"lockfileType,omitempty"`
	HasVendor      bool           `json:"hasVendor,omitempty"`
	VendorHash     string         `json:"vendorHash,omitempty"`
	Confidence     float64        `json:"confidence"`
	DetectedBy     string         `json:"detectedBy,omitempty"`
}

// Validate checks that the profile has the minimum required fields.
func (p *AppProfile) Validate() error {
	var errs []string
	if p.Language == "" {
		errs = append(errs, "language is required")
	}
	if p.PackageManager == "" {
		errs = append(errs, "packageManager is required")
	}
	if p.Confidence < 0 || p.Confidence > 1 {
		errs = append(errs, "confidence must be between 0 and 1")
	}
	if len(errs) > 0 {
		return fmt.Errorf("invalid profile: %s", strings.Join(errs, "; "))
	}
	return nil
}

// Merge overlays another profile onto this one. Non-zero fields in other
// take precedence. Slices are appended and deduplicated. Confidence takes
// the higher value.
func (p *AppProfile) Merge(other *AppProfile) {
	if other == nil {
		return
	}
	if other.Language != "" {
		p.Language = other.Language
	}
	if other.Version != "" {
		p.Version = other.Version
	}
	if other.PackageManager != "" {
		p.PackageManager = other.PackageManager
	}
	if other.Framework != "" {
		p.Framework = other.Framework
	}
	if other.BuildCommand != "" {
		p.BuildCommand = other.BuildCommand
	}
	if other.StartCommand != "" {
		p.StartCommand = other.StartCommand
	}
	if other.OutputDir != "" {
		p.OutputDir = other.OutputDir
	}
	if other.Port != 0 {
		p.Port = other.Port
	}
	if other.HasLockfile {
		p.HasLockfile = true
	}
	if other.LockfileType != "" {
		p.LockfileType = other.LockfileType
	}
	if other.Confidence > p.Confidence {
		p.Confidence = other.Confidence
	}
	if other.DetectedBy != "" {
		p.DetectedBy = other.DetectedBy
	}
	p.SystemDeps = mergeUnique(p.SystemDeps, other.SystemDeps)
	p.EnvVars = mergeUnique(p.EnvVars, other.EnvVars)
}

func mergeUnique(a, b []string) []string {
	seen := make(map[string]struct{}, len(a))
	for _, s := range a {
		seen[s] = struct{}{}
	}
	for _, s := range b {
		if _, ok := seen[s]; !ok {
			a = append(a, s)
			seen[s] = struct{}{}
		}
	}
	return a
}
