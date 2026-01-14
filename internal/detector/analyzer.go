package detector

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/gaoubak/Makegen/internal/utils"
)

// Result contains all detection results
type Result struct {
	Language        string
	Frameworks      []Framework
	DockerDetected  bool
	DockerServices  []string
	TestDirFound    bool
	BuildDirFound   bool
	HasVendor       bool
	HasModules      bool
	DependencyFiles []string
	ConfigFiles     []string
	MainEntrypoint  string
	ProjectRoot     string
}

// Framework represents a detected framework
type Framework struct {
	Name     string
	Type     string // "web", "cli", "orm", "frontend", etc.
	Files    []string
	Commands map[string]string
	Port     int
	DevTools []string
}

// Analyzer is the main detection engine
type Analyzer struct {
	logger *utils.Logger
}

// NewAnalyzer creates a new analyzer
func NewAnalyzer(logger *utils.Logger) *Analyzer {
	return &Analyzer{
		logger: logger,
	}
}

// Analyze performs complete project analysis
func (a *Analyzer) Analyze(projectPath string) (*Result, error) {
	result := &Result{
		ProjectRoot: projectPath,
	}

	// Detect language
	a.logger.Debug("Detecting language...")
	if err := a.detectLanguage(projectPath, result); err != nil {
		a.logger.Warn("Language detection error: %v", err)
	}

	// Detect frameworks
	a.logger.Debug("Detecting frameworks...")
	if err := a.detectFrameworks(projectPath, result); err != nil {
		a.logger.Warn("Framework detection error: %v", err)
	}

	// Detect Docker
	a.logger.Debug("Detecting Docker...")
	if err := a.detectDocker(projectPath, result); err != nil {
		a.logger.Warn("Docker detection error: %v", err)
	}

	// Analyze project structure
	a.logger.Debug("Analyzing project structure...")
	if err := a.analyzeProjectStructure(projectPath, result); err != nil {
		a.logger.Warn("Project analysis error: %v", err)
	}

	return result, nil
}

// ============================================================================
// FRAMEWORK DETECTION
// ============================================================================

// detectFrameworks detects installed frameworks
func (a *Analyzer) detectFrameworks(path string, result *Result) error {
	result.Frameworks = []Framework{}

	switch result.Language {
	case "go":
		a.detectGoFrameworks(path, result)
	case "javascript", "typescript":
		a.detectJavaScriptFrameworks(path, result)
	case "python":
		a.detectPythonFrameworks(path, result)
	case "rust":
		a.detectRustFrameworks(path, result)
	case "java":
		a.detectJavaFrameworks(path, result)
	case "ruby":
		a.detectRubyFrameworks(path, result)
	}

	return nil
}

// detectGoFrameworks detects Go frameworks
func (a *Analyzer) detectGoFrameworks(path string, result *Result) {
	goModPath := filepath.Join(path, "go.mod")
	content, err := readFile(goModPath)
	if err != nil {
		return
	}

	if hasContent(content, "github.com/gin-gonic/gin") {
		result.Frameworks = append(result.Frameworks, Framework{
			Name: "Gin",
			Type: "web",
			Port: 3000,
		})
	}

	if hasContent(content, "github.com/labstack/echo") {
		result.Frameworks = append(result.Frameworks, Framework{
			Name: "Echo",
			Type: "web",
			Port: 8080,
		})
	}

	if hasContent(content, "github.com/gofiber/fiber") {
		result.Frameworks = append(result.Frameworks, Framework{
			Name: "Fiber",
			Type: "web",
			Port: 3000,
		})
	}

	if hasContent(content, "gorm.io/gorm") {
		result.Frameworks = append(result.Frameworks, Framework{
			Name: "GORM",
			Type: "orm",
		})
	}
}

