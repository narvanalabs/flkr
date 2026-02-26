package detector

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/narvanalabs/flkr/pkg/flkr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeDetector_NextJS(t *testing.T) {
	fsys := fstest.MapFS{
		"package.json": &fstest.MapFile{
			Data: []byte(`{
				"name": "my-app",
				"scripts": {"build": "next build", "start": "next start"},
				"dependencies": {"next": "14.0.0", "react": "18.2.0"},
				"engines": {"node": ">=20.0.0"}
			}`),
		},
		"package-lock.json": &fstest.MapFile{Data: []byte(`{}`)},
	}

	d := &NodeDetector{}
	profile, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.True(t, matched)
	assert.Equal(t, flkr.LangNode, profile.Language)
	assert.Equal(t, flkr.PkgNPM, profile.PackageManager)
	assert.Equal(t, flkr.FrameworkNextJS, profile.Framework)
	assert.Equal(t, "20.0.0", profile.Version)
	assert.Equal(t, ".next", profile.OutputDir)
	assert.Equal(t, "next build", profile.BuildCommand)
	assert.Equal(t, "next start", profile.StartCommand)
	assert.True(t, profile.HasLockfile)
	assert.Equal(t, 3000, profile.Port)
	assert.InDelta(t, 0.9, profile.Confidence, 0.01)
}

func TestNodeDetector_Yarn(t *testing.T) {
	fsys := fstest.MapFS{
		"package.json": &fstest.MapFile{
			Data: []byte(`{"name": "app", "dependencies": {"vite": "5.0.0"}}`),
		},
		"yarn.lock": &fstest.MapFile{Data: []byte(`# yarn lockfile`)},
	}

	d := &NodeDetector{}
	profile, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.True(t, matched)
	assert.Equal(t, flkr.PkgYarn, profile.PackageManager)
	assert.Equal(t, flkr.FrameworkVite, profile.Framework)
	assert.True(t, profile.HasLockfile)
}

func TestNodeDetector_PNPM(t *testing.T) {
	fsys := fstest.MapFS{
		"package.json":    &fstest.MapFile{Data: []byte(`{"name": "app"}`)},
		"pnpm-lock.yaml": &fstest.MapFile{Data: []byte(`lockfileVersion: 6`)},
	}

	d := &NodeDetector{}
	profile, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.True(t, matched)
	assert.Equal(t, flkr.PkgPNPM, profile.PackageManager)
}

func TestNodeDetector_NoPackageJSON(t *testing.T) {
	fsys := fstest.MapFS{}
	d := &NodeDetector{}
	_, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.False(t, matched)
}
