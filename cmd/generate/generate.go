package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	aifactory "github.com/danielrossellosanchez/claude-init/internal/ai"
	"github.com/danielrossellosanchez/claude-init/internal/claude"
	"github.com/danielrossellosanchez/claude-init/internal/logger"
	"github.com/danielrossellosanchez/claude-init/internal/survey"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	log           *logger.Logger
	forceFlag     bool
	dryRunFlag    bool
	configDirFlag string
	outputDirFlag string
	agentsFlag    bool
	skillsFlag    bool
	commandsFlag  bool
	guidesFlag    bool
)

var generateCmd = &cobra.Command{
	Use:   "generate [project-path]",
	Short: "Generate Claude Code configuration using Claude CLI",
	Long: `Generate Claude Code configuration (.claude/) for a project.

This command loads project configuration from .claude/project.yaml
and uses Claude CLI to generate agents, skills, and commands.

Requires Claude CLI to be installed. Visit: https://claude.com/claude-code

Files generated:
  - agents/       : Agent configurations (architect, developer, tester, etc.)
  - skills/       : Language/framework specific skills
  - commands/     : Custom commands

Usage:
  claude-init generate                    # Use current directory
  claude-init generate ./my-project      # Use specific path
  claude-init generate --force            # Overwrite existing files
  claude-init generate --dry-run          # Show what would be generated
  claude-init generate --only-agents      # Generate only agents
`,
	Example: `  # Generate full configuration in current directory
  claude-init generate

  # Generate for a specific project
  claude-init generate ./path/to/project

  # Generate only agents and skills
  claude-init generate --only-agents --only-skills

  # Show what would be generated without creating files
  claude-init generate --dry-run

  # Overwrite existing files
  claude-init generate --force`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGenerate,
}

// Execute añade el comando generate al root command.
func Execute(rootCmd *cobra.Command, l *logger.Logger) {
	log = l
	rootCmd.AddCommand(generateCmd)
}

// GetGenerateCmd retorna el comando generate para testing.
func GetGenerateCmd() *cobra.Command {
	return generateCmd
}

