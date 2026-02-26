package parser

import (
	"io/fs"
	"strings"
)

// ParseEnvFile extracts variable names from a .env-style file.
func ParseEnvFile(root fs.FS, path string) []string {
	data, err := fs.ReadFile(root, path)
	if err != nil {
		return nil
	}
	var keys []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if k, _, ok := strings.Cut(line, "="); ok {
			k = strings.TrimSpace(k)
			if k != "" {
				keys = append(keys, k)
			}
		}
	}
	return keys
}
