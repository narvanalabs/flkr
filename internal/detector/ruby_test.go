package detector

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/narvanalabs/flkr/pkg/flkr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRubyDetector_Rails(t *testing.T) {
	fsys := fstest.MapFS{
		"Gemfile": &fstest.MapFile{
			Data: []byte(`source "https://rubygems.org"
gem "rails", "~> 7.1"
`),
		},
		"Gemfile.lock":     &fstest.MapFile{Data: []byte("GEM\n")},
		"config/routes.rb": &fstest.MapFile{Data: []byte("Rails.application.routes.draw do\nend\n")},
		".ruby-version":    &fstest.MapFile{Data: []byte("3.2.2\n")},
	}

	d := &RubyDetector{}
	profile, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.True(t, matched)
	assert.Equal(t, flkr.LangRuby, profile.Language)
	assert.Equal(t, flkr.FrameworkRails, profile.Framework)
	assert.Equal(t, "3.2.2", profile.Version)
	assert.True(t, profile.HasLockfile)
}
