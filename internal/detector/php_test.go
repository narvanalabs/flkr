package detector

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/narvanalabs/flkr/pkg/flkr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPHPDetector_Laravel(t *testing.T) {
	fsys := fstest.MapFS{
		"composer.json": &fstest.MapFile{
			Data: []byte(`{"require": {"php": "^8.2", "laravel/framework": "^10.0"}}`),
		},
		"composer.lock": &fstest.MapFile{Data: []byte(`{}`)},
	}

	d := &PHPDetector{}
	profile, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.True(t, matched)
	assert.Equal(t, flkr.LangPHP, profile.Language)
	assert.Equal(t, flkr.FrameworkLaravel, profile.Framework)
	assert.Equal(t, "8.2", profile.Version)
	assert.Equal(t, "public", profile.OutputDir)
	assert.True(t, profile.HasLockfile)
}
