package detector

import (
	"context"
	"io/fs"

	"github.com/narvanalabs/flkr/pkg/flkr"
)

// Detector analyzes a repository filesystem and returns an AppProfile
// if the ecosystem is detected.
type Detector interface {
	// Name returns a human-readable identifier for this detector.
	Name() string

	// Detect inspects the filesystem and returns a profile if detected.
	// The boolean indicates whether the detector matched at all.
	Detect(ctx context.Context, root fs.FS) (*flkr.AppProfile, bool, error)

	// Priority controls execution order; lower values run first.
	Priority() int
}
