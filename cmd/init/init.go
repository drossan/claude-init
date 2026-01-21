// Package init implementa el comando init para inicializar proyectos.
//
// Este comando ejecuta un survey interactivo y genera la estructura .claude/
// usando Claude CLI directamente.
package init

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	gSurvey "github.com/AlecAivazis/survey/v2"
	"github.com/drossan/claude-init/internal/ai"
	aifactory "github.com/drossan/claude-init/internal/ai"
	"github.com/drossan/claude-init/internal/claude"
	"github.com/drossan/claude-init/internal/config"
	"github.com/drossan/claude-init/internal/logger"
	"github.com/drossan/claude-init/internal/survey"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	// DefaultConfigDir es el directorio de configuración por defecto.
	DefaultConfigDir = ".claude"
)

var (
	// log es el logger compartido con el comando root.
	log *logger.Logger
	// initCmd es el comando init.
	initCmd *cobra.Command
)

// InitOptions contiene las opciones del comando init.
type InitOptions struct {
	// Force indica si se deben sobrescribir archivos existentes.
	Force bool
	// DryRun indica si solo se debe mostrar qué se generaría sin crear archivos.
	DryRun bool
	// ConfigDir es el directorio de configuración (default: .claude).
	ConfigDir string
}

// Execute añade el comando init al root command.
func Execute(rootCmd *cobra.Command, l *logger.Logger) {
	log = l
	initCmd = NewInitCommand()
	rootCmd.AddCommand(initCmd)
}

// GetInitCmd retorna el comando init para testing.
func GetInitCmd() *cobra.Command {
	return initCmd
}