// detectJavaScriptFrameworks detects JavaScript frameworks
func (a *Analyzer) detectJavaScriptFrameworks(path string, result *Result) {
	packagePath := filepath.Join(path, "package.json")
	content, err := readFile(packagePath)
	if err != nil {
		return
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal([]byte(content), &pkg); err != nil {
		return
	}

	// Get dependencies
	deps := make(map[string]interface{})
	if d, ok := pkg["dependencies"].(map[string]interface{}); ok {
		for k, v := range d {
			deps[k] = v
		}
	}
	if d, ok := pkg["devDependencies"].(map[string]interface{}); ok {
		for k, v := range d {
			deps[k] = v
		}
	}

	// Check for Next.js
	if _, ok := deps["next"]; ok {
		result.Frameworks = append(result.Frameworks, Framework{
			Name: "Next.js",
			Type: "web",
			Port: 3000,
		})
	}

	// Check for React
	if _, ok := deps["react"]; ok {
		result.Frameworks = append(result.Frameworks, Framework{
			Name: "React",
			Type: "frontend",
			Port: 3000,
		})
	}

	// Check for Vue
	if _, ok := deps["vue"]; ok {
		result.Frameworks = append(result.Frameworks, Framework{
			Name: "Vue",
			Type: "frontend",
			Port: 5173,
		})
	}

	// Check for Express
	if _, ok := deps["express"]; ok {
		result.Frameworks = append(result.Frameworks, Framework{
			Name: "Express",
			Type: "web",
			Port: 3000,
		})
	}

	// Check for Fastify
	if _, ok := deps["fastify"]; ok {
		result.Frameworks = append(result.Frameworks, Framework{
			Name: "Fastify",
			Type: "web",
			Port: 3000,
		})
	}

	// Check for NestJS
	if _, ok := deps["@nestjs/core"]; ok {
		result.Frameworks = append(result.Frameworks, Framework{
			Name: "NestJS",
			Type: "web",
			Port: 3000,
		})
	}
}

// detectPythonFrameworks detects Python frameworks
func (a *Analyzer) detectPythonFrameworks(path string, result *Result) {
	// Check requirements.txt
	reqPath := filepath.Join(path, "requirements.txt")
	if content, err := readFile(reqPath); err == nil {
		if hasContent(content, "django") {
			result.Frameworks = append(result.Frameworks, Framework{
				Name: "Django",
				Type: "web",
				Port: 8000,
			})
		}
		if hasContent(content, "flask") {
			result.Frameworks = append(result.Frameworks, Framework{
				Name: "Flask",
				Type: "web",
				Port: 5000,
			})
		}
		if hasContent(content, "fastapi") {
			result.Frameworks = append(result.Frameworks, Framework{
				Name: "FastAPI",
				Type: "web",
				Port: 8000,
			})
		}
		if hasContent(content, "sqlalchemy") {
			result.Frameworks = append(result.Frameworks, Framework{
				Name: "SQLAlchemy",
				Type: "orm",
			})
		}
	}

	// Check pyproject.toml
	pyprojPath := filepath.Join(path, "pyproject.toml")
	if content, err := readFile(pyprojPath); err == nil {
		if hasContent(content, "django") {
			result.Frameworks = append(result.Frameworks, Framework{
				Name: "Django",
				Type: "web",
				Port: 8000,
			})
		}
		if hasContent(content, "flask") {
			result.Frameworks = append(result.Frameworks, Framework{
				Name: "Flask",
				Type: "web",
				Port: 5000,
			})
		}
		if hasContent(content, "fastapi") {
			result.Frameworks = append(result.Frameworks, Framework{
				Name: "FastAPI",
				Type: "web",
				Port: 8000,
			})
		}
	}
}

// detectRustFrameworks detects Rust frameworks
func (a *Analyzer) detectRustFrameworks(path string, result *Result) {
	cargoPath := filepath.Join(path, "Cargo.toml")
	content, err := readFile(cargoPath)
	if err != nil {
		return
	}

	if hasContent(content, "actix-web") {
		result.Frameworks = append(result.Frameworks, Framework{
			Name: "Actix",
			Type: "web",
			Port: 8000,
		})
	}

	if hasContent(content, "rocket") {
		result.Frameworks = append(result.Frameworks, Framework{
			Name: "Rocket",
			Type: "web",
			Port: 8000,
		})
	}

	if hasContent(content, "axum") {
		result.Frameworks = append(result.Frameworks, Framework{
			Name: "Axum",
			Type: "web",
			Port: 8000,
		})
	}
}

// detectJavaFrameworks detects Java frameworks
func (a *Analyzer) detectJavaFrameworks(path string, result *Result) {
	pomPath := filepath.Join(path, "pom.xml")
	if content, err := readFile(pomPath); err == nil {
		if hasContent(content, "spring-boot") {
			result.Frameworks = append(result.Frameworks, Framework{
				Name: "Spring Boot",
				Type: "web",
				Port: 8080,
			})
		}
		return
	}

	gradlePath := filepath.Join(path, "build.gradle")
	if content, err := readFile(gradlePath); err == nil {
		if hasContent(content, "spring-boot") {
			result.Frameworks = append(result.Frameworks, Framework{
				Name: "Spring Boot",
				Type: "web",
				Port: 8080,
			})
		}
	}
}

// detectRubyFrameworks detects Ruby frameworks
func (a *Analyzer) detectRubyFrameworks(path string, result *Result) {
	gemfilePath := filepath.Join(path, "Gemfile")
	content, err := readFile(gemfilePath)
	if err != nil {
		return
	}

	if hasContent(content, "rails") {
		result.Frameworks = append(result.Frameworks, Framework{
			Name: "Rails",
			Type: "web",
			Port: 3000,
		})
	}

	if hasContent(content, "sinatra") {
		result.Frameworks = append(result.Frameworks, Framework{
			Name: "Sinatra",
			Type: "web",
			Port: 4567,
		})
	}
}

// ============================================================================
// DOCKER DETECTION
// ============================================================================

// detectDocker detects Docker configuration
func (a *Analyzer) detectDocker(path string, result *Result) error {
	// Check for Dockerfile
	dockerfilePath := filepath.Join(path, "Dockerfile")
	if fileExists(dockerfilePath) {
		result.DockerDetected = true
		a.logger.Debug("Found Dockerfile")
	}

	// Check for docker-compose.yml
	composePath := filepath.Join(path, "docker-compose.yml")
	if fileExists(composePath) {
		result.DockerDetected = true
		a.logger.Debug("Found docker-compose.yml")
		a.parseDockerCompose(composePath, result)
	}

	// Also check for docker-compose.yaml
	composeYamlPath := filepath.Join(path, "docker-compose.yaml")
	if fileExists(composeYamlPath) {
		result.DockerDetected = true
		a.logger.Debug("Found docker-compose.yaml")
		a.parseDockerCompose(composeYamlPath, result)
	}

	return nil
}

// parseDockerCompose parses docker-compose file to extract services
func (a *Analyzer) parseDockerCompose(path string, result *Result) {
	content, err := readFile(path)
	if err != nil {
		a.logger.Warn("Failed to read docker-compose: %v", err)
		return
	}

	lines := strings.Split(content, "\n")
	inServices := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "services:") {
			inServices = true
			continue
		}

		if inServices {
			if len(line) > 0 && line[0] != ' ' && line[0] != '\t' {
				inServices = false
				continue
			}

			if len(line) >= 2 && line[0] == ' ' && line[1] != ' ' {
				serviceName := strings.TrimSpace(line)
				if strings.HasSuffix(serviceName, ":") {
					serviceName = serviceName[:len(serviceName)-1]
					if serviceName != "" && !strings.Contains(serviceName, " ") {
						result.DockerServices = append(result.DockerServices, serviceName)
						a.logger.Debug("Found Docker service: %s", serviceName)
					}
				}
			}
		}
	}

	result.DockerServices = removeDuplicates(result.DockerServices)
}

