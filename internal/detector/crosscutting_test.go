package detector

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCrosscutting_EnvExample(t *testing.T) {
	fsys := fstest.MapFS{
		".env.example": &fstest.MapFile{
			Data: []byte("DATABASE_URL=postgres://localhost\n# Comment\nSECRET_KEY=changeme\n\nAPI_URL=\n"),
		},
	}

	d := &CrosscuttingDetector{}
	profile, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.True(t, matched)
	assert.Equal(t, []string{"DATABASE_URL", "SECRET_KEY", "API_URL"}, profile.EnvVars)
}

func TestCrosscutting_Procfile(t *testing.T) {
	fsys := fstest.MapFS{
		"Procfile": &fstest.MapFile{
			Data: []byte("web: node server.js\nworker: node worker.js\n"),
		},
	}

	d := &CrosscuttingDetector{}
	profile, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.True(t, matched)
	assert.Equal(t, "node server.js", profile.StartCommand)
}

func TestCrosscutting_NoFiles(t *testing.T) {
	fsys := fstest.MapFS{}
	d := &CrosscuttingDetector{}
	_, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.False(t, matched)
}
