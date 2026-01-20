package init

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/danielrossellosanchez/claude-init/internal/logger"
	"github.com/danielrossellosanchez/claude-init/internal/survey"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewInitCommand_CreatesCommand verifica que NewInitCommand crea un comando válido.
func TestNewInitCommand_CreatesCommand(t *testing.T) {
	cmd := NewInitCommand()

	assert.NotNil(t, cmd)
	assert.Contains(t, cmd.Use, "init")
	assert.Contains(t, cmd.Short, "Initialize")
	assert.NotNil(t, cmd.RunE)
}

// TestNewInitCommand_HasCorrectFlags verifica que el comando tiene los flags correctos (sin IA).
func TestNewInitCommand_HasCorrectFlags(t *testing.T) {
	cmd := NewInitCommand()

	// Flags que deben existir
	requiredFlags := []string{
		"force",
		"dry-run",
		"config-dir",
	}

	for _, flag := range requiredFlags {
		assert.NotNil(t, cmd.Flags().Lookup(flag), "missing flag: %s", flag)
	}

	// Verificar shorthand para force
	forceFlag := cmd.Flags().Lookup("force")
	assert.Equal(t, "f", forceFlag.Shorthand)
}

// TestNewInitCommand_DoesNotHaveAIFlags verifica que los flags de IA fueron eliminados.
func TestNewInitCommand_DoesNotHaveAIFlags(t *testing.T) {
	cmd := NewInitCommand()

	// Flags que NO deben existir (fueron eliminados)
	removedFlags := []string{
		"ai-provider",
		"api-key",
		"no-ai",
	}

	for _, flag := range removedFlags {
		flagRef := cmd.Flags().Lookup(flag)
		assert.Nil(t, flagRef, "flag %s should be removed", flag)
	}
}

