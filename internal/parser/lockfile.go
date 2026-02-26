package parser

import "io/fs"

// LockfileInfo describes a detected lockfile.
type LockfileInfo struct {
	Type string // e.g. "package-lock.json", "yarn.lock", "pnpm-lock.yaml"
	Path string
}

// lockfiles maps lockfile names to their type descriptions.
var lockfiles = map[string]string{
	"package-lock.json": "npm",
	"yarn.lock":         "yarn",
	"pnpm-lock.yaml":   "pnpm",
	"Pipfile.lock":     "pipenv",
	"poetry.lock":      "poetry",
	"uv.lock":          "uv",
	"go.sum":           "gomod",
	"Cargo.lock":       "cargo",
	"Gemfile.lock":     "bundler",
	"mix.lock":         "mix",
	"composer.lock":    "composer",
}

// DetectLockfile checks for known lockfiles in the root.
func DetectLockfile(root fs.FS) *LockfileInfo {
	for path, typ := range lockfiles {
		if _, err := fs.Stat(root, path); err == nil {
			return &LockfileInfo{Type: typ, Path: path}
		}
	}
	return nil
}
