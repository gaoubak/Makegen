package detector

import (
	"os"
	"path/filepath"
)

func (a *Analyzer) detectLanguage(path string, result *Result) error {
	// Check for Go
	if fileExists(filepath.Join(path, "go.mod")) {
		result.Language = "go"
		result.HasModules = true
		return nil
	}

	// Check for Python
	if fileExists(filepath.Join(path, "requirements.txt")) ||
		fileExists(filepath.Join(path, "setup.py")) ||
		fileExists(filepath.Join(path, "pyproject.toml")) {
		result.Language = "python"
		return nil
	}

	// Check for Node.js/JavaScript/TypeScript
	if fileExists(filepath.Join(path, "package.json")) {
		// Check if TypeScript
		if fileExists(filepath.Join(path, "tsconfig.json")) {
			result.Language = "typescript"
		} else {
			result.Language = "javascript"
		}
		result.HasModules = true
		return nil
	}

	// Check for Rust
	if fileExists(filepath.Join(path, "Cargo.toml")) {
		result.Language = "rust"
		return nil
	}

	// Check for Java
	if fileExists(filepath.Join(path, "pom.xml")) {
		result.Language = "java"
		return nil
	}

	if fileExists(filepath.Join(path, "build.gradle")) ||
		fileExists(filepath.Join(path, "build.gradle.kts")) {
		result.Language = "java"
		return nil
	}

	// Check for Ruby
	if fileExists(filepath.Join(path, "Gemfile")) {
		result.Language = "ruby"
		return nil
	}

	// Check for PHP
	if fileExists(filepath.Join(path, "composer.json")) {
		result.Language = "php"
		return nil
	}

	// Check for C/C++
	if fileExists(filepath.Join(path, "CMakeLists.txt")) ||
		fileExists(filepath.Join(path, "Makefile")) {
		result.Language = "cpp"
		return nil
	}

	// Default: unknown
	result.Language = "unknown"
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