func init() {
	generateCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "overwrite existing files")
	generateCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "show what would be generated without creating files")
	generateCmd.Flags().StringVar(&configDirFlag, "config-dir", "", "config directory (default: .claude/)")
	generateCmd.Flags().StringVar(&outputDirFlag, "output-dir", "", "output directory (default: <project>/.claude/)")
	generateCmd.Flags().BoolVar(&agentsFlag, "only-agents", false, "generate only agents")
	generateCmd.Flags().BoolVar(&skillsFlag, "only-skills", false, "generate only skills")
	generateCmd.Flags().BoolVar(&commandsFlag, "only-commands", false, "generate only commands")
	generateCmd.Flags().BoolVar(&guidesFlag, "only-guides", false, "generate only guides")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Inicializar logger si es nil
	if log == nil {
		log = logger.New(os.Stdout, logger.INFOLevel)
	}

	// Verificar que Claude CLI está instalado
	log.Info("Checking Claude CLI installation...")
	if err := claude.CheckInstalled(); err != nil {
		return err
	}
	log.Info("✓ Claude CLI detected")

	// Determinar el path del proyecto
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}

	// Convertir a path absoluto
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Verificar que el directorio existe
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("project path does not exist: %s", absPath)
	}

	log.Info("Generating Claude Code configuration...")
	log.Debugf("Project path: %s", absPath)

	// Cargar configuración existente del proyecto
	projectConfig, err := loadProjectConfig(absPath)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w (run 'claude-init init' first)", err)
	}

	// Convertir ProjectConfig a survey.Answers para usar con el generador
	answers := &survey.Answers{
		ProjectName:     projectConfig.ProjectName,
		Description:     projectConfig.Description,
		Language:        projectConfig.Language,
		Framework:       projectConfig.Framework,
		Architecture:    projectConfig.Architecture,
		Database:        projectConfig.Database,
		ProjectCategory: projectConfig.ProjectCategory,
		BusinessContext: projectConfig.BusinessContext,
		AIProvider:      projectConfig.AIProvider,
	}

	// Determinar el directorio de salida
	outputDir := filepath.Join(absPath, ".claude")
	if outputDirFlag != "" {
		outputDir = outputDirFlag
	}

	// Verificar si ya existe la configuración
	if _, err := os.Stat(outputDir); err == nil && !forceFlag {
		return fmt.Errorf("configuration directory already exists: %s (use --force to overwrite)", outputDir)
	}

	// Crear cliente según el provider configurado
	factory := aifactory.NewClientFactory()
	client, err := factory.CreateClientFromString(projectConfig.AIProvider)
	if err != nil {
		return fmt.Errorf("error creating AI client: %w", err)
	}

	// Crear generador usando el cliente
	generator := claude.NewGenerator(absPath, answers, client)
	generator.SetLogger(log)

	// Obtener recomendación usando AI provider
	log.Info("Getting structure recommendations from AI provider...")
	recommendation, err := generator.GetRecommendation()
	if err != nil {
		log.Warn("Failed to get recommendation from AI provider: %v", err)
		log.Info("Using default structure...")
		recommendation = getDefaultRecommendation(answers)
	}

	// Determinar qué generar
	generateAgents := agentsFlag
	generateSkills := skillsFlag
	generateCommands := commandsFlag
	generateGuides := guidesFlag

	// Si no se especificó nada, generar todo
	if !generateAgents && !generateSkills && !generateCommands && !generateGuides {
		generateAgents = true
		generateSkills = true
		generateCommands = true
		generateGuides = true
	}

	// Generar según lo solicitado
	if dryRunFlag {
		return runDryRun(recommendation, outputDir, generateAgents, generateSkills, generateCommands, generateGuides)
	}

	// Crear directorio base
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generar agentes
	if generateAgents {
		for _, agent := range recommendation.Agents {
			if err := generator.GenerateAgent(agent); err != nil {
				log.Warn("Failed to generate agent %s: %v", agent, err)
			}
		}
		log.Info("✓ Agents generated")
	}

	// Generar skills
	if generateSkills {
		for _, skill := range recommendation.Skills {
			// Determinar tipo de skill basado en el contexto
			skillType := "language"
			if answers.Framework != "" && strings.EqualFold(skill, answers.Framework) {
				skillType = "framework"
			}
			if err := generator.GenerateSkill(skillType, skill); err != nil {
				log.Warn("Failed to generate skill %s: %v", skill, err)
			}
		}
		log.Info("✓ Skills generated")
	}

	// Generar comandos
	if generateCommands {
		for _, command := range recommendation.Commands {
			if err := generator.GenerateCommand(command); err != nil {
				log.Warn("Failed to generate command %s: %v", command, err)
			}
		}
		log.Info("✓ Commands generated")
	}

	// Generar guías (no implementado todavía, solo placeholder)
	if generateGuides {
		log.Info("Guides generation not yet implemented")
	}

	log.Info("✓ Configuration generated successfully at: %s", outputDir)
	return nil
}

func runDryRun(rec *claude.Recommendation, outputDir string, agents, skills, commands, guides bool) error {
	log.Info("Dry run mode - showing what would be generated:")
	log.Info("Output directory: %s", outputDir)
	log.Info("")

	if agents {
		log.Info("Would generate agents:")
		for _, agent := range rec.Agents {
			log.Info("  - %s.md", agent)
		}
	}

	if skills {
		log.Info("Would generate skills:")
		for _, skill := range rec.Skills {
			log.Info("  - %s.md", skill)
		}
	}

	if commands {
		log.Info("Would generate commands:")
		for _, cmd := range rec.Commands {
			log.Info("  - %s.md", cmd)
		}
	}

	if guides {
		log.Info("Would generate guides:")
		log.Info("  - development_guide.md")
	}

	return nil
}

// ProjectConfig representa la configuración del proyecto guardada por init.
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

// loadProjectConfig carga la configuración del proyecto desde .claude/
func loadProjectConfig(projectPath string) (*ProjectConfig, error) {
	configPath := filepath.Join(projectPath, ".claude", "project.yaml")

	// Verificar si existe el archivo
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("project config not found at %s", configPath)
	}

	// Leer el archivo YAML
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read project config: %w", err)
	}

	// Parsear YAML
	var config ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse project config: %w", err)
	}

	return &config, nil
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

// normalizeSkillName normaliza el nombre de un skill (ej: "Node.js" -> "nodejs").
func normalizeSkillName(name string) string {
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
	return strings.ToLower(name)
}
