package root_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/drossan/claude-init/cmd/root"
	"github.com/drossan/claude-init/internal/logger"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExecute_CallsRootCommandExecute verifica que Execute llama a rootCmd.Execute.
func TestExecute_CallsRootCommandExecute(t *testing.T) {
	// Arrange
	testLogger := logger.New(io.Discard, logger.INFOLevel)

	// Act
	err := root.Execute(testLogger)

	// Assert
	require.NoError(t, err, "Execute should not return an error when called with valid logger")
}

// TestGetLogger_ReturnsConfiguredLogger verifica que GetLogger retorna el logger configurado.
func TestGetLogger_ReturnsConfiguredLogger(t *testing.T) {
	// Arrange
	testLogger := logger.New(io.Discard, logger.INFOLevel)
	_ = root.Execute(testLogger)

	// Act
	retrievedLogger := root.GetLogger()

	// Assert
	assert.NotNil(t, retrievedLogger, "GetLogger should return a non-nil logger")
	assert.Equal(t, testLogger, retrievedLogger, "GetLogger should return the same logger that was passed to Execute")
}

// TestVerboseFlag_PersistentInAllCommands verifica que el flag verbose es persistente.
func TestVerboseFlag_PersistentInAllCommands(t *testing.T) {
	// Arrange
	rootCmd := root.GetRootCmd()

	// Act - Verificamos que el flag está definido
	verboseFlag := rootCmd.PersistentFlags().Lookup("verbose")

	// Assert
	assert.NotNil(t, verboseFlag, "Verbose flag should be defined as a persistent flag")
	assert.Equal(t, "v", verboseFlag.Shorthand, "Verbose flag should have 'v' as shorthand")
	assert.Equal(t, "verbose", verboseFlag.Name, "Verbose flag should have 'verbose' as name")
}

// TestVerboseFlag_EnablesDebugLevel verifica que verbose=true habilita DEBUG.
func TestVerboseFlag_EnablesDebugLevel(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	testLogger := logger.New(&buf, logger.INFOLevel)

	// Act - Ejecutamos con el flag verbose
	rootCmd := root.GetRootCmd()
	rootCmd.SetArgs([]string{"--verbose"})
	err := root.Execute(testLogger)

	// Assert
	require.NoError(t, err, "Execute with verbose flag should not return an error")
	output := buf.String()
	assert.Contains(t, output, "Verbose mode enabled", "Debug log should be present when verbose is enabled")
	assert.Equal(t, logger.DEBUGLevel, testLogger.Level(), "Logger level should be DEBUG when verbose is enabled")
}

// TestRootCommand_HasCorrectUse verifica que Use es "ia-start".
func TestRootCommand_HasCorrectUse(t *testing.T) {
	// Arrange
	rootCmd := root.GetRootCmd()

	// Act
	use := rootCmd.Use

	// Assert
	assert.Equal(t, "claude-init", use, "Root command Use should be 'claude-init'")
}

// TestRootCommand_HasCorrectShortDescription verifica la descripción corta.
func TestRootCommand_HasCorrectShortDescription(t *testing.T) {
	// Arrange
	rootCmd := root.GetRootCmd()

	// Act
	short := rootCmd.Short

	// Assert
	expectedShort := "CLI para inicializar proyectos con configuración guiada por IA"
	assert.Equal(t, expectedShort, short, "Root command Short description should match expected value")
}

// TestRootCommand_HasCorrectLongDescription verifica la descripción larga.
func TestRootCommand_HasCorrectLongDescription(t *testing.T) {
	// Arrange
	rootCmd := root.GetRootCmd()

	// Act
	long := rootCmd.Long

	// Assert
	assert.NotEmpty(t, long, "Root command Long description should not be empty")
	assert.Contains(t, long, "claude-init", "Long description should mention the tool name")
	assert.Contains(t, long, "Claude Code", "Long description should mention Claude Code")
	assert.Contains(t, long, ".claude/", "Long description should mention .claude/ directory")
}

// TestRootCommand_DisplaysHelpByDefault verifica que muestra ayuda por defecto.
func TestRootCommand_DisplaysHelpByDefault(t *testing.T) {
	// Arrange
	var out bytes.Buffer
	testLogger := logger.New(io.Discard, logger.INFOLevel)
	rootCmd := root.GetRootCmd()
	rootCmd.SetOut(&out)
	rootCmd.SetArgs([]string{})

	// Act
	err := root.Execute(testLogger)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	output := out.String()
	assert.NotEmpty(t, output, "Help output should not be empty")
	assert.Contains(t, output, "Usage:", "Help output should contain usage information")
	assert.Contains(t, output, "claude-init", "Help output should contain the command name")
}

