package generator

import (
	"embed"
	"text/template"

	"github.com/narvanalabs/flkr/pkg/flkr"
)

//go:embed templates/flake.nix.tmpl
var templateFS embed.FS

var flakeTemplate = template.Must(
	template.ParseFS(templateFS, "templates/flake.nix.tmpl"),
)

// templateData is the view model passed to the flake.nix template.
type templateData struct {
	Name            string
	Ecosystem       string
	Version         string
	PackageManager  string
	Framework       string
	BuildCommand    string
	StartCommand    string
	OutputDir       string
	Port            int
	SystemDeps      []string
	EnvVars         []string
	TemplateVersion string
	AppVersion      string
	VendorHash      string // Nix expression: "null" for vendor/, or quoted hash string
}

// newTemplateData converts an AppProfile into template data.
func newTemplateData(profile *flkr.AppProfile, templateVersion string) templateData {
	name := string(profile.Language)
	if profile.Framework != "" {
		name = string(profile.Framework)
	}

	// For Go projects, emit vendorHash.
	// Quoted SRI hash if computed, "null" (use vendor/ dir) as fallback.
	var vendorHash string
	if profile.Language == flkr.LangGo {
		if profile.VendorHash != "" {
			vendorHash = `"` + profile.VendorHash + `"`
		} else {
			vendorHash = "null"
		}
	}

	return templateData{
		Name:            name + "-app",
		Ecosystem:       string(profile.Language),
		Version:         profile.Version,
		PackageManager:  string(profile.PackageManager),
		Framework:       string(profile.Framework),
		BuildCommand:    profile.BuildCommand,
		StartCommand:    profile.StartCommand,
		OutputDir:       profile.OutputDir,
		Port:            profile.Port,
		SystemDeps:      profile.SystemDeps,
		EnvVars:         profile.EnvVars,
		AppVersion:      profile.AppVersion,
		TemplateVersion: templateVersion,
		VendorHash:      vendorHash,
	}
}
