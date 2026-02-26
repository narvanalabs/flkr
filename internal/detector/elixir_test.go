package detector

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/narvanalabs/flkr/pkg/flkr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestElixirDetector_Phoenix(t *testing.T) {
	fsys := fstest.MapFS{
		"mix.exs": &fstest.MapFile{
			Data: []byte(`defmodule MyApp.MixProject do
  defp deps do
    [{:phoenix, "~> 1.7"}]
  end
end
`),
		},
		"mix.lock": &fstest.MapFile{Data: []byte(`%{"phoenix": {:hex}}`)},
	}

	d := &ElixirDetector{}
	profile, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.True(t, matched)
	assert.Equal(t, flkr.LangElixir, profile.Language)
	assert.Equal(t, flkr.FrameworkPhoenix, profile.Framework)
	assert.Equal(t, 4000, profile.Port)
	assert.True(t, profile.HasLockfile)
}