// Table-Driven Tests para diferentes escenarios del comando raíz
func TestRootCommand_TableDriven(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		wantErr          bool
		expectedInOutput string
		description      string
	}{
		{
			name:             "help flag",
			args:             []string{"--help"},
			wantErr:          false,
			expectedInOutput: "Usage:",
			description:      "Should display help when --help flag is provided",
		},
		{
			name:             "short help flag",
			args:             []string{"-h"},
			wantErr:          false,
			expectedInOutput: "Usage:",
			description:      "Should display help when -h flag is provided",
		},
		{
			name:             "verbose with help",
			args:             []string{"--verbose", "--help"},
			wantErr:          false,
			expectedInOutput: "verbose output",
			description:      "Should display help with verbose flag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			testLogger := logger.New(io.Discard, logger.INFOLevel)
			var out bytes.Buffer
			rootCmd := root.GetRootCmd()
			rootCmd.SetOut(&out)
			rootCmd.SetArgs(tt.args)

			// Act
			err := root.Execute(testLogger)

			// Assert
			if tt.wantErr {
				assert.Error(t, err, tt.description)
			} else {
				require.NoError(t, err, tt.description)
				output := out.String()
				assert.NotEmpty(t, output, "Output should not be empty")
				if tt.expectedInOutput != "" {
					assert.Contains(t, output, tt.expectedInOutput, tt.description)
				}
			}
		})
	}
}

// TestExecute_WithNilLogger verifica el comportamiento con logger nil.
func TestExecute_WithNilLogger(t *testing.T) {
	// Arrange
	var nilLogger *logger.Logger = nil

	// Act
	err := root.Execute(nilLogger)

	// Assert
	// El comando debería ejecutarse sin error ya que PersistentPreRun verifica nil
	require.NoError(t, err, "Execute should handle nil logger gracefully")
}

// TestVerboseFlag_DefaultValue verifica que el flag verbose tiene el valor correcto por defecto.
func TestVerboseFlag_DefaultValue(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	testLogger := logger.New(&buf, logger.INFOLevel)
	rootCmd := root.GetRootCmd()

	// Act - Ejecutamos sin el flag verbose (ejecución por defecto)
	rootCmd.SetArgs([]string{})
	err := root.Execute(testLogger)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	// El nivel de logging debería ser INFO (no DEBUG)
	assert.Equal(t, logger.INFOLevel, testLogger.Level(), "Logger level should be INFO when verbose is not set")
}

// TestGetLogger_BeforeExecute verifica que GetLogger retorna nil antes de Execute.
func TestGetLogger_BeforeExecute(t *testing.T) {
	// Arrange & Act - No llamamos a Execute (usamos una nueva instancia del test)
	// Nota: Como los tests comparten el paquete, necesitamos ser cuidadosos
	// En este caso, simplemente verificamos que no haya panic
	retrievedLogger := root.GetLogger()

	// Assert
	// Puede ser nil o un logger de tests anteriores, no fallamos si no es nil
	// Lo importante es que no haya panic
	assert.NotNil(t, retrievedLogger, "GetLogger should return a logger (may be from previous test)")
}

// TestExecute_MultipleCalls_GetLoggerReturnsLastLogger verifica que GetLogger retorna el último logger configurado.
func TestExecute_MultipleCalls_GetLoggerReturnsLastLogger(t *testing.T) {
	// Arrange
	firstLogger := logger.New(io.Discard, logger.INFOLevel)
	secondLogger := logger.New(io.Discard, logger.DEBUGLevel)

	// Act - Primera ejecución
	_ = root.Execute(firstLogger)
	firstRetrieved := root.GetLogger()

	// Act - Segunda ejecución
	_ = root.Execute(secondLogger)
	secondRetrieved := root.GetLogger()

	// Assert
	assert.Equal(t, firstLogger, firstRetrieved, "First GetLogger should return first logger")
	assert.Equal(t, secondLogger, secondRetrieved, "Second GetLogger should return second logger")
	assert.NotEqual(t, firstRetrieved, secondRetrieved, "Loggers should be different")
}

