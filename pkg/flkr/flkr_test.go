package flkr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppProfile_Validate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		p := &AppProfile{
			Language:       LangNode,
			PackageManager: PkgNPM,
			Confidence:     0.9,
		}
		require.NoError(t, p.Validate())
	})

	t.Run("missing language", func(t *testing.T) {
		p := &AppProfile{PackageManager: PkgNPM, Confidence: 0.5}
		assert.ErrorContains(t, p.Validate(), "language is required")
	})

	t.Run("missing package manager", func(t *testing.T) {
		p := &AppProfile{Language: LangNode, Confidence: 0.5}
		assert.ErrorContains(t, p.Validate(), "packageManager is required")
	})

	t.Run("invalid confidence", func(t *testing.T) {
		p := &AppProfile{Language: LangNode, PackageManager: PkgNPM, Confidence: 1.5}
		assert.ErrorContains(t, p.Validate(), "confidence must be between")
	})
}

func TestAppProfile_Merge(t *testing.T) {
	base := &AppProfile{
		Language:       LangNode,
		PackageManager: PkgNPM,
		Port:           3000,
		Confidence:     0.7,
		EnvVars:        []string{"DB_URL"},
	}

	other := &AppProfile{
		Framework:  FrameworkNextJS,
		Confidence: 0.9,
		EnvVars:    []string{"DB_URL", "SECRET"},
	}

	base.Merge(other)
	assert.Equal(t, LangNode, base.Language)
	assert.Equal(t, FrameworkNextJS, base.Framework)
	assert.InDelta(t, 0.9, base.Confidence, 0.01)
	assert.Equal(t, 3000, base.Port)
	assert.Equal(t, []string{"DB_URL", "SECRET"}, base.EnvVars)
}

func TestAppProfile_Merge_Nil(t *testing.T) {
	p := &AppProfile{Language: LangNode}
	p.Merge(nil)
	assert.Equal(t, LangNode, p.Language)
}
