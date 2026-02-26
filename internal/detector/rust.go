package detector

import (
	"context"
	"io/fs"

	"github.com/narvanalabs/flkr/internal/parser"
	"github.com/narvanalabs/flkr/pkg/flkr"
)

// RustDetector detects Rust applications.
type RustDetector struct{}

func (d *RustDetector) Name() string  { return "rust" }
func (d *RustDetector) Priority() int { return 40 }

func (d *RustDetector) Detect(ctx context.Context, root fs.FS) (*flkr.AppProfile, bool, error) {
	if !fileExists(root, "Cargo.toml") {
		return nil, false, nil
	}

	profile := &flkr.AppProfile{
		Language:       flkr.LangRust,
		PackageManager: flkr.PkgCargo,
		Confidence:     0.8,
		DetectedBy:     d.Name(),
		BuildCommand:   "cargo build --release",
		StartCommand:   "./target/release/app",
		Port:           8080,
	}

	if fileExists(root, "Cargo.lock") {
		profile.HasLockfile = true
		profile.LockfileType = "cargo"
	}

	// Parse Cargo.toml for edition and deps.
	cargo, err := parser.ParseCargoTOML(root, "Cargo.toml")
	if err == nil {
		if cargo.Package.Edition != "" {
			profile.Version = cargo.Package.Edition
		}
		if cargo.Package.Name != "" {
			profile.StartCommand = "./target/release/" + cargo.Package.Name
		}
		if cargo.HasDep("actix-web") {
			profile.Framework = flkr.FrameworkActix
			profile.Confidence = 0.9
		}
	}

	// Check rust-toolchain.toml.
	tc, err := parser.ParseRustToolchainTOML(root, "rust-toolchain.toml")
	if err == nil && tc.Toolchain.Channel != "" {
		profile.Version = tc.Toolchain.Channel
	}

	return profile, true, nil
}