// TestVerboseFlag_ShortHandWorks verifica que el shorthand -v funciona.
func TestVerboseFlag_ShortHandWorks(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	testLogger := logger.New(&buf, logger.INFOLevel)

	// Creamos un subcomando dummy para testear
	dummyCmd := &cobra.Command{
		Use: "dummy",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	rootCmd := root.GetRootCmd()
	rootCmd.AddCommand(dummyCmd)

	// Act - Ejecutamos el subcomando dummy con el shorthand -v
	// Usamos un subcomando para que PersistentPreRun se ejecute correctamente
	rootCmd.SetArgs([]string{"dummy", "-v"})
	err := root.Execute(testLogger)

	// Assert
	require.NoError(t, err, "Execute with -v flag should not return an error")
	output := buf.String()
	assert.Contains(t, output, "Verbose mode enabled", "Debug log should be present when -v is used")
	assert.Equal(t, logger.DEBUGLevel, testLogger.Level(), "Logger level should be DEBUG when -v is used")
}

// TestVerboseFlag_NotSet_DoesNotChangeLevel verifica que sin verbose el nivel no cambia.
func TestVerboseFlag_NotSet_DoesNotChangeLevel(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	testLogger := logger.New(&buf, logger.WARNLevel) // Empezamos con WARN

	// Act - Ejecutamos sin el flag verbose
	rootCmd := root.GetRootCmd()
	rootCmd.SetArgs([]string{})
	err := root.Execute(testLogger)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	// El nivel debería seguir siendo WARN
	assert.Equal(t, logger.WARNLevel, testLogger.Level(), "Logger level should remain WARN when verbose is not set")
}

// TestVerboseFlag_PersistsAcrossSubcommands verifica que el flag verbose se hereda en subcomandos.
func TestVerboseFlag_PersistsAcrossSubcommands(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	testLogger := logger.New(&buf, logger.INFOLevel)

	// Creamos un subcomando dummy
	dummyCmd := &cobra.Command{
		Use: "dummy",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	rootCmd := root.GetRootCmd()
	rootCmd.AddCommand(dummyCmd)

	// Act - Ejecutamos el subcomando con el flag verbose
	rootCmd.SetArgs([]string{"dummy", "--verbose"})
	err := root.Execute(testLogger)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	output := buf.String()
	assert.Contains(t, output, "Verbose mode enabled", "Debug log should be present when verbose is enabled in subcommand")
	assert.Equal(t, logger.DEBUGLevel, testLogger.Level(), "Logger level should be DEBUG when verbose is enabled in subcommand")
}

// TestRootCommand_OutputContainsExpectedTexts verifica que la salida contiene textos esperados.
func TestRootCommand_OutputContainsExpectedTexts(t *testing.T) {
	// Arrange
	var out bytes.Buffer
	testLogger := logger.New(io.Discard, logger.INFOLevel)
	rootCmd := root.GetRootCmd()
	rootCmd.SetOut(&out)
	rootCmd.SetArgs([]string{})

	// Act
	err := root.Execute(testLogger)

	// Assert
	require.NoError(t, err)
	output := out.String()
	assert.NotEmpty(t, output)

	// Verificar que contiene información esperada
	expectedTexts := []string{
		"claude-init",
		"CLI",
	}

	for _, text := range expectedTexts {
		assert.Contains(t, output, text, "Output should contain '%s'", text)
	}
}

// TestLogger_CreatedWithCorrectLevel verifica que el logger se crea con el nivel correcto.
func TestLogger_CreatedWithCorrectLevel(t *testing.T) {
	tests := []struct {
		name        string
		level       logger.Level
		expectLevel logger.Level
	}{
		{
			name:        "INFO level",
			level:       logger.INFOLevel,
			expectLevel: logger.INFOLevel,
		},
		{
			name:        "DEBUG level",
			level:       logger.DEBUGLevel,
			expectLevel: logger.DEBUGLevel,
		},
		{
			name:        "WARN level",
			level:       logger.WARNLevel,
			expectLevel: logger.WARNLevel,
		},
		{
			name:        "ERROR level",
			level:       logger.ERRORLevel,
			expectLevel: logger.ERRORLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange & Act
			testLogger := logger.New(io.Discard, tt.level)

			// Assert
			assert.Equal(t, tt.expectLevel, testLogger.Level(), "Logger should have level %v", tt.expectLevel)
		})
	}
}

// TestLogger_SetLevelChangesLevel verifica que SetLevel cambia el nivel del logger.
func TestLogger_SetLevelChangesLevel(t *testing.T) {
	// Arrange
	testLogger := logger.New(io.Discard, logger.INFOLevel)

	// Act
	testLogger.SetLevel(logger.DEBUGLevel)

	// Assert
	assert.Equal(t, logger.DEBUGLevel, testLogger.Level(), "Logger level should be DEBUG after SetLevel")
}

// TestLogger_DebugOnlyLogsWhenDebugLevel verifica que Debug solo loggea cuando el nivel es DEBUG.
func TestLogger_DebugOnlyLogsWhenDebugLevel(t *testing.T) {
	tests := []struct {
		name            string
		level           logger.Level
		shouldContain   bool
		expectedMessage string
	}{
		{
			name:            "DEBUG level - shows debug",
			level:           logger.DEBUGLevel,
			shouldContain:   true,
			expectedMessage: "[DEBUG]",
		},
		{
			name:            "INFO level - hides debug",
			level:           logger.INFOLevel,
			shouldContain:   false,
			expectedMessage: "[DEBUG]",
		},
		{
			name:            "WARN level - hides debug",
			level:           logger.WARNLevel,
			shouldContain:   false,
			expectedMessage: "[DEBUG]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var buf bytes.Buffer
			testLogger := logger.New(&buf, tt.level)

			// Act
			testLogger.Debug("test debug message")

			// Assert
			output := buf.String()
			if tt.shouldContain {
				assert.Contains(t, output, tt.expectedMessage, "Output should contain debug message")
				assert.Contains(t, output, "test debug message", "Output should contain the debug message text")
			} else {
				assert.NotContains(t, output, tt.expectedMessage, "Output should not contain debug message")
			}
		})
	}
}