// removeDuplicates removes duplicate strings from a slice
func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// ============================================================================
// PROJECT STRUCTURE ANALYSIS
// ============================================================================

// analyzeProjectStructure analyzes project directories
func (a *Analyzer) analyzeProjectStructure(path string, result *Result) error {
	a.findTestDirs(path, result)
	a.findBuildDirs(path, result)
	result.HasVendor = dirExists(filepath.Join(path, "vendor"))
	a.findMainEntrypoint(path, result)
	a.findDependencyFiles(path, result)
	a.findConfigFiles(path, result)
	return nil
}

// findTestDirs finds test directories in the project
func (a *Analyzer) findTestDirs(path string, result *Result) {
	testDirs := []string{
		"test",
		"tests",
		"spec",
		"specs",
		"__tests__",
		".test",
		".tests",
	}

	for _, testDir := range testDirs {
		fullPath := filepath.Join(path, testDir)
		if dirExists(fullPath) {
			result.TestDirFound = true
			a.logger.Debug("Found test directory: %s", testDir)
			return
		}
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.Contains(name, "_test.") || strings.HasSuffix(name, ".test.js") {
			result.TestDirFound = true
			a.logger.Debug("Found test files in root")
			return
		}
	}
}

// findBuildDirs finds build/dist directories
func (a *Analyzer) findBuildDirs(path string, result *Result) {
	buildDirs := []string{
		"build",
		"dist",
		"out",
		"bin",
		"target",
		"release",
		"debug",
		".build",
		"__pycache__",
		"node_modules/.bin",
	}

	for _, buildDir := range buildDirs {
		fullPath := filepath.Join(path, buildDir)
		if dirExists(fullPath) {
			result.BuildDirFound = true
			a.logger.Debug("Found build directory: %s", buildDir)
			return
		}
	}
}

