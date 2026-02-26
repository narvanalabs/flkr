package detector

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/narvanalabs/flkr/pkg/flkr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRustDetector_Actix(t *testing.T) {
	fsys := fstest.MapFS{
		"Cargo.toml": &fstest.MapFile{
			Data: []byte(`[package]
name = "my-app"
edition = "2021"

[dependencies]
actix-web = "4"
`),
		},
		"Cargo.lock": &fstest.MapFile{Data: []byte("version = 3\n")},
	}

	d := &RustDetector{}
	profile, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.True(t, matched)
	assert.Equal(t, flkr.LangRust, profile.Language)
	assert.Equal(t, flkr.PkgCargo, profile.PackageManager)
	assert.Equal(t, flkr.FrameworkActix, profile.Framework)
	assert.Equal(t, "./target/release/my-app", profile.StartCommand)
	assert.True(t, profile.HasLockfile)
}
