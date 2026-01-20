package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/drossan/claude-init/internal/ai"
)

// ProjectAnalysis contiene el análisis del proyecto extraído por Claude.
type ProjectAnalysis struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	Language         string `json:"language"`
	Framework        string `json:"framework,omitempty"`
	Architecture     string `json:"architecture"`
	Database         string `json:"database,omitempty"`
	ProjectCategory  string `json:"project_category"`
	BusinessContext  string `json:"business_context"`
	GitSystem        string `json:"git_system,omitempty"`
	TestingFramework string `json:"testing_framework,omitempty"`
}

// Analyzer analiza proyectos existentes usando un Client de IA.
type Analyzer struct {
	projectPath string
	logger      Logger
	client      ai.Client
}

// Logger es la interfaz que debe cumplir el logger.
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
}

// NewAnalyzer crea un nuevo Analyzer.
func NewAnalyzer(projectPath string, client ai.Client) *Analyzer {
	return &Analyzer{
		projectPath: projectPath,
		client:      client,
	}
}

// SetLogger establece el logger.
func (a *Analyzer) SetLogger(logger Logger) {
	a.logger = logger
}

// Analyze ejecuta el análisis del proyecto.
func (a *Analyzer) Analyze() (*ProjectAnalysis, error) {
	// Primero escanear el proyecto localmente
	projectInfo := a.scanProject()

	prompt := a.buildAnalysisPrompt(projectInfo)
	systemPrompt := a.buildSystemPrompt()

	a.logDebug("Analyzing project at: %s", a.projectPath)

	output, err := a.client.SendMessage(systemPrompt, prompt)
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	a.logDebug("AI Response (raw): %s", truncateString(output, 500))

	analysis, err := a.parseAnalysis(output)
	if err != nil {
		return nil, fmt.Errorf("parse failed: %w", err)
	}

	a.logDebug("Analysis completed successfully")
	return analysis, nil
}

// buildSystemPrompt construye el system prompt para Claude.
func (a *Analyzer) buildSystemPrompt() string {
	return `You are an expert software project analyst. Your task is to analyze existing projects and extract structured information about them.

Be precise and concise in your analysis. If unsure about any field, use an empty string or a generic appropriate value.`
}

// buildAnalysisPrompt construye el prompt para Claude.
func (a *Analyzer) buildAnalysisPrompt(projectInfo string) string {
	return fmt.Sprintf(`Analyze this project based on the following information:

%s

Based on this file structure and configuration, identify:
1. Project name (from directory or config files)
2. Main programming language
3. Framework(s) used
4. Architecture type (Monolith, Microservices, etc.)
5. Database if applicable
6. Project category (API, Web App, CLI, Library, etc.)
7. Business context and purpose
8. Testing framework if present

Respond with a JSON object using this exact structure:
{
  "name": "project name",
  "description": "brief project description",
  "language": "main programming language (Go, Python, JavaScript, TypeScript, etc.)",
  "framework": "framework used (optional, or empty string if not applicable)",
  "architecture": "architecture type (Monolith, Microservices, Hexagonal, Layered, etc.)",
  "database": "database used (optional, or empty string if not applicable)",
  "project_category": "project type (REST API, Web App, CLI, Library, etc.)",
  "business_context": "business context and project purpose",
  "git_system": "version control system (git, svn, etc., or empty string)",
  "testing_framework": "testing framework (optional, or empty string)"
}

CRITICAL: Respond with ONLY the raw JSON object. Do not include markdown code blocks, explanations, or any additional text.`, projectInfo)
}

// scanProject escanea el directorio del proyecto y recopila información.
func (a *Analyzer) scanProject() string {
	var info strings.Builder

	// Nombre del proyecto desde el directorio
	projectName := filepath.Base(a.projectPath)
	info.WriteString(fmt.Sprintf("Project directory: %s\n\n", projectName))

	// Detectar sistema de control de versiones
	if _, err := os.Stat(filepath.Join(a.projectPath, ".git")); err == nil {
		info.WriteString("Version control: git\n")
	}

	// Escanear estructura de directorios
	info.WriteString("\nDirectory structure:\n")
	a.scanDirectory(a.projectPath, "", &info, 0)

	// Leer archivos de configuración importantes
	info.WriteString("\n\nConfiguration files:\n")
	configFiles := []string{
		"package.json",
		"go.mod",
		"requirements.txt",
		"pyproject.toml",
		"Cargo.toml",
		"pom.xml",
		"build.gradle",
		"composer.json",
		"Gemfile",
	}

	for _, configFile := range configFiles {
		path := filepath.Join(a.projectPath, configFile)
		if content, err := os.ReadFile(path); err == nil {
			info.WriteString(fmt.Sprintf("\n--- %s ---\n%s\n", configFile, truncateString(string(content), 500)))
		}
	}

	return info.String()
}

