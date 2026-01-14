package app

import (
	"fmt"

	"github.com/gaoubak/Makegen/internal/detector"
	"github.com/gaoubak/Makegen/internal/generator"
	"github.com/gaoubak/Makegen/internal/storage"
	"github.com/gaoubak/Makegen/internal/ui"
	"github.com/gaoubak/Makegen/internal/utils"
)

// App is the main application struct
type App struct {
	logger    *utils.Logger
	workDir   string
	detector  *detector.Analyzer
	storage   storage.FileSystem
	generator *generator.Builder
}

// NewApp creates a new application instance
func NewApp(logger *utils.Logger, workDir string) *App {
	return &App{
		logger:    logger,
		workDir:   workDir,
		detector:  detector.NewAnalyzer(logger),
		storage:   storage.NewLocalFileSystem(logger),
		generator: generator.NewBuilder(logger),
	}
}

// Run executes the main application flow
func (a *App) Run() error {
	a.logger.Info("🔨 Makefile Generator - Interactive Setup")
	a.logger.Info("=====================================\n")

	// Phase 1: Detect Project
	a.logger.Info("📊 Analyzing project...")
	detection, err := a.detector.Analyze(a.workDir)
	if err != nil {
		return fmt.Errorf("detection failed: %w", err)
	}

	a.logDetectionResults(detection)

	// Phase 2: Interactive Questions
	a.logger.Info("\n❓ Configuration Questions")
	a.logger.Info("=======================\n")

	questionnaire := ui.NewQuestionnaire(a.logger, detection)
	config, err := questionnaire.Ask()
	if err != nil {
		return fmt.Errorf("questionnaire failed: %w", err)
	}

	// Phase 3: Generate Makefile
	a.logger.Info("\n📝 Generating Makefile...")
	makefile, err := a.generator.Build(config)
	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	// Phase 4: Preview and Save
	a.logger.Info("\n✨ Preview:")
	a.logger.Info("===========\n")
	fmt.Println(makefile)
	a.logger.Info("\n===========\n")

	// Phase 5: Save to File
	shouldSave := ui.PromptYesNo("Save to Makefile?", true)
	if shouldSave {
		if err := a.storage.WriteMakefile(a.workDir, makefile); err != nil {
			return fmt.Errorf("failed to save Makefile: %w", err)
		}
		a.logger.Success("✅ Makefile saved successfully!")
	} else {
		a.logger.Info("❌ Makefile not saved")
	}

	return nil
}

// logDetectionResults logs what was detected
func (a *App) logDetectionResults(detection *detector.Result) {
	a.logger.Info("✓ Language: %s", detection.Language)
	a.logger.Info("✓ Frameworks found: %d", len(detection.Frameworks))
	a.logger.Info("✓ Docker detected: %v", detection.DockerDetected)
	if detection.DockerDetected {
		a.logger.Info("  Services: %v", detection.DockerServices)
	}
	a.logger.Info("")
}
