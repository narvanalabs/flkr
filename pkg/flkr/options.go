package flkr

// DetectOptions configures detection behavior.
type DetectOptions struct {
	// Path is the root directory to scan. Defaults to ".".
	Path string

	// Verbose enables detailed detection logging.
	Verbose bool
}

// GenerateOptions configures flake generation behavior.
type GenerateOptions struct {
	// OutputPath is where to write the flake.nix. Defaults to "<detect-path>/flake.nix".
	OutputPath string

	// TemplateVersion pins the flkr-templates revision. Empty means "main".
	TemplateVersion string

	// DryRun prints the generated flake to stdout instead of writing a file.
	DryRun bool
}
