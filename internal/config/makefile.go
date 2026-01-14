package config

// MakefileConfig represents the complete Makefile configuration
type MakefileConfig struct {
	ProjectName    string
	Language       string
	Framework      *FrameworkConfig
	HasDocker      bool
	DockerImage    string
	DockerServices []string
	DockerCompose  bool
	EnableCI       bool
	EnableDeploy   bool
	BuildTools     []string
	TestFramework  string
	LintTools      []string
	FormatTools    []string
	CustomTargets  map[string]Target
}

// FrameworkConfig represents a selected framework
type FrameworkConfig struct {
	Name     string
	Type     string
	Commands map[string]string
	Port     int
}

// Target represents a Makefile target
type Target struct {
	Name         string
	Dependencies []string
	Commands     []string
	Description  string
	Phony        bool
}

// Variable represents a Makefile variable
type Variable struct {
	Name  string
	Value string
}

// NewMakefileConfig creates a new configuration
func NewMakefileConfig() *MakefileConfig {
	return &MakefileConfig{
		ProjectName:    "myproject",
		CustomTargets:  make(map[string]Target),
		BuildTools:     []string{},
		LintTools:      []string{},
		FormatTools:    []string{},
		DockerServices: []string{},
	}
}

// NewTarget creates a new target
func NewTarget(name string) *Target {
	return &Target{
		Name:         name,
		Dependencies: []string{},
		Commands:     []string{},
		Phony:        true,
	}
}

// AddCommand adds a command to a target
func (t *Target) AddCommand(cmd string) {
	t.Commands = append(t.Commands, cmd)
}

// AddDependency adds a dependency to a target
func (t *Target) AddDependency(dep string) {
	t.Dependencies = append(t.Dependencies, dep)
}
