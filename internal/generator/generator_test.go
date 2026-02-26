package generator

import (
	"strings"
	"testing"

	"github.com/narvanalabs/flkr/pkg/flkr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultGenerator_NextJS(t *testing.T) {
	profile := &flkr.AppProfile{
		Language:       flkr.LangNode,
		Version:        "20.0.0",
		PackageManager: flkr.PkgNPM,
		Framework:      flkr.FrameworkNextJS,
		BuildCommand:   "next build",
		StartCommand:   "next start",
		OutputDir:      ".next",
		Port:           3000,
		EnvVars:        []string{"DATABASE_URL", "SECRET_KEY"},
	}

	gen := &DefaultGenerator{}
	result, err := gen.Generate(profile, Options{DryRun: true})
	require.NoError(t, err)
	require.NotNil(t, result)

	content := result.FlakeContent
	assert.Contains(t, content, `ecosystem = "node"`)
	assert.Contains(t, content, `framework = "nextjs"`)
	assert.Contains(t, content, `version = "20.0.0"`)
	assert.Contains(t, content, `port = 3000`)
	assert.Contains(t, content, `"DATABASE_URL"`)
	assert.Contains(t, content, "flkr-templates")
	assert.Empty(t, result.OutputPath)
}

func TestDefaultGenerator_TemplateVersion(t *testing.T) {
	profile := &flkr.AppProfile{
		Language:       flkr.LangGo,
		PackageManager: flkr.PkgGoMod,
		BuildCommand:   "go build -o app .",
	}

	gen := &DefaultGenerator{}
	result, err := gen.Generate(profile, Options{
		DryRun:          true,
		TemplateVersion: "v1.0.0",
	})
	require.NoError(t, err)
	assert.Contains(t, result.FlakeContent, "flkr-templates/v1.0.0")
}

func TestDefaultGenerator_MinimalProfile(t *testing.T) {
	profile := &flkr.AppProfile{
		Language:       flkr.LangRust,
		PackageManager: flkr.PkgCargo,
	}

	gen := &DefaultGenerator{}
	result, err := gen.Generate(profile, Options{DryRun: true})
	require.NoError(t, err)
	// Should not contain empty optional fields.
	assert.False(t, strings.Contains(result.FlakeContent, `version = ""`))
	assert.False(t, strings.Contains(result.FlakeContent, `framework = ""`))
}
