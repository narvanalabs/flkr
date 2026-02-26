package flkr

import (
	"context"
	"fmt"
	"os"
)

// DetectFunc is the function signature that the internal detector registry
// exposes. It is set during initialization by the cmd package to avoid
// an import cycle.
var DetectFunc func(ctx context.Context, root string) (*AppProfile, error)

// Detect scans the repository at the given path and returns an AppProfile.
func Detect(ctx context.Context, opts DetectOptions) (*AppProfile, error) {
	if opts.Path == "" {
		opts.Path = "."
	}
	info, err := os.Stat(opts.Path)
	if err != nil {
		return nil, fmt.Errorf("cannot access path %q: %w", opts.Path, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path %q is not a directory", opts.Path)
	}
	if DetectFunc == nil {
		return nil, fmt.Errorf("detection engine not initialized")
	}
	return DetectFunc(ctx, opts.Path)
}
