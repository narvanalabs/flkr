package detector

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/narvanalabs/flkr/pkg/flkr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPythonDetector_FastAPI(t *testing.T) {
	fsys := fstest.MapFS{
		"pyproject.toml": &fstest.MapFile{
			Data: []byte(`[project]
name = "my-app"
requires-python = ">=3.11"
dependencies = ["fastapi>=0.100.0", "uvicorn[standard]>=0.23.0"]
`),
		},
		"uv.lock": &fstest.MapFile{Data: []byte("version = 1\n")},
	}

	d := &PythonDetector{}
	profile, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.True(t, matched)
	assert.Equal(t, flkr.LangPython, profile.Language)
	assert.Equal(t, flkr.PkgUV, profile.PackageManager)
	assert.Equal(t, flkr.FrameworkFastAPI, profile.Framework)
	assert.Equal(t, "3.11", profile.Version)
	assert.True(t, profile.HasLockfile)
}

func TestPythonDetector_Django_Requirements(t *testing.T) {
	fsys := fstest.MapFS{
		"requirements.txt": &fstest.MapFile{
			Data: []byte("django>=4.2\npsycopg2-binary\n"),
		},
	}

	d := &PythonDetector{}
	profile, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.True(t, matched)
	assert.Equal(t, flkr.FrameworkDjango, profile.Framework)
	assert.Equal(t, flkr.PkgPip, profile.PackageManager)
}

func TestPythonDetector_NoMatch(t *testing.T) {
	fsys := fstest.MapFS{}
	d := &PythonDetector{}
	_, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.False(t, matched)
}