// TestInitCommand_WithNonExistentPath_ReturnsError verifica que retorna error con path no existente.
func TestInitCommand_WithNonExistentPath_ReturnsError(t *testing.T) {
	t.Skip("Skipping - requires Claude CLI to be installed")

	// Crear un path que no existe
	nonExistentPath := filepath.Join(os.TempDir(), "non-existent-path-xyz123")

	cmd := NewInitCommand()
	cmd.SetArgs([]string{nonExistentPath})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

// TestInitCommand_WithDryRun_DoesNotCreateFiles verifica que --dry-run no crea archivos.
// SKIP: Este test requiere interacción del usuario, se omite por ahora.
func TestInitCommand_WithDryRun_DoesNotCreateFiles(t *testing.T) {
	t.Skip("Skipping interactive test - requires user input")

	tempDir := t.TempDir()

	// Inicializar logger para el test
	log = logger.New(os.Stdout, logger.INFOLevel)

	cmd := NewInitCommand()
	cmd.SetArgs([]string{"--dry-run", tempDir})

	err := cmd.Execute()
	assert.NoError(t, err)

	// Verificar que NO se creó el directorio .claude/
	claudeDir := filepath.Join(tempDir, ".claude")
	_, err = os.Stat(claudeDir)
	assert.True(t, os.IsNotExist(err), "directory should not exist with --dry-run")
}

// TestValidatePath_ValidPath_ReturnsNil verifica que validatePath retorna nil para path válido.
func TestValidatePath_ValidPath_ReturnsNil(t *testing.T) {
	tempDir := t.TempDir()

	err := validatePath(tempDir)
	assert.NoError(t, err)
}

// TestValidatePath_NonExistentPath_ReturnsError verifica que validatePath retorna error para path no existente.
func TestValidatePath_NonExistentPath_ReturnsError(t *testing.T) {
	nonExistentPath := "/non/existent/path/xyz123"

	err := validatePath(nonExistentPath)
	assert.Error(t, err)
}

// TestCheckExistingConfig_WithExistingDir_ReturnsError verifica que checkExistingConfig retorna error si existe.
func TestCheckExistingConfig_WithExistingDir_ReturnsError(t *testing.T) {
	tempDir := t.TempDir()

	// Crear directorio .claude/ existente
	claudeDir := filepath.Join(tempDir, ".claude")
	err := os.MkdirAll(claudeDir, 0755)
	require.NoError(t, err)

	err = checkExistingConfig(tempDir, ".claude", false)
	assert.Error(t, err)
}

// TestCheckExistingConfig_WithForceFlag_ReturnsNil verifica que checkExistingConfig retorna nil con --force.
func TestCheckExistingConfig_WithForceFlag_ReturnsNil(t *testing.T) {
	tempDir := t.TempDir()

	// Crear directorio .claude/ existente
	claudeDir := filepath.Join(tempDir, ".claude")
	err := os.MkdirAll(claudeDir, 0755)
	require.NoError(t, err)

	err = checkExistingConfig(tempDir, ".claude", true)
	assert.NoError(t, err)
}

// TestCheckExistingConfig_WithoutExistingDir_ReturnsNil verifica que checkExistingConfig retorna nil si no existe.
func TestCheckExistingConfig_WithoutExistingDir_ReturnsNil(t *testing.T) {
	tempDir := t.TempDir()

	err := checkExistingConfig(tempDir, ".claude", false)
	assert.NoError(t, err)
}

// TestGetDefaultRecommendation_GeneratesValidRecommendation verifica que getDefaultRecommendation genera recomendaciones válidas.
func TestGetDefaultRecommendation_GeneratesValidRecommendation(t *testing.T) {
	tests := []struct {
		name    string
		answers *survey.Answers
		wantLen struct {
			agents   int
			commands int
		}
		mustContainAgents []string
		mustContainSkills []string
	}{
		{
			name: "Go project with clean architecture",
			answers: &survey.Answers{
				Language:        "Go",
				Architecture:    "Clean",
				ProjectCategory: "API REST",
			},
			wantLen: struct {
				agents   int
				commands int
			}{agents: 4, commands: 3},
			mustContainAgents: []string{"architect", "developer", "tester", "reviewer"},
			mustContainSkills: []string{"go"},
		},
		{
			name: "Node.js project with framework and microservices",
			answers: &survey.Answers{
				Language:        "Node.js",
				Framework:       "Express",
				Architecture:    "Microservicios",
				ProjectCategory: "API REST",
			},
			wantLen: struct {
				agents   int
				commands int
			}{agents: 5, commands: 3}, // 5 agents because of microservices architecture
			mustContainAgents: []string{"debugger"},
			mustContainSkills: []string{"nodejs", "express"},
		},
		{
			name: "Python project with monolith architecture",
			answers: &survey.Answers{
				Language:        "Python",
				Architecture:    "Monolito",
				ProjectCategory: "Web App",
			},
			wantLen: struct {
				agents   int
				commands int
			}{agents: 4, commands: 3}, // 4 agents (no debugger for monolith)
			mustContainAgents: []string{"architect", "developer", "tester", "reviewer"},
			mustContainSkills: []string{"python"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := getDefaultRecommendation(tt.answers)

			assert.Len(t, rec.Agents, tt.wantLen.agents)
			assert.Len(t, rec.Commands, tt.wantLen.commands)
			assert.NotEmpty(t, rec.Description)

			for _, agent := range tt.mustContainAgents {
				assert.Contains(t, rec.Agents, agent)
			}

			for _, skill := range tt.mustContainSkills {
				assert.Contains(t, rec.Skills, skill)
			}
		})
	}
}

// TestGenerateClaudeStructure_CreatesAllComponents verifica que generateClaudeStructure crea todos los componentes.
func TestGenerateClaudeStructure_CreatesAllComponents(t *testing.T) {
	t.Skip("Skipping - requires Claude CLI to be installed")

	tempDir := t.TempDir()

	answers := &survey.Answers{
		ProjectName:     "test-project",
		Description:     "Test project",
		Language:        "Go",
		Architecture:    "Clean",
		ProjectCategory: "API REST",
	}

	opts := &InitOptions{
		ConfigDir: ".claude",
		DryRun:    false,
	}

	err := generateClaudeStructure(tempDir, opts, answers)
	assert.NoError(t, err)

	// Verificar que se creó la estructura
	claudeDir := filepath.Join(tempDir, ".claude")
	_, err = os.Stat(claudeDir)
	assert.NoError(t, err)

	// Verificar que se crearon los subdirectorios
	agentsDir := filepath.Join(claudeDir, "agents")
	_, err = os.Stat(agentsDir)
	assert.NoError(t, err)

	commandsDir := filepath.Join(claudeDir, "commands")
	_, err = os.Stat(commandsDir)
	assert.NoError(t, err)

	skillsDir := filepath.Join(claudeDir, "skills")
	_, err = os.Stat(skillsDir)
	assert.NoError(t, err)
}

// TestPrintSummary_DoesNotPanic verifica que printSummary no panic.
func TestPrintSummary_DoesNotPanic(t *testing.T) {
	cmd := &cobra.Command{}
	answers := &survey.Answers{
		ProjectName:     "test",
		Description:     "test project",
		Language:        "Go",
		Architecture:    "Clean",
		ProjectCategory: "API REST",
	}
	opts := &InitOptions{DryRun: true}

	// Capturar stdout para evitar imprimir en tests
	log = logger.New(os.Stdout, logger.INFOLevel)

	assert.NotPanics(t, func() {
		printSummary(cmd, answers, opts)
	})
}

// TestSaveProjectConfig_SavesValidYAML verifica que saveProjectConfig guarda YAML válido.
func TestSaveProjectConfig_SavesValidYAML(t *testing.T) {
	tempDir := t.TempDir()

	answers := &survey.Answers{
		ProjectName:     "test-project",
		Description:     "Test description",
		Language:        "Go",
		Framework:       "Gin",
		Architecture:    "Clean",
		Database:        "PostgreSQL",
		ProjectCategory: "API REST",
		BusinessContext: "Test business context",
	}

	err := saveProjectConfig(tempDir, ".claude", answers)
	assert.NoError(t, err)

	// Verificar que se creó el archivo
	configPath := filepath.Join(tempDir, ".claude", "project.yaml")
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// Leer y verificar contenido
	content, err := os.ReadFile(configPath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "test-project")
	assert.Contains(t, string(content), "Test description")
}
