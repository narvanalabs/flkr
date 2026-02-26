package detector

import (
	"io/fs"
)

// fileExists checks whether a file exists in the given filesystem.
func fileExists(root fs.FS, path string) bool {
	_, err := fs.Stat(root, path)
	return err == nil
}

// readFile reads the entire contents of a file from the filesystem.
func readFile(root fs.FS, path string) ([]byte, error) {
	return fs.ReadFile(root, path)
}

// readFileString reads a file as a string, returning empty string on error.
func readFileString(root fs.FS, path string) string {
	data, err := readFile(root, path)
	if err != nil {
		return ""
	}
	return string(data)
}