// findMainEntrypoint finds the main entry point of the project
func (a *Analyzer) findMainEntrypoint(path string, result *Result) {
	switch result.Language {
	case "go":
		if fileExists(filepath.Join(path, "main.go")) {
			result.MainEntrypoint = "main.go"
			a.logger.Debug("Found Go entrypoint: main.go")
		}

	case "javascript", "typescript":
		entryPoints := []string{"index.js", "main.js", "app.js", "server.js", "index.ts", "main.ts"}
		for _, ep := range entryPoints {
			if fileExists(filepath.Join(path, ep)) {
				result.MainEntrypoint = ep
				a.logger.Debug("Found JS entrypoint: %s", ep)
				return
			}
		}

		pkgPath := filepath.Join(path, "package.json")
		if content, err := readFile(pkgPath); err == nil && strings.Contains(content, "\"main\"") {
			result.MainEntrypoint = "package.json (main field)"
			a.logger.Debug("Found JS entrypoint in package.json")
		}

	case "python":
		entryPoints := []string{"main.py", "app.py", "__main__.py", "run.py", "wsgi.py"}
		for _, ep := range entryPoints {
			if fileExists(filepath.Join(path, ep)) {
				result.MainEntrypoint = ep
				a.logger.Debug("Found Python entrypoint: %s", ep)
				return
			}
		}

	case "rust":
		if fileExists(filepath.Join(path, "src", "main.rs")) {
			result.MainEntrypoint = "src/main.rs"
			a.logger.Debug("Found Rust entrypoint: src/main.rs")
		}

	case "java":
		if fileExists(filepath.Join(path, "src", "main", "java")) {
			result.MainEntrypoint = "src/main/java"
			a.logger.Debug("Found Java entrypoint directory")
		}

	case "ruby":
		entryPoints := []string{"app.rb", "main.rb", "server.rb", "config.ru"}
		for _, ep := range entryPoints {
			if fileExists(filepath.Join(path, ep)) {
				result.MainEntrypoint = ep
				a.logger.Debug("Found Ruby entrypoint: %s", ep)
				return
			}
		}
	}
}

// findDependencyFiles finds dependency files
func (a *Analyzer) findDependencyFiles(path string, result *Result) {
	depFiles := []string{
		"go.mod",
		"go.sum",
		"package.json",
		"package-lock.json",
		"yarn.lock",
		"pnpm-lock.yaml",
		"requirements.txt",
		"setup.py",
		"pyproject.toml",
		"Pipfile",
		"Gemfile",
		"Gemfile.lock",
		"Cargo.toml",
		"Cargo.lock",
		"pom.xml",
		"build.gradle",
		"composer.json",
		"composer.lock",
	}

	for _, depFile := range depFiles {
		fullPath := filepath.Join(path, depFile)
		if fileExists(fullPath) {
			result.DependencyFiles = append(result.DependencyFiles, depFile)
			a.logger.Debug("Found dependency file: %s", depFile)
		}
	}
}

// findConfigFiles finds configuration files
func (a *Analyzer) findConfigFiles(path string, result *Result) {
	configFiles := []string{
		".env",
		".env.local",
		".env.example",
		"config.yaml",
		"config.yml",
		"config.json",
		".eslintrc",
		".eslintrc.json",
		".prettierrc",
		"jest.config.js",
		"tsconfig.json",
		".pylintrc",
		"setup.cfg",
		"tox.ini",
		".gitignore",
		"Dockerfile",
		"docker-compose.yml",
		"docker-compose.yaml",
		".github/workflows",
		".gitlab-ci.yml",
		".travis.yml",
		"Jenkinsfile",
	}

	for _, configFile := range configFiles {
		fullPath := filepath.Join(path, configFile)
		if fileExists(fullPath) || dirExists(fullPath) {
			result.ConfigFiles = append(result.ConfigFiles, configFile)
			a.logger.Debug("Found config file: %s", configFile)
		}
	}
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// readFile reads file content
func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// hasContent checks if content contains a string (case insensitive)
func hasContent(content, search string) bool {
	return strings.Contains(strings.ToLower(content), strings.ToLower(search))
}
