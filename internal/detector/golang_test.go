package detector

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/narvanalabs/flkr/pkg/flkr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoDetector_Gin(t *testing.T) {
	fsys := fstest.MapFS{
		"go.mod": &fstest.MapFile{
			Data: []byte("module myapp\n\ngo 1.22.0\n\nrequire github.com/gin-gonic/gin v1.9.1\n"),
		},
		"go.sum": &fstest.MapFile{Data: []byte("github.com/gin-gonic/gin v1.9.1 h1:abc=\n")},
	}

	d := &GoDetector{}
	profile, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.True(t, matched)
	assert.Equal(t, flkr.LangGo, profile.Language)
	assert.Equal(t, flkr.PkgGoMod, profile.PackageManager)
	assert.Equal(t, flkr.FrameworkGin, profile.Framework)
	assert.Equal(t, "1.22.0", profile.Version)
	assert.True(t, profile.HasLockfile)
}

func TestGoDetector_Plain(t *testing.T) {
	fsys := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte("module myapp\n\ngo 1.22.0\n")},
	}

	d := &GoDetector{}
	profile, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.True(t, matched)
	assert.Equal(t, flkr.FrameworkNone, profile.Framework)
}
