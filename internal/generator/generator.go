package generator

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/narvanalabs/flkr/pkg/flkr"
)

// Options configures generation.
type Options struct {
	OutputPath      string
	TemplateVersion string
	DryRun          bool
}

// Generator renders flake.nix files.
type Generator interface {
	Generate(profile *flkr.AppProfile, opts Options) (*flkr.GenerateResult, error)
}

// DefaultGenerator uses text/template to render flake.nix.
type DefaultGenerator struct{}

// Generate renders a flake.nix from the given profile.
func (g *DefaultGenerator) Generate(profile *flkr.AppProfile, opts Options) (*flkr.GenerateResult, error) {
	data := newTemplateData(profile, opts.TemplateVersion)

	var buf bytes.Buffer
	if err := flakeTemplate.Execute(&buf, data); err != nil {
		return nil, err
	}

	result := &flkr.GenerateResult{
		FlakeContent: buf.String(),
	}

	if !opts.DryRun && opts.OutputPath != "" {
		dir := filepath.Dir(opts.OutputPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(opts.OutputPath, buf.Bytes(), 0o644); err != nil {
			return nil, err
		}
		result.OutputPath = opts.OutputPath
	}

	return result, nil
}
