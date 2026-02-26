// Package nixhash computes Nix-compatible hashes for build inputs.
package nixhash

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GoVendorHash computes the vendorHash for buildGoModule by running
// `go mod vendor` into a temp directory and hashing it with `nix hash path`.
// Returns the SRI hash string (e.g. "sha256-abc...=") or an error.
func GoVendorHash(projectDir string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "flkr-vendor-*")
	if err != nil {
		return "", fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	vendorDir := filepath.Join(tmpDir, "vendor")

	// Run go mod vendor into the temp directory.
	goCmd := exec.Command("go", "mod", "vendor", "-o", vendorDir)
	goCmd.Dir = projectDir
	goCmd.Env = append(os.Environ(),
		"GOFLAGS=-mod=mod",
		"GOPATH="+filepath.Join(tmpDir, "gopath"),
		"GOMODCACHE="+filepath.Join(tmpDir, "gomodcache"),
	)
	var stderr bytes.Buffer
	goCmd.Stderr = &stderr
	if err := goCmd.Run(); err != nil {
		return "", fmt.Errorf("go mod vendor: %s: %w", stderr.String(), err)
	}

	// Hash the vendor directory with nix hash path (produces SRI format).
	nixCmd := exec.Command("nix", "hash", "path", vendorDir)
	var out bytes.Buffer
	nixCmd.Stdout = &out
	nixCmd.Stderr = &stderr
	if err := nixCmd.Run(); err != nil {
		return "", fmt.Errorf("nix hash path: %s: %w", stderr.String(), err)
	}

	hash := strings.TrimSpace(out.String())
	if hash == "" {
		return "", fmt.Errorf("nix hash path returned empty output")
	}

	return hash, nil
}
