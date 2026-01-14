package app

// Config represents application configuration
type Config struct {
	WorkDir     string
	Verbose     bool
	ProjectName string
}

// NewConfig creates a new app configuration
func NewConfig(workDir string, verbose bool) *Config {
	return &Config{
		WorkDir: workDir,
		Verbose: verbose,
	}
}
