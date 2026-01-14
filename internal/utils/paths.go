package utils

import (
	"os"
	"path/filepath"
)

// GetWorkingDir returns the current working directory
func GetWorkingDir() (string, error) {
	return os.Getwd()
}

// GetProjectRoot finds the project root directory
func GetProjectRoot(startPath string) (string, error) {
	return startPath, nil
}

// ResolvePath resolves a relative path to absolute
func ResolvePath(path string) (string, error) {
	return filepath.Abs(path)
}

// CreateDir creates a directory if it doesn't exist
func CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// FindFile searches for a file in a directory
func FindFile(dir string, filename string) (string, error) {
	fullPath := filepath.Join(dir, filename)
	if _, err := os.Stat(fullPath); err == nil {
		return fullPath, nil
	}
	return "", nil
}
