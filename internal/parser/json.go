package parser

import (
	"encoding/json"
	"io/fs"
)

// PackageJSON represents a Node.js package.json file.
type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Engines         struct {
		Node string `json:"node"`
	} `json:"engines"`
}

// HasDep returns true if the package has the given dependency (prod or dev).
func (p *PackageJSON) HasDep(name string) bool {
	if _, ok := p.Dependencies[name]; ok {
		return true
	}
	if _, ok := p.DevDependencies[name]; ok {
		return true
	}
	return false
}

// ParsePackageJSON reads and parses a package.json from the given fs.
func ParsePackageJSON(root fs.FS, path string) (*PackageJSON, error) {
	data, err := fs.ReadFile(root, path)
	if err != nil {
		return nil, err
	}
	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}
	return &pkg, nil
}

// ComposerJSON represents a PHP composer.json file.
type ComposerJSON struct {
	Name    string            `json:"name"`
	Version string            `json:"version"`
	Require map[string]string `json:"require"`
	Scripts map[string]any `json:"scripts"`
	Extra   map[string]any    `json:"extra"`
}

// HasRequire checks if a composer package is required.
func (c *ComposerJSON) HasRequire(name string) bool {
	_, ok := c.Require[name]
	return ok
}

// ParseComposerJSON reads and parses a composer.json from the given fs.
func ParseComposerJSON(root fs.FS, path string) (*ComposerJSON, error) {
	data, err := fs.ReadFile(root, path)
	if err != nil {
		return nil, err
	}
	var comp ComposerJSON
	if err := json.Unmarshal(data, &comp); err != nil {
		return nil, err
	}
	return &comp, nil
}
