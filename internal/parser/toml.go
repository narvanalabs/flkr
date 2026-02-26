package parser

import (
	"io/fs"

	"github.com/BurntSushi/toml"
)

// PyprojectTOML represents a Python pyproject.toml file.
type PyprojectTOML struct {
	Project struct {
		Name           string   `toml:"name"`
		Version        string   `toml:"version"`
		RequiresPython string   `toml:"requires-python"`
		Dependencies   []string `toml:"dependencies"`
	} `toml:"project"`
	Tool struct {
		Poetry struct {
			Name         string            `toml:"name"`
			Dependencies map[string]any    `toml:"dependencies"`
			Scripts      map[string]string `toml:"scripts"`
		} `toml:"poetry"`
	} `toml:"tool"`
}

// HasDep checks if a dependency name appears in the project or poetry deps.
func (p *PyprojectTOML) HasDep(name string) bool {
	for _, d := range p.Project.Dependencies {
		// Dependencies can be "name>=version" etc.
		if len(d) >= len(name) && d[:len(name)] == name {
			if len(d) == len(name) || d[len(name)] == '>' || d[len(name)] == '<' || d[len(name)] == '=' || d[len(name)] == '!' || d[len(name)] == '[' || d[len(name)] == ';' {
				return true
			}
		}
	}
	if _, ok := p.Tool.Poetry.Dependencies[name]; ok {
		return true
	}
	return false
}

// ParsePyprojectTOML reads and parses a pyproject.toml.
func ParsePyprojectTOML(root fs.FS, path string) (*PyprojectTOML, error) {
	data, err := fs.ReadFile(root, path)
	if err != nil {
		return nil, err
	}
	var proj PyprojectTOML
	if err := toml.Unmarshal(data, &proj); err != nil {
		return nil, err
	}
	return &proj, nil
}

// CargoTOML represents a Rust Cargo.toml file.
type CargoTOML struct {
	Package struct {
		Name    string `toml:"name"`
		Version string `toml:"version"`
		Edition string `toml:"edition"`
	} `toml:"package"`
	Dependencies map[string]any `toml:"dependencies"`
}

// HasDep checks if a cargo dependency exists.
func (c *CargoTOML) HasDep(name string) bool {
	_, ok := c.Dependencies[name]
	return ok
}

// ParseCargoTOML reads and parses a Cargo.toml.
func ParseCargoTOML(root fs.FS, path string) (*CargoTOML, error) {
	data, err := fs.ReadFile(root, path)
	if err != nil {
		return nil, err
	}
	var cargo CargoTOML
	if err := toml.Unmarshal(data, &cargo); err != nil {
		return nil, err
	}
	return &cargo, nil
}

// RustToolchainTOML represents a rust-toolchain.toml file.
type RustToolchainTOML struct {
	Toolchain struct {
		Channel string `toml:"channel"`
	} `toml:"toolchain"`
}

// ParseRustToolchainTOML reads and parses a rust-toolchain.toml.
func ParseRustToolchainTOML(root fs.FS, path string) (*RustToolchainTOML, error) {
	data, err := fs.ReadFile(root, path)
	if err != nil {
		return nil, err
	}
	var tc RustToolchainTOML
	if err := toml.Unmarshal(data, &tc); err != nil {
		return nil, err
	}
	return &tc, nil
}
