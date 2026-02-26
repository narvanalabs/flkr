package detector

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/narvanalabs/flkr/pkg/flkr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry_DetectAll_Node(t *testing.T) {
	fsys := fstest.MapFS{
		"package.json": &fstest.MapFile{
			Data: []byte(`{"name": "app", "dependencies": {"next": "14.0.0"}}`),
		},
	}

	reg := NewRegistry()
	profiles, err := reg.DetectAll(context.Background(), fsys)
	require.NoError(t, err)
	require.Len(t, profiles, 1)
	assert.Equal(t, flkr.LangNode, profiles[0].Language)
}

func TestRegistry_DetectAll_MultiLanguage(t *testing.T) {
	fsys := fstest.MapFS{
		"package.json": &fstest.MapFile{
			Data: []byte(`{"name": "frontend", "dependencies": {"vite": "5.0.0"}}`),
		},
		"go.mod": &fstest.MapFile{
			Data: []byte("module myapp\n\ngo 1.22.0\n"),
		},
	}

	reg := NewRegistry()
	profiles, err := reg.DetectAll(context.Background(), fsys)
	require.NoError(t, err)
	assert.Len(t, profiles, 2)
}

func TestRegistry_DetectBest_WithCrosscutting(t *testing.T) {
	fsys := fstest.MapFS{
		"package.json": &fstest.MapFile{
			Data: []byte(`{"name": "app", "dependencies": {"next": "14.0.0"}}`),
		},
		".env.example": &fstest.MapFile{
			Data: []byte("DATABASE_URL=\nSECRET_KEY=\n"),
		},
		"Procfile": &fstest.MapFile{
			Data: []byte("web: npm start\n"),
		},
	}

	reg := NewRegistry()
	profile, err := reg.DetectBest(context.Background(), fsys)
	require.NoError(t, err)
	require.NotNil(t, profile)
	assert.Equal(t, flkr.FrameworkNextJS, profile.Framework)
	assert.Contains(t, profile.EnvVars, "DATABASE_URL")
	assert.Contains(t, profile.EnvVars, "SECRET_KEY")
	assert.Equal(t, "npm start", profile.StartCommand)
}

func TestRegistry_DetectAll_Empty(t *testing.T) {
	fsys := fstest.MapFS{}
	reg := NewRegistry()
	profiles, err := reg.DetectAll(context.Background(), fsys)
	require.NoError(t, err)
	assert.Empty(t, profiles)
}
