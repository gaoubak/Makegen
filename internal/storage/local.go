package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gaoubak/Makegen/internal/utils"
)

// FileSystem interface defines file operations
type FileSystem interface {
	WriteMakefile(dir, content string) error
	ReadMakefile(dir string) (string, error)
	FileExists(path string) bool
	ListFiles(dir string, extensions []string) ([]string, error)
}

// LocalFileSystem implements FileSystem using local filesystem
type LocalFileSystem struct {
	logger *utils.Logger
}

// NewLocalFileSystem creates a new local filesystem
func NewLocalFileSystem(logger *utils.Logger) *LocalFileSystem {
	return &LocalFileSystem{
		logger: logger,
	}
}

// WriteMakefile writes the Makefile to disk
func (lfs *LocalFileSystem) WriteMakefile(dir, content string) error {
	makefilePath := filepath.Join(dir, "Makefile")

	// Check if file exists
	if lfs.FileExists(makefilePath) {
		lfs.logger.Warn("Makefile already exists at %s", makefilePath)
		lfs.logger.Warn("It will be overwritten")
	}

	// Write file
	err := os.WriteFile(makefilePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write Makefile: %w", err)
	}

	lfs.logger.Info("Makefile written to %s", makefilePath)
	return nil
}

// ReadMakefile reads an existing Makefile
func (lfs *LocalFileSystem) ReadMakefile(dir string) (string, error) {
	makefilePath := filepath.Join(dir, "Makefile")

	if !lfs.FileExists(makefilePath) {
		return "", fmt.Errorf("Makefile not found at %s", makefilePath)
	}

	content, err := os.ReadFile(makefilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read Makefile: %w", err)
	}

	return string(content), nil
}

// FileExists checks if a file exists
func (lfs *LocalFileSystem) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ListFiles lists files in a directory with given extensions
func (lfs *LocalFileSystem) ListFiles(dir string, extensions []string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		for _, ext := range extensions {
			if filepath.Ext(filename) == ext {
				fullPath := filepath.Join(dir, filename)
				files = append(files, fullPath)
				break
			}
		}
	}

	return files, nil
}
