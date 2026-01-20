package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/danielrossellosanchez/claude-init/internal/claude"
	"github.com/danielrossellosanchez/claude-init/internal/logger"
	"github.com/danielrossellosanchez/claude-init/internal/survey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestNewGenerateCommand(t *testing.T) {
	cmd := GetGenerateCmd()

	if cmd == nil {
		t.Fatal("Generate command should not be nil")
	}

	if cmd.Use != "generate [project-path]" {
		t.Errorf("Expected use 'generate [project-path]', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Short description should not be empty")
	}
}

func TestGenerateFlags(t *testing.T) {
	cmd := GetGenerateCmd()

	// Flags que deben existir
	requiredFlags := []string{
		"force",
		"dry-run",
		"config-dir",
		"output-dir",
		"only-agents",
		"only-skills",
		"only-commands",
		"only-guides",
	}

	for _, flag := range requiredFlags {
		if f := cmd.Flags().Lookup(flag); f == nil {
			t.Errorf("Flag '%s' should be defined", flag)
		}
	}

	// Verificar shorthand para force
	forceFlag := cmd.Flags().Lookup("force")
	assert.Equal(t, "f", forceFlag.Shorthand)
}

func TestGenerateCommand_DoesNotHaveAIFlags(t *testing.T) {
	cmd := GetGenerateCmd()

	// Flags que NO deben existir (fueron eliminados)
	removedFlags := []string{
		"ai-provider",
		"api-key",
	}

	for _, flag := range removedFlags {
		flagRef := cmd.Flags().Lookup(flag)
		assert.Nil(t, flagRef, "flag %s should be removed", flag)
	}
}

func TestLoadProjectConfig_WithValidConfig_ReturnsConfig(t *testing.T) {
	tempDir := t.TempDir()

	// Crear directorio .claude
	claudeDir := filepath.Join(tempDir, ".claude")
	err := os.MkdirAll(claudeDir, 0755)
	require.NoError(t, err)

	// Crear project.yaml
	config := &ProjectConfig{
		ProjectOrigin:   "new",
		ProjectName:     "test-project",
		Description:     "Test description",
		Language:        "Go",
		Framework:       "Gin",
		Architecture:    "Clean",
		Database:        "PostgreSQL",
		ProjectCategory: "API REST",
		BusinessContext: "Test business context",
	}

	data, err := yaml.Marshal(config)
	require.NoError(t, err)

	configPath := filepath.Join(claudeDir, "project.yaml")
	err = os.WriteFile(configPath, data, 0644)
	require.NoError(t, err)

	// Cargar configuración
	loaded, err := loadProjectConfig(tempDir)
	assert.NoError(t, err)
	assert.Equal(t, config.ProjectName, loaded.ProjectName)
	assert.Equal(t, config.Description, loaded.Description)
	assert.Equal(t, config.Language, loaded.Language)
}

func TestLoadProjectConfig_WithoutConfig_ReturnsError(t *testing.T) {
	tempDir := t.TempDir()

	// No crear project.yaml

	// Intentar cargar configuración
	_, err := loadProjectConfig(tempDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

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

func TestRunDryRun_DoesNotCreateFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Crear logger
	log = logger.New(os.Stdout, logger.INFOLevel)

	// Crear recomendación de prueba
	rec := &claude.Recommendation{
		Agents:      []string{"architect", "developer"},
		Commands:    []string{"test", "build"},
		Skills:      []string{"go"},
		Description: "Test recommendation",
	}

	// Ejecutar dry-run
	err := runDryRun(rec, filepath.Join(tempDir, ".claude"), true, true, true, true)
	assert.NoError(t, err)

	// Verificar que NO se creó el directorio
	claudeDir := filepath.Join(tempDir, ".claude")
	_, err = os.Stat(claudeDir)
	assert.True(t, os.IsNotExist(err), "directory should not exist with --dry-run")
}

func TestGenerateCommand_Integration(t *testing.T) {
	t.Skip("Skipping - requires Claude CLI to be installed")

	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create .claude directory with project.yaml
	claudeDir := filepath.Join(tempDir, ".claude")
	err := os.MkdirAll(claudeDir, 0755)
	require.NoError(t, err)

	config := &ProjectConfig{
		ProjectOrigin:   "new",
		ProjectName:     "test-project",
		Description:     "Test description",
		Language:        "Go",
		Framework:       "",
		Architecture:    "Clean",
		Database:        "",
		ProjectCategory: "API",
		BusinessContext: "Test context",
	}

	data, err := yaml.Marshal(config)
	require.NoError(t, err)

	configPath := filepath.Join(claudeDir, "project.yaml")
	err = os.WriteFile(configPath, data, 0644)
	require.NoError(t, err)

	// Create logger
	testLog := logger.NewDefault()

	// Setup command
	cmd := GetGenerateCmd()
	cmd.SetArgs([]string{"--dry-run", tempDir})

	// Set logger for the command
	log = testLog

	// Execute
	err = cmd.Execute()
	assert.NoError(t, err)
}