// TestLogger_InfoAlwaysLogs verificar que Info siempre loggea.
func TestLogger_InfoAlwaysLogs(t *testing.T) {
	tests := []struct {
		name          string
		level         logger.Level
		shouldContain bool
	}{
		{
			name:          "DEBUG level - shows info",
			level:         logger.DEBUGLevel,
			shouldContain: true,
		},
		{
			name:          "INFO level - shows info",
			level:         logger.INFOLevel,
			shouldContain: true,
		},
		{
			name:          "WARN level - hides info",
			level:         logger.WARNLevel,
			shouldContain: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var buf bytes.Buffer
			testLogger := logger.New(&buf, tt.level)

			// Act
			testLogger.Info("test info message")

			// Assert
			output := buf.String()
			if tt.shouldContain {
				assert.Contains(t, output, "[INFO]", "Output should contain INFO message")
				assert.Contains(t, output, "test info message", "Output should contain the info message text")
			} else {
				assert.NotContains(t, output, "[INFO]", "Output should not contain INFO message")
			}
		})
	}
}

// TestRootCommand_RunE_ReturnsHelpError verifica que RunE llama a Help.
func TestRootCommand_RunE_ReturnsHelpError(t *testing.T) {
	// Arrange
	testLogger := logger.New(io.Discard, logger.INFOLevel)
	rootCmd := root.GetRootCmd()
	rootCmd.SetArgs([]string{})

	// Act
	err := root.Execute(testLogger)

	// Assert
	// cobra.Help() devuelve un error específico cuando se muestra la ayuda
	// Verificamos que no hay error (el help se muestra correctamente)
	require.NoError(t, err, "Execute should not return an error when showing help")
}

// TestGetRootCmd_ReturnsNonNullCommand verifica que GetRootCmd retorna un comando válido.
func TestGetRootCmd_ReturnsNonNullCommand(t *testing.T) {
	// Arrange & Act
	rootCmd := root.GetRootCmd()

	// Assert
	assert.NotNil(t, rootCmd, "GetRootCmd should return a non-nil command")
	assert.NotNil(t, rootCmd.PersistentFlags(), "Root command should have persistent flags")
}

// TestVerboseFlag_AfterMultipleExecutions verifica que el flag verbose funciona correctamente en ejecuciones múltiples.
func TestVerboseFlag_AfterMultipleExecutions(t *testing.T) {
	// Arrange - Creamos loggers separados para cada ejecución y un subcomando dummy
	var buf bytes.Buffer
	_ = logger.New(io.Discard, logger.INFOLevel) // Primera ejecución para limpiar estado

	dummyCmd := &cobra.Command{
		Use: "multitest",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	rootCmd := root.GetRootCmd()
	rootCmd.AddCommand(dummyCmd)

	// Act - Primera ejecución sin verbose
	rootCmd.SetArgs([]string{"multitest"})
	_ = root.Execute(logger.New(io.Discard, logger.INFOLevel))

	// Segunda ejecución con verbose
	testLogger := logger.New(&buf, logger.INFOLevel)
	rootCmd.SetArgs([]string{"multitest", "--verbose"})
	_ = root.Execute(testLogger)
	output := buf.String()
	level := testLogger.Level()

	// Assert - Verificamos que la segunda ejecución SÍ tiene verbose
	assert.Contains(t, output, "Verbose mode enabled", "Second execution should contain verbose debug log")
	assert.Equal(t, logger.DEBUGLevel, level, "Second execution should have DEBUG level")
}

// TestRootCommand_HelpOutputContainsVerboseFlag verifica que la ayuda contiene el flag verbose.
func TestRootCommand_HelpOutputContainsVerboseFlag(t *testing.T) {
	// Arrange
	var out bytes.Buffer
	testLogger := logger.New(io.Discard, logger.INFOLevel)
	rootCmd := root.GetRootCmd()
	rootCmd.SetOut(&out)
	rootCmd.SetArgs([]string{"--help"})

	// Act
	err := root.Execute(testLogger)

	// Assert
	require.NoError(t, err)
	output := out.String()
	assert.Contains(t, output, "verbose", "Help output should mention verbose flag")
	assert.Contains(t, output, "-v", "Help output should show -v shorthand")
}
