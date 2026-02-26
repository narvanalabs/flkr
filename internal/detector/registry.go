package detector

import (
	"context"
	"io/fs"
	"os"
	"sort"

	"github.com/narvanalabs/flkr/pkg/flkr"
)

// Registry manages all available detectors and orchestrates detection.
type Registry struct {
	detectors []Detector
}

// NewRegistry creates a registry with all built-in detectors.
func NewRegistry() *Registry {
	return &Registry{
		detectors: []Detector{
			&NodeDetector{},
			&PythonDetector{},
			&GoDetector{},
			&RustDetector{},
			&RubyDetector{},
			&ElixirDetector{},
			&PHPDetector{},
			&JavaDetector{},
		},
	}
}

// DetectAll runs every detector and returns all matching profiles, sorted
// by confidence (highest first).
func (r *Registry) DetectAll(ctx context.Context, root fs.FS) ([]*flkr.AppProfile, error) {
	// Sort detectors by priority.
	sorted := make([]Detector, len(r.detectors))
	copy(sorted, r.detectors)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority() < sorted[j].Priority()
	})

	var profiles []*flkr.AppProfile
	for _, d := range sorted {
		profile, matched, err := d.Detect(ctx, root)
		if err != nil {
			return nil, err
		}
		if matched && profile != nil {
			profiles = append(profiles, profile)
		}
	}

	// Sort by confidence descending.
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Confidence > profiles[j].Confidence
	})

	return profiles, nil
}

// DetectBest runs all detectors and returns the highest-confidence match,
// enriched with cross-cutting data.
func (r *Registry) DetectBest(ctx context.Context, root fs.FS) (*flkr.AppProfile, error) {
	profiles, err := r.DetectAll(ctx, root)
	if err != nil {
		return nil, err
	}
	if len(profiles) == 0 {
		return nil, nil
	}

	best := profiles[0]

	// Enrich with cross-cutting data.
	cc := &CrosscuttingDetector{}
	enrichment, matched, err := cc.Detect(ctx, root)
	if err != nil {
		return nil, err
	}
	if matched {
		best.Merge(enrichment)
	}

	return best, nil
}

// DetectFromPath is a convenience that opens an OS directory and runs DetectBest.
func (r *Registry) DetectFromPath(ctx context.Context, path string) (*flkr.AppProfile, error) {
	return r.DetectBest(ctx, os.DirFS(path))
}
