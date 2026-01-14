package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gaoubak/Makegen/internal/config"
	"github.com/gaoubak/Makegen/internal/detector"
	"github.com/gaoubak/Makegen/internal/utils"
)

// Questionnaire manages the interactive question flow
type Questionnaire struct {
	logger    *utils.Logger
	detection *detector.Result
	reader    *bufio.Reader
	config    *config.MakefileConfig
}

// NewQuestionnaire creates a new questionnaire
func NewQuestionnaire(logger *utils.Logger, detection *detector.Result) *Questionnaire {
	return &Questionnaire{
		logger:    logger,
		detection: detection,
		reader:    bufio.NewReader(os.Stdin),
		config:    config.NewMakefileConfig(),
	}
}

// Ask runs the interactive questionnaire
func (q *Questionnaire) Ask() (*config.MakefileConfig, error) {
	// Phase 1: Project Info
	q.askProjectName()
	q.askFramework()

	// Phase 2: Docker
	q.askDocker()

	// Phase 3: Build & Test
	q.askBuildTargets()
	q.askTestSetup()

	// Phase 4: Quality
	q.askLinting()
	q.askFormatting()

	// Phase 5: Advanced
	q.askCICD()
	q.askDeployment()
	q.askCustomTargets()

	return q.config, nil
}

// Helper prompts
func (q *Questionnaire) askProjectName() {
	fmt.Print("\n📝 Project name: ")
	name, _ := q.reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name != "" {
		q.config.ProjectName = name
	} else {
		q.config.ProjectName = "myproject"
	}
	q.logger.Info("✓ Project: %s", q.config.ProjectName)
}

func (q *Questionnaire) askFramework() {
	if len(q.detection.Frameworks) == 0 {
		return
	}

	fmt.Println("\n🎯 Detected Frameworks:")
	for i, fw := range q.detection.Frameworks {
		fmt.Printf("  %d. %s (%s)\n", i+1, fw.Name, fw.Type)
	}

	if PromptYesNo("Use a detected framework?", true) {
		// TODO: Implement framework selection
	}
}

func (q *Questionnaire) askDocker() {
	if !q.detection.DockerDetected {
		if PromptYesNo("\n🐳 Add Docker support?", false) {
			q.config.HasDocker = true
		}
		return
	}

	fmt.Println("\n🐳 Docker detected!")
	if len(q.detection.DockerServices) > 0 {
		fmt.Printf("   Services: %v\n", q.detection.DockerServices)
	}

	if PromptYesNo("Add Docker targets?", true) {
		q.config.HasDocker = true
		q.config.DockerServices = q.detection.DockerServices

		fmt.Print("Docker image name: ")
		name, _ := q.reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if name != "" {
			q.config.DockerImage = name
		}

		if PromptYesNo("Add docker-compose targets?", true) {
			q.config.DockerCompose = true
		}
	}
}

func (q *Questionnaire) askBuildTargets() {
	fmt.Println("\n🔨 Build Configuration")

	// TODO: Language-specific build targets
	if PromptYesNo("Add 'build' target?", true) {
		// Add build target
	}
	if PromptYesNo("Add 'clean' target?", true) {
		// Add clean target
	}
	if PromptYesNo("Add 'run' target?", true) {
		// Add run target
	}
}

func (q *Questionnaire) askTestSetup() {
	fmt.Println("\n🧪 Testing Configuration")

	if !q.detection.TestDirFound {
		if !PromptYesNo("No test directory found. Add test target anyway?", false) {
			return
		}
	}

	if PromptYesNo("Add 'test' target?", true) {
		// TODO: Test framework selection
		if PromptYesNo("Add coverage target?", true) {
			// Add coverage target
		}
	}
}

func (q *Questionnaire) askLinting() {
	fmt.Println("\n🔍 Linting Configuration")

	if PromptYesNo("Add 'lint' target?", true) {
		// TODO: Linter selection
	}
}

func (q *Questionnaire) askFormatting() {
	fmt.Println("\n✨ Code Formatting")

	if PromptYesNo("Add 'format' target?", true) {
		// TODO: Formatter selection
	}
}

func (q *Questionnaire) askCICD() {
	fmt.Println("\n🔄 CI/CD Configuration")

	if PromptYesNo("Add GitHub Actions CI target?", false) {
		q.config.EnableCI = true
		// TODO: CI/CD configuration
	}
}

func (q *Questionnaire) askDeployment() {
	fmt.Println("\n🚀 Deployment Configuration")

	if PromptYesNo("Add deployment targets?", false) {
		q.config.EnableDeploy = true
		// TODO: Deployment target selection
	}
}

func (q *Questionnaire) askCustomTargets() {
	fmt.Println("\n✨ Custom Targets")

	for PromptYesNo("Add custom target?", false) {
		// TODO: Custom target input
	}
}

// PromptYesNo asks a yes/no question
func PromptYesNo(message string, defaultYes bool) bool {
	suffix := "[Y/n]"
	if !defaultYes {
		suffix = "[y/N]"
	}
	fmt.Printf("%s %s: ", message, suffix)

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	if response == "" {
		return defaultYes
	}
	return response == "y" || response == "yes"
}