// scanDirectory escanea recursivamente un directorio y escribe su estructura.
func (a *Analyzer) scanDirectory(dir, prefix string, info *strings.Builder, depth int) {
	if depth > 5 { // Limitar profundidad
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	// Directorios a ignorar
	ignoreDirs := map[string]bool{
		"node_modules": true,
		"vendor":       true,
		".git":         true,
		"dist":         true,
		"build":        true,
		"target":       true,
		"bin":          true,
		"obj":          true,
		".venv":        true,
		"venv":         true,
		"__pycache__":  true,
		".claude":      true,
		".idea":        true,
		".vscode":      true,
	}

	dirs := []string{}
	files := []string{}

	for _, entry := range entries {
		name := entry.Name()
		if ignoreDirs[name] {
			continue
		}

		if entry.IsDir() {
			dirs = append(dirs, name)
		} else {
			files = append(files, name)
		}
	}

	// Mostrar primero los directorios
	for _, dirName := range dirs {
		fullPath := filepath.Join(dir, dirName)
		info.WriteString(fmt.Sprintf("%s%s/\n", prefix, dirName))
		a.scanDirectory(fullPath, prefix+"  ", info, depth+1)
	}

	// Luego mostrar archivos (limitados)
	for i, fileName := range files {
		if i >= 10 { // Limitar a 10 archivos por directorio
			info.WriteString(fmt.Sprintf("%s... (%d more files)\n", prefix, len(files)-10))
			break
		}
		info.WriteString(fmt.Sprintf("%s%s\n", prefix, fileName))
	}
}

// parseAnalysis extrae y parsea el JSON de la respuesta.
func (a *Analyzer) parseAnalysis(output string) (*ProjectAnalysis, error) {
	jsonStr := a.extractJSON(output)
	if jsonStr == "" {
		return nil, fmt.Errorf("no JSON found in response")
	}

	var analysis ProjectAnalysis
	if err := json.Unmarshal([]byte(jsonStr), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Usar valores por defecto en lugar de fallar
	if analysis.Name == "" {
		analysis.Name = "Unknown Project"
		a.logDebug("Missing field 'name', using default: 'Unknown Project'")
	}
	if analysis.Language == "" {
		analysis.Language = "Unknown"
		a.logDebug("Missing field 'language', using default: 'Unknown'")
	}
	if analysis.Architecture == "" {
		analysis.Architecture = "Monolith"
		a.logDebug("Missing field 'architecture', using default: 'Monolith'")
	}
	if analysis.ProjectCategory == "" {
		analysis.ProjectCategory = "General"
		a.logDebug("Missing field 'project_category', using default: 'General'")
	}
	if analysis.BusinessContext == "" {
		analysis.BusinessContext = "General purpose software project"
		a.logDebug("Missing field 'business_context', using default")
	}

	a.logDebug("Parsed analysis: name=%s, language=%s, category=%s", analysis.Name, analysis.Language, analysis.ProjectCategory)
	return &analysis, nil
}

// extractJSON extrae el objeto JSON de una respuesta mixta.
func (a *Analyzer) extractJSON(s string) string {
	s = strings.TrimSpace(s)

	// Buscar el primer { y el último } para encontrar el JSON completo
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")

	if start == -1 || end == -1 || end <= start {
		a.logDebug("No JSON found in response. Response preview: %s", truncateString(s, 200))
		return ""
	}

	jsonStr := s[start : end+1]

	// Validar que el JSON extraído es válido intentando parsearlo
	var test interface{}
	if err := json.Unmarshal([]byte(jsonStr), &test); err != nil {
		a.logDebug("Extracted JSON is invalid: %v. Trying to fix...", err)
		// Si no es válido, intentar encontrar bloques de código markdown
		jsonStr = a.extractJSONFromMarkdown(s)
	}

	return jsonStr
}

// extractJSONFromMarkdown intenta extraer JSON de bloques de código markdown.
func (a *Analyzer) extractJSONFromMarkdown(s string) string {
	// Buscar ```json ... ```
	patterns := []string{
		"```json",
		"```JSON",
		"```",
	}

	for _, pattern := range patterns {
		startIdx := strings.Index(s, pattern)
		if startIdx != -1 {
			// Saltar el pattern
			startIdx += len(pattern)
			// Encontrar el fin del bloque
			endIdx := strings.Index(s[startIdx:], "```")
			if endIdx != -1 {
				candidate := strings.TrimSpace(s[startIdx : startIdx+endIdx])
				// Validar que sea JSON válido
				var test interface{}
				if err := json.Unmarshal([]byte(candidate), &test); err == nil {
					return candidate
				}
			}
		}
	}

	return ""
}

// truncateString corta un string a una longitud máxima.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// logDebug logs debug messages if logger is set.
func (a *Analyzer) logDebug(format string, args ...interface{}) {
	if a.logger != nil {
		a.logger.Debug(format, args...)
	}
}