// NewInitCommand crea un nuevo comando init con todos los flags configurados.
//
// El comando init ejecuta un survey interactivo, usa Claude CLI para generar
// recomendaciones, y crea la estructura .claude/ basada en las respuestas.
func NewInitCommand() *cobra.Command {
	opts := &InitOptions{
		ConfigDir: DefaultConfigDir,
	}

	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Initialize .claude/ directory for AI-assisted development",
		Long: `Initialize .claude/ directory for AI-assisted development.

This command runs an interactive survey to understand your project,
uses Claude CLI to generate recommendations, and creates
the .claude/ structure with agents, commands, and skills.

Requires Claude CLI to be installed. Visit: https://claude.com/claude-code

If no path is provided, the current directory is used.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(cmd, opts, args)
		},
	}

	// Configurar flags (eliminados flags de IA: --ai-provider, --api-key, --no-ai)
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Overwrite existing files")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Show what would be generated without creating files")
	cmd.Flags().StringVar(&opts.ConfigDir, "config-dir", DefaultConfigDir, "Config directory name")

	return cmd
}

// runInit ejecuta el comando init con el flujo simplificado usando Claude CLI.
func runInit(cmd *cobra.Command, opts *InitOptions, args []string) error {
	// Inicializar logger si es nil
	if log == nil {
		log = logger.New(os.Stdout, logger.INFOLevel)
	}

	// 1. Determinar el path del proyecto
	projectPath, err := getProjectPath(args)
	if err != nil {
		return fmt.Errorf("failed to determine project path: %w", err)
	}

	// 2. Validar que el path existe
	if err := validatePath(projectPath); err != nil {
		return err
	}

	// 3. Verificar si ya existe configuración
	if err := checkExistingConfig(projectPath, opts.ConfigDir, opts.Force); err != nil {
		return err
	}

	// 4. PREGUNTAR POR PROVIDER DE IA
	log.Info("\nAI Provider Selection")
	aiProvider, err := askAIProvider()
	if err != nil {
		return fmt.Errorf("failed to ask AI provider: %w", err)
	}

	// Crear cliente según provider seleccionado
	factory := aifactory.NewClientFactory()
	client, err := factory.CreateClientFromString(aiProvider)
	if err != nil {
		// Si el error es por falta de configuración, pedirla interactivamente
		if strings.Contains(err.Error(), "not configured") {
			log.Info("AI provider not configured. Let's set it up!")
			if err := configureProvider(aiProvider); err != nil {
				return fmt.Errorf("failed to configure provider: %w", err)
			}
			// Reintentar crear el cliente después de configurar
			factory = aifactory.NewClientFactory()
			client, err = factory.CreateClientFromString(aiProvider)
			if err != nil {
				return fmt.Errorf("error creating AI client after configuration: %w", err)
			}
		} else {
			return fmt.Errorf("error creating AI client: %w", err)
		}
	}

	// Verificar que el provider está disponible
	available, err := client.IsAvailable()
	if err != nil {
		return fmt.Errorf("error checking provider availability: %w", err)
	}
	if !available {
		return fmt.Errorf("selected provider is not available. Please run: claude-init config --provider %s", aiProvider)
	}

	log.Info("✓ AI provider configured: %s", aiProvider)

	// 5. NUEVA PREGUNTA: ¿Proyecto nuevo o existente?
	log.Info("\nStarting project configuration...")

	projectOrigin, err := askProjectOrigin()
	if err != nil {
		return fmt.Errorf("failed to ask project origin: %w", err)
	}

	var answers *survey.Answers

	// 6. Branch según el origen del proyecto
	if projectOrigin == "Existente" {
		answers, err = runExistingProjectFlow(projectPath, client)
	} else {
		answers, err = runNewProjectFlow(client)
	}

	if err != nil {
		return err
	}

	// Asignar el provider seleccionado
	answers.AIProvider = aiProvider

	// Validar respuestas
	if err := answers.Validate(); err != nil {
		return fmt.Errorf("invalid answers: %w", err)
	}

	// Guardar configuración del proyecto
	if err := saveProjectConfig(projectPath, opts.ConfigDir, answers); err != nil {
		return fmt.Errorf("failed to save project config: %w", err)
	}

	log.Info("\n✓ Project information collected successfully!")

	// 8. Generar estructura usando AI provider
	log.Info("\nGenerating .claude/ structure with AI provider...")

	if err := generateClaudeStructure(projectPath, opts, answers, client); err != nil {
		return fmt.Errorf("failed to generate structure: %w", err)
	}

	// 7. Mostrar resumen
	printSummary(cmd, answers, opts)

	return nil
}

// generateClaudeStructure genera la estructura .claude/ usando un Client de IA.
func generateClaudeStructure(projectPath string, opts *InitOptions, answers *survey.Answers, client ai.Client) error {
	outputDir := filepath.Join(projectPath, opts.ConfigDir)

	if opts.DryRun {
		// En modo dry-run, solo mostrar qué se generaría
		fmt.Printf("\nWould generate .claude/ structure at: %s\n", outputDir)
		fmt.Printf("  Project: %s\n", answers.ProjectName)
		fmt.Printf("  Language: %s\n", answers.Language)
		if answers.Framework != "" {
			fmt.Printf("  Framework: %s\n", answers.Framework)
		}
		fmt.Printf("  Architecture: %s\n", answers.Architecture)
		fmt.Printf("  Category: %s\n", answers.ProjectCategory)
		return nil
	}

	// Crear directorio base .claude
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create .claude directory: %w", err)
	}

	// Crear generador usando el client apropiado
	generator := claude.NewGenerator(projectPath, answers, client)
	generator.SetLogger(log)

	// Obtener recomendación usando AI provider
	log.Info("Getting structure recommendations from AI provider...")
	recommendation, err := generator.GetRecommendation()
	if err != nil {
		log.Warn("Failed to get recommendation from AI provider: %v", err)
		log.Info("Using default structure...")
		recommendation = getDefaultRecommendation(answers)
	}

	// Generar estructura completa basada en recomendación
	if err := generator.GenerateAll(recommendation); err != nil {
		return fmt.Errorf("failed to generate structure: %w", err)
	}

	log.Info("✓ Structure generated successfully")
	return nil
}

// getDefaultRecommendation retorna una recomendación por defecto basada en las respuestas.
func getDefaultRecommendation(answers *survey.Answers) *claude.Recommendation {
	agents := []string{"architect", "developer", "tester", "reviewer"}
	commands := []string{"test", "lint", "build"}
	skills := []string{normalizeSkillName(answers.Language)}

	// Agregar agent de debugger solo para arquitecturas complejas
	complexArchitectures := []string{"Microservicios", "DDD", "Hexagonal", "Event-Driven", "Serverless"}
	for _, arch := range complexArchitectures {
		if answers.Architecture == arch {
			agents = append(agents, "debugger")
			break
		}
	}

	// Agregar skill de framework si existe
	if answers.Framework != "" {
		skills = append(skills, strings.ToLower(answers.Framework))
	}

	description := fmt.Sprintf("Default structure for %s %s project", answers.Language, answers.ProjectCategory)

	return &claude.Recommendation{
		Agents:      agents,
		Commands:    commands,
		Skills:      skills,
		Description: description,
	}
}

// normalizeSkillName normaliza el nombre de un skill a kebab-case (ej: "Node.js" -> "nodejs", "Code Quality Agent" -> "code-quality-agent").
func normalizeSkillName(name string) string {
	// Normalizaciones específicas primero
	normalizations := map[string]string{
		"Node.js":    "nodejs",
		"JavaScript": "nodejs",
		"TypeScript": "typescript",
		"Go":         "go",
		"Python":     "python",
		"Rust":       "rust",
	}

	if normalized, ok := normalizations[name]; ok {
		return normalized
	}

	// Convertir a kebab-case: minúsculas y reemplazar espacios por guiones
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")

	// Eliminar múltiples guiones consecutivos
	reg := regexp.MustCompile(`-+`)
	name = reg.ReplaceAllString(name, "-")

	// Eliminar caracteres no alfanuméricos excepto guiones
	reg = regexp.MustCompile(`[^a-z0-9\-]`)
	name = reg.ReplaceAllString(name, "")

	// Eliminar guiones al inicio y al final
	name = strings.Trim(name, "-")

	if name == "" {
		return "unnamed"
	}

	return name
}

// getProjectPath retorna el path del proyecto desde los argumentos o el directorio actual.
func getProjectPath(args []string) (string, error) {
	if len(args) > 0 {
		// Convertir a path absoluto
		path := args[0]
		if !filepath.IsAbs(path) {
			abs, err := filepath.Abs(path)
			if err != nil {
				return "", fmt.Errorf("failed to get absolute path: %w", err)
			}
			return abs, nil
		}
		return path, nil
	}

	// Usar directorio actual
	return os.Getwd()
}

// validatePath valida que el path existe y es un directorio.
func validatePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}
		return fmt.Errorf("failed to stat path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}

	return nil
}

// checkExistingConfig verifica si ya existe configuración y retorna error si es necesario.
func checkExistingConfig(projectPath, configDir string, force bool) error {
	configPath := filepath.Join(projectPath, configDir)
	if _, err := os.Stat(configPath); err == nil {
		if !force {
			return fmt.Errorf("config directory already exists: %s (use --force to overwrite)", configPath)
		}
		// Si force es true, eliminar directorio existente
		if err := os.RemoveAll(configPath); err != nil {
			return fmt.Errorf("failed to remove existing config directory: %w", err)
		}
	}
	return nil
}

// printSummary imprime un resumen de lo generado.
func printSummary(cmd *cobra.Command, answers *survey.Answers, opts *InitOptions) {
	if log == nil {
		log = logger.New(os.Stdout, logger.INFOLevel)
	}

	if opts.DryRun {
		log.Info("\nDry run completed successfully!")
		return
	}

	log.Info("\n✓ .claude/ directory initialized successfully!")
	log.Info("  Project: %s", answers.ProjectName)
	log.Info("  Description: %s", answers.Description)
	log.Info("  Language: %s", answers.Language)
	if answers.Framework != "" {
		log.Info("  Framework: %s", answers.Framework)
	}
	log.Info("  Architecture: %s", answers.Architecture)
	if answers.Database != "" {
		log.Info("  Database: %s", answers.Database)
	}
	log.Info("  Category: %s", answers.ProjectCategory)

	log.Info("\nYou can now start using AI-assisted development with Claude Code!")
	log.Info("Try running: claude -p \"help me understand this codebase\"")
}

// ProjectConfig representa la configuración del proyecto guardada.
type ProjectConfig struct {
	ProjectOrigin   string `yaml:"project_origin" json:"project_origin"`
	ProjectName     string `yaml:"project_name" json:"project_name"`
	Description     string `yaml:"description" json:"description"`
	Language        string `yaml:"language" json:"language"`
	Framework       string `yaml:"framework" json:"framework"`
	Architecture    string `yaml:"architecture" json:"architecture"`
	Database        string `yaml:"database" json:"database"`
	ProjectCategory string `yaml:"project_category" json:"project_category"`
	BusinessContext string `yaml:"business_context" json:"business_context"`
	AIProvider      string `yaml:"ai_provider" json:"ai_provider"`
	CreatedAt       string `yaml:"created_at" json:"created_at"`
}

// saveProjectConfig guarda las respuestas del survey en un archivo YAML.
func saveProjectConfig(projectPath, configDir string, answers *survey.Answers) error {
	// Crear directorio de configuración si no existe
	configPath := filepath.Join(projectPath, configDir)
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Crear estructura de configuración
	config := ProjectConfig{
		ProjectOrigin:   answers.ProjectOrigin,
		ProjectName:     answers.ProjectName,
		Description:     answers.Description,
		Language:        answers.Language,
		Framework:       answers.Framework,
		Architecture:    answers.Architecture,
		Database:        answers.Database,
		ProjectCategory: answers.ProjectCategory,
		BusinessContext: answers.BusinessContext,
		AIProvider:      answers.AIProvider,
		CreatedAt:       time.Now().Format(time.RFC3339),
	}

	// Convertir a YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal project config: %w", err)
	}

	// Guardar en archivo
	projectFile := filepath.Join(configPath, "project.yaml")
	if err := os.WriteFile(projectFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write project config file: %w", err)
	}

	log.Debugf("Project config saved to: %s", projectFile)
	return nil
}

// askProjectOrigin pregunta si es un proyecto nuevo o existente.
func askProjectOrigin() (string, error) {
	var origin string
	prompt := &gSurvey.Select{
		Message: "¿Es este un proyecto NUEVO o existente?",
		Options: []string{"Nuevo", "Existente"},
		Default: "Nuevo",
	}

	if err := gSurvey.AskOne(prompt, &origin); err != nil {
		return "", err
	}

	return origin, nil
}

// askDocumentationDirs pregunta al usuario por directorios de documentación adicionales.
func askDocumentationDirs(projectPath string, answers *survey.Answers) (*survey.Answers, error) {
	// Primero, detectar automáticamente directorios comunes
	commonDocDirs := []string{"docs", "documentation", "guide", "guides", "wiki", "help"}
	foundDirs := make([]string, 0)

	for _, dir := range commonDocDirs {
		fullPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(fullPath); err == nil {
			foundDirs = append(foundDirs, dir)
		}
	}

	// Mostrar directorios detectados
	if len(foundDirs) > 0 {
		log.Info("\nDirectorios de documentación detectados:")
		for _, dir := range foundDirs {
			log.Info("  ✓ %s/", dir)
		}
		log.Info("")
	}

	// Preguntar si hay directorios adicionales
	var addMore bool
	prompt := &gSurvey.Confirm{
		Message: "¿Hay directorios adicionales de documentación que quieras incluir?",
		Default: false,
	}

	if err := gSurvey.AskOne(prompt, &addMore); err != nil {
		return answers, err
	}

	if !addMore {
		// Solo usar los detectados automáticamente
		answers.DocumentationDirs = foundDirs
		return answers, nil
	}

	// Preguntar por directorios adicionales (separados por coma)
	var dirsInput string
	inputPrompt := &gSurvey.Input{
		Message: "Directorios adicionales (separados por coma):",
		Help:    "Ejemplo: docs/architecture,docs/api,guides",
	}

	if err := gSurvey.AskOne(inputPrompt, &dirsInput); err != nil {
		// Si hay error o input vacío, usar solo los detectados
		answers.DocumentationDirs = foundDirs
		return answers, nil
	}

	// Procesar los directorios ingresados
	extraDirs := make([]string, 0)
	inputParts := strings.Split(dirsInput, ",")
	for _, part := range inputParts {
		dir := strings.TrimSpace(part)
		if dir == "" {
			continue
		}

		// Verificar que el directorio existe
		fullPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(fullPath); err != nil {
			log.Warn("  ⚠ Directorio '%s' no existe, se omitirá", dir)
			continue
		}

		extraDirs = append(extraDirs, dir)
		log.Info("  ✓ %s/ agregado", dir)
	}

	// Combinar directorios detectados y adicionales
	answers.DocumentationDirs = append(foundDirs, extraDirs...)

	if len(answers.DocumentationDirs) > 0 {
		log.Info("\nDirectorios de documentación a incluir:")
		for _, dir := range answers.DocumentationDirs {
			log.Info("  - %s/", dir)
		}
	}

	return answers, nil
}

// runExistingProjectFlow analiza y pre-llena el survey para proyectos existentes.
func runExistingProjectFlow(projectPath string, client ai.Client) (*survey.Answers, error) {
	log.Info("\nAnalizando proyecto existente...")
	log.Info("Esto puede tomar unos segundos...\n")

	analyzer := claude.NewAnalyzer(projectPath, client)
	analyzer.SetLogger(log)

	analysis, err := analyzer.Analyze()
	if err != nil {
		log.Warn("Project analysis failed: %v", err)
		log.Info("Falling back to manual survey...\n")
		return runNewProjectFlow(client)
	}

	// Convertir análisis a Answers
	prefill := &survey.Answers{
		ProjectOrigin:   "Existente",
		ProjectName:     analysis.Name,
		Description:     analysis.Description,
		Language:        analysis.Language,
		Framework:       analysis.Framework,
		Architecture:    analysis.Architecture,
		Database:        analysis.Database,
		ProjectCategory: analysis.ProjectCategory,
		BusinessContext: analysis.BusinessContext,
	}

	// Mostrar resultados del análisis
	log.Info("Análisis completado:")
	log.Info("  Nombre: %s", analysis.Name)
	log.Info("  Lenguaje: %s", analysis.Language)
	if analysis.Framework != "" {
		log.Info("  Framework: %s", analysis.Framework)
	}
	log.Info("  Arquitectura: %s", analysis.Architecture)
	if analysis.Database != "" {
		log.Info("  Base de datos: %s", analysis.Database)
	}
	log.Info("  Categoría: %s", analysis.ProjectCategory)
	log.Info("\nPor favor, revisa y edita la información pre-llenada:\n")

	// Ejecutar survey con valores pre-llenados
	questions := getProjectQuestions()
	runner := survey.NewRunner(questions)
	answers, err := runner.RunWithPrefill(prefill)
	if err != nil {
		return nil, err
	}

	// Preguntar por directorios de documentación adicionales
	answers, err = askDocumentationDirs(projectPath, answers)
	if err != nil {
		log.Warn("No se pudieron agregar directorios de documentación: %v", err)
		// Continuar sin los directorios de documentación
	}

	return answers, nil
}

// runNewProjectFlow ejecuta el survey normal para proyectos nuevos.
func runNewProjectFlow(client ai.Client) (*survey.Answers, error) {
	log.Info("Please answer the following questions to configure your project.\n")

	questions := getProjectQuestions()
	runner := survey.NewRunner(questions)
	answers, err := runner.Run()
	if err != nil {
		return nil, fmt.Errorf("survey failed: %w", err)
	}

	answers.ProjectOrigin = "Nuevo"
	return answers, nil
}

// getProjectQuestions retorna todas las preguntas del survey.
func getProjectQuestions() []*survey.Question {
	return survey.GetProjectQuestions()
}

// configureProvider configura interactivamente un provider de IA.
func configureProvider(provider string) error {
	// Cargar configuración existente
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.GlobalConfig{
			Provider:  provider,
			Providers: make(map[string]config.ProviderConfig),
		}
	}

	// Actualizar provider por defecto
	cfg.Provider = provider

	// Pedir API key
	var apiKey string
	prompt := &gSurvey.Password{
		Message: fmt.Sprintf("Enter your %s API key:", provider),
		Help:    fmt.Sprintf("Get your API key from: %s", getAPIKeyURL(provider)),
	}
	if err := gSurvey.AskOne(prompt, &apiKey, gSurvey.WithStdio(os.Stdin, os.Stderr, os.Stdout)); err != nil {
		return err
	}

	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	// Crear configuración del provider
	providerCfg := config.ProviderConfig{
		APIKey: apiKey,
	}

	// Preguntar si quiere configuración avanzada
	var advanced bool
	advancedPrompt := &gSurvey.Confirm{
		Message: "Do you want to configure advanced settings (base URL, model, max tokens)?",
		Default: false,
	}
	if err := gSurvey.AskOne(advancedPrompt, &advanced); err == nil && advanced {
		// Base URL
		var baseURL string
		baseURLPrompt := &gSurvey.Input{
			Message: "Base URL (press Enter for default):",
			Default: getDefaultBaseURL(provider),
		}
		if err := gSurvey.AskOne(baseURLPrompt, &baseURL); err == nil && baseURL != "" {
			providerCfg.BaseURL = baseURL
		}

		// Model
		var model string
		modelPrompt := &gSurvey.Input{
			Message: "Model (press Enter for default):",
			Default: getDefaultModel(provider),
		}
		if err := gSurvey.AskOne(modelPrompt, &model); err == nil && model != "" {
			providerCfg.Model = model
		}

		// Max Tokens
		var maxTokens int
		maxTokensPrompt := &gSurvey.Input{
			Message: "Max tokens (press Enter for default):",
			Default: "4096",
		}
		if err := gSurvey.AskOne(maxTokensPrompt, &maxTokens); err == nil && maxTokens > 0 {
			providerCfg.MaxTokens = maxTokens
		}
	}

	// Guardar configuración del provider
	cfg.Providers[provider] = providerCfg

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	log.Info("✓ Configuration saved for provider: %s", provider)
	return nil
}

// getAPIKeyURL retorna la URL para conseguir el API key de un provider.
func getAPIKeyURL(provider string) string {
	switch provider {
	case "claude-api":
		return "https://console.anthropic.com/settings/keys"
	case "openai":
		return "https://platform.openai.com/account/api-keys"
	case "zai":
		return "https://z.ai"
	case "gemini":
		return "https://aistudio.google.com/apikey"
	case "groq":
		return "https://console.groq.com/keys"
	default:
		return ""
	}
}

// getDefaultBaseURL retorna la base URL por defecto para un provider.
func getDefaultBaseURL(provider string) string {
	switch provider {
	case "claude-api":
		return "https://api.anthropic.com/v1/messages"
	case "openai":
		return "https://api.openai.com/v1"
	case "zai":
		return "https://api.z.ai/v1"
	case "gemini":
		return "https://generativelanguage.googleapis.com/v1beta/models"
	case "groq":
		return "https://api.groq.com/openai/v1"
	default:
		return ""
	}
}

// getDefaultModel retorna el modelo por defecto para un provider.
func getDefaultModel(provider string) string {
	switch provider {
	case "claude-api":
		return "claude-opus-4"
	case "openai":
		return "gpt-4o-mini"
	case "zai":
		return "glm-4.7"
	case "gemini":
		return "gemini-2.5-flash"
	case "groq":
		return "llama-3.3-70b-versatile"
	default:
		return ""
	}
}

// askAIProvider pregunta al usuario qué provider de IA quiere usar.
func askAIProvider() (string, error) {
	var provider string
	prompt := &gSurvey.Select{
		Message: "Selecciona el provider de IA a usar:",
		Options: []string{
			"Claude CLI (gratis con Claude Code PRO) - Opción por defecto",
			"OpenAI API (Requiere API key)",
			"Google Gemini API (Requiere API key) - ⚠️ Free tier NO válido para este CLI",
			"Groq API (Requiere API key) - ⚠️ Free tier NO válido para este CLI",
			// "Claude API (Requiere API key, más rápido)",
			// "Z.AI API (Requiere API key, más rápido)",
		},
		Default: "Claude CLI (gratis con Claude Code PRO) - Opción por defecto",
		Help:    "Claude CLI usa tu suscripción PRO. Gemini y Groq free tiers tienen límites muy bajos (~15-50 req/día).",
	}

	if err := gSurvey.AskOne(prompt, &provider); err != nil {
		return "", err
	}

	// Mapear la opción seleccionada al valor interno
	switch provider {
	case "Claude CLI (gratis con Claude Code PRO) - Opción por defecto":
		return "cli", nil
	case "OpenAI API (Requiere API key)":
		return "openai", nil
	case "Google Gemini API (Requiere API key) - ⚠️ Free tier NO válido para este CLI":
		return "gemini", nil
	case "Groq API (Requiere API key) - ⚠️ Free tier NO válido para este CLI":
		return "groq", nil
	// Comentado temporalmente - se usará más adelante
	// case "Claude API (Requiere API key, más rápido)":
	// 	return "claude-api", nil
	// case "Z.AI API (Requiere API key, más rápido)":
	// 	return "zai", nil
	default:
		return "cli", nil
	}
}
