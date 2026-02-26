package flkr

// GenerateResult holds the output of flake generation.
type GenerateResult struct {
	// FlakeContent is the rendered flake.nix content.
	FlakeContent string

	// OutputPath is the file path where the flake was written (empty if DryRun).
	OutputPath string
}

// GenerateFunc is set during initialization to avoid import cycles.
var GenerateFunc func(profile *AppProfile, opts GenerateOptions) (*GenerateResult, error)

// Generate renders a flake.nix for the given profile.
func Generate(profile *AppProfile, opts GenerateOptions) (*GenerateResult, error) {
	if err := profile.Validate(); err != nil {
		return nil, err
	}
	if GenerateFunc == nil {
		return nil, nil
	}
	return GenerateFunc(profile, opts)
}
