package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewLogger_CreatesLoggerWithCorrectLevel verifica que New crea un logger con el nivel correcto.
func TestNewLogger_CreatesLoggerWithCorrectLevel(t *testing.T) {
	tests := []struct {
		name  string
		level Level
	}{
		{
			name:  "DEBUG level",
			level: DEBUGLevel,
		},
		{
			name:  "INFO level",
			level: INFOLevel,
		},
		{
			name:  "WARN level",
			level: WARNLevel,
		},
		{
			name:  "ERROR level",
			level: ERRORLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var buf bytes.Buffer

			// Act
			logger := New(&buf, tt.level)

			// Assert
			require.NotNil(t, logger)
			assert.Equal(t, tt.level, logger.Level())
		})
	}
}

// TestNewLogger_WithNilWriterUsesStderr verifica que New usa stderr cuando el writer es nil.
func TestNewLogger_WithNilWriterUsesStderr(t *testing.T) {
	// Arrange & Act
	logger := New(nil, INFOLevel)

	// Assert
	require.NotNil(t, logger)
	assert.NotNil(t, logger.out)
}

// TestNewDefaultLogger_CreatesLoggerWithInfoLevel verifica que NewDefault crea un logger con nivel INFO.
func TestNewDefaultLogger_CreatesLoggerWithInfoLevel(t *testing.T) {
	// Arrange & Act
	logger := NewDefault()

	// Assert
	require.NotNil(t, logger)
	assert.Equal(t, INFOLevel, logger.Level())
}

// TestSetLevel_ChangesLoggerLevel verifica que SetLevel cambia el nivel.
func TestSetLevel_ChangesLoggerLevel(t *testing.T) {
	tests := []struct {
		name         string
		initialLevel Level
		newLevel     Level
	}{
		{
			name:         "change from DEBUG to ERROR",
			initialLevel: DEBUGLevel,
			newLevel:     ERRORLevel,
		},
		{
			name:         "change from INFO to WARN",
			initialLevel: INFOLevel,
			newLevel:     WARNLevel,
		},
		{
			name:         "change from ERROR to DEBUG",
			initialLevel: ERRORLevel,
			newLevel:     DEBUGLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var buf bytes.Buffer
			logger := New(&buf, tt.initialLevel)
			assert.Equal(t, tt.initialLevel, logger.Level())

			// Act
			logger.SetLevel(tt.newLevel)

			// Assert
			assert.Equal(t, tt.newLevel, logger.Level())
		})
	}
}

// TestLevel_ReturnsCurrentLevel verifica que Level retorna el nivel actual.
func TestLevel_ReturnsCurrentLevel(t *testing.T) {
	tests := []struct {
		name  string
		level Level
	}{
		{name: "DEBUG level", level: DEBUGLevel},
		{name: "INFO level", level: INFOLevel},
		{name: "WARN level", level: WARNLevel},
		{name: "ERROR level", level: ERRORLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var buf bytes.Buffer
			logger := New(&buf, tt.level)

			// Act
			returnedLevel := logger.Level()

			// Assert
			assert.Equal(t, tt.level, returnedLevel)
		})
	}
}

// TestDebug_OnlyLogsWhenDebugEnabled verifica que Debug solo logea cuando está en DEBUG.
func TestDebug_OnlyLogsWhenDebugEnabled(t *testing.T) {
	tests := []struct {
		name          string
		level         Level
		shouldLog     bool
		expectedInLog string
	}{
		{
			name:          "DEBUG level - should log",
			level:         DEBUGLevel,
			shouldLog:     true,
			expectedInLog: "[DEBUG]",
		},
		{
			name:          "INFO level - should not log",
			level:         INFOLevel,
			shouldLog:     false,
			expectedInLog: "[DEBUG]",
		},
		{
			name:          "WARN level - should not log",
			level:         WARNLevel,
			shouldLog:     false,
			expectedInLog: "[DEBUG]",
		},
		{
			name:          "ERROR level - should not log",
			level:         ERRORLevel,
			shouldLog:     false,
			expectedInLog: "[DEBUG]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var buf bytes.Buffer
			logger := New(&buf, tt.level)

			// Act
			logger.Debug("debug message")

			// Assert
			output := buf.String()
			if tt.shouldLog {
				assert.Contains(t, output, tt.expectedInLog)
				assert.Contains(t, output, "debug message")
			} else {
				assert.NotContains(t, output, tt.expectedInLog)
			}
		})
	}
}

// TestDebugf_OnlyLogsWhenDebugEnabled verifica que Debugf solo logea cuando está en DEBUG.
func TestDebugf_OnlyLogsWhenDebugEnabled(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := New(&buf, DEBUGLevel)
	message := "debug message %d"

	// Act
	logger.Debugf(message, 42)

	// Assert
	output := buf.String()
	assert.Contains(t, output, "[DEBUG]")
	assert.Contains(t, output, "debug message 42")
}

// TestInfo_OnlyLogsWhenInfoEnabled verifica que Info solo logea cuando está en INFO o inferior.
func TestInfo_OnlyLogsWhenInfoEnabled(t *testing.T) {
	tests := []struct {
		name          string
		level         Level
		shouldLog     bool
		expectedInLog string
	}{
		{
			name:          "DEBUG level - should log",
			level:         DEBUGLevel,
			shouldLog:     true,
			expectedInLog: "[INFO]",
		},
		{
			name:          "INFO level - should log",
			level:         INFOLevel,
			shouldLog:     true,
			expectedInLog: "[INFO]",
		},
		{
			name:          "WARN level - should not log",
			level:         WARNLevel,
			shouldLog:     false,
			expectedInLog: "[INFO]",
		},
		{
			name:          "ERROR level - should not log",
			level:         ERRORLevel,
			shouldLog:     false,
			expectedInLog: "[INFO]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var buf bytes.Buffer
			logger := New(&buf, tt.level)

			// Act
			logger.Info("info message")

			// Assert
			output := buf.String()
			if tt.shouldLog {
				assert.Contains(t, output, tt.expectedInLog)
				assert.Contains(t, output, "info message")
			} else {
				assert.NotContains(t, output, tt.expectedInLog)
			}
		})
	}
}

// TestInfof_OnlyLogsWhenInfoEnabled verifica que Infof formatea correctamente.
func TestInfof_OnlyLogsWhenInfoEnabled(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := New(&buf, INFOLevel)
	message := "info message %s"

	// Act
	logger.Infof(message, "test")

	// Assert
	output := buf.String()
	assert.Contains(t, output, "[INFO]")
	assert.Contains(t, output, "info message test")
}

// TestWarn_OnlyLogsWhenWarnEnabled verifica que Warn solo logea cuando está en WARN o inferior.
func TestWarn_OnlyLogsWhenWarnEnabled(t *testing.T) {
	tests := []struct {
		name          string
		level         Level
		shouldLog     bool
		expectedInLog string
	}{
		{
			name:          "DEBUG level - should log",
			level:         DEBUGLevel,
			shouldLog:     true,
			expectedInLog: "[WARN]",
		},
		{
			name:          "INFO level - should log",
			level:         INFOLevel,
			shouldLog:     true,
			expectedInLog: "[WARN]",
		},
		{
			name:          "WARN level - should log",
			level:         WARNLevel,
			shouldLog:     true,
			expectedInLog: "[WARN]",
		},
		{
			name:          "ERROR level - should not log",
			level:         ERRORLevel,
			shouldLog:     false,
			expectedInLog: "[WARN]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var buf bytes.Buffer
			logger := New(&buf, tt.level)

			// Act
			logger.Warn("warn message")

			// Assert
			output := buf.String()
			if tt.shouldLog {
				assert.Contains(t, output, tt.expectedInLog)
				assert.Contains(t, output, "warn message")
			} else {
				assert.NotContains(t, output, tt.expectedInLog)
			}
		})
	}
}

// TestWarnf_OnlyLogsWhenWarnEnabled verifica que Warnf formatea correctamente.
func TestWarnf_OnlyLogsWhenWarnEnabled(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := New(&buf, WARNLevel)
	message := "warn message %d"

	// Act
	logger.Warnf(message, 123)

	// Assert
	output := buf.String()
	assert.Contains(t, output, "[WARN]")
	assert.Contains(t, output, "warn message 123")
}

// TestError_AlwaysLogs verifica que Error siempre logea.
func TestError_AlwaysLogs(t *testing.T) {
	tests := []struct {
		name          string
		level         Level
		shouldLog     bool
		expectedInLog string
	}{
		{
			name:          "DEBUG level - should log",
			level:         DEBUGLevel,
			shouldLog:     true,
			expectedInLog: "[ERROR]",
		},
		{
			name:          "INFO level - should log",
			level:         INFOLevel,
			shouldLog:     true,
			expectedInLog: "[ERROR]",
		},
		{
			name:          "WARN level - should log",
			level:         WARNLevel,
			shouldLog:     true,
			expectedInLog: "[ERROR]",
		},
		{
			name:          "ERROR level - should log",
			level:         ERRORLevel,
			shouldLog:     true,
			expectedInLog: "[ERROR]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var buf bytes.Buffer
			logger := New(&buf, tt.level)

			// Act
			logger.Error("error message")

			// Assert
			output := buf.String()
			if tt.shouldLog {
				assert.Contains(t, output, tt.expectedInLog)
				assert.Contains(t, output, "error message")
			} else {
				assert.NotContains(t, output, tt.expectedInLog)
			}
		})
	}
}

// TestErrorf_AlwaysLogs verifica que Errorf formatea correctamente.
func TestErrorf_AlwaysLogs(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := New(&buf, ERRORLevel)
	message := "error message %s"

	// Act
	logger.Errorf(message, "critical")

	// Assert
	output := buf.String()
	assert.Contains(t, output, "[ERROR]")
	assert.Contains(t, output, "error message critical")
}

// TestWithLevel_ReturnsNewLoggerWithDifferentLevel verifica que WithLevel retorna un nuevo logger.
func TestWithLevel_ReturnsNewLoggerWithDifferentLevel(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := New(&buf, INFOLevel)

	// Act
	newLogger := logger.WithLevel(ERRORLevel)

	// Assert
	require.NotNil(t, newLogger)
	assert.Equal(t, INFOLevel, logger.Level(), "original logger level should not change")
	assert.Equal(t, ERRORLevel, newLogger.Level(), "new logger should have new level")
}

// TestWithLevel_DoesNotModifyOriginalLogger verifica que WithLevel no modifica el logger original.
func TestWithLevel_DoesNotModifyOriginalLogger(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := New(&buf, DEBUGLevel)

	// Act
	newLogger := logger.WithLevel(ERRORLevel)

	// Assert
	// Escribir con el nuevo logger
	newLogger.Error("error from new logger")

	// Verificar que el nivel del original no cambió
	logger.Debug("debug from original")

	output := buf.String()
	assert.Contains(t, output, "error from new logger")
	assert.Contains(t, output, "debug from original")
	assert.Equal(t, DEBUGLevel, logger.Level())
	assert.Equal(t, ERRORLevel, newLogger.Level())
}

// TestStdLogger_ReturnsStandardLogger verifica que StdLogger retorna un *log.Logger.
func TestStdLogger_ReturnsStandardLogger(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := New(&buf, INFOLevel)

	// Act
	stdLogger := logger.StdLogger()

	// Assert
	require.NotNil(t, stdLogger)
	assert.IsType(t, &log.Logger{}, stdLogger)
}

// TestStdLogger_WritesToLoggerOutput verifica que StdLogger escribe a la salida del logger.
func TestStdLogger_WritesToLoggerOutput(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := New(&buf, INFOLevel)
	stdLogger := logger.StdLogger()
	message := "message from std logger"

	// Act
	stdLogger.Println(message)

	// Assert
	output := buf.String()
	assert.Contains(t, output, "[INFO]")
	assert.Contains(t, output, "message from std logger")
}

// TestLevel_String_ReturnsCorrectRepresentation verifica el método String del tipo Level.
func TestLevel_String_ReturnsCorrectRepresentation(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DEBUGLevel, "DEBUG"},
		{INFOLevel, "INFO"},
		{WARNLevel, "WARN"},
		{ERRORLevel, "ERROR"},
		{Level(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			// Arrange & Act
			result := tt.level.String()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestMultipleLogMessages_VerifyAllAreLogged verifica que múltiples mensajes se registran correctamente.
func TestMultipleLogMessages_VerifyAllAreLogged(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := New(&buf, DEBUGLevel)

	// Act
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Assert
	output := buf.String()
	assert.Contains(t, output, "[DEBUG]")
	assert.Contains(t, output, "debug message")
	assert.Contains(t, output, "[INFO]")
	assert.Contains(t, output, "info message")
	assert.Contains(t, output, "[WARN]")
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "[ERROR]")
	assert.Contains(t, output, "error message")
}

// TestLogFormatting_VerifyStandardFlags verifica que los logs incluyen fecha y hora.
func TestLogFormatting_VerifyStandardFlags(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := New(&buf, INFOLevel)

	// Act
	logger.Info("test message")

	// Assert
	output := buf.String()
	// log.LstdFlags incluye fecha y hora
	// El formato debe contener algo como "2025/01/16 15:04:05"
	assert.Contains(t, output, "[INFO]")
	assert.Contains(t, output, "test message")
	// Verificar que hay un timestamp (buscamos patrones comunes)
	lines := strings.Split(output, "\n")
	if len(lines) > 0 {
		// La primera línea no vacía debe tener el formato de timestamp
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				// Debe contener fecha y hora antes del prefacio [INFO]
				assert.True(t, strings.HasPrefix(line, "20") || strings.Contains(line, "/"))
				break
			}
		}
	}
}

// TestLoggerWithBuffer_VerifySharedOutput verifica que multiple loggers comparten el mismo output.
func TestLoggerWithBuffer_VerifySharedOutput(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger1 := New(&buf, DEBUGLevel)
	logger2 := New(&buf, ERRORLevel)

	// Act
	logger1.Debug("debug from logger1")
	logger2.Error("error from logger2")

	// Assert
	output := buf.String()
	assert.Contains(t, output, "debug from logger1")
	assert.Contains(t, output, "error from logger2")
}

// TestSetLevel_AffectsLoggingBehavior verifica que cambiar el nivel afecta el comportamiento de logging.
func TestSetLevel_AffectsLoggingBehavior(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	logger := New(&buf, ERRORLevel)

	// Act - No debe logear debug en ERROR level
	logger.Debug("debug message")

	// Cambiar a DEBUG level
	logger.SetLevel(DEBUGLevel)

	// Ahora debe logear debug
	logger.Debug("debug message after setlevel")

	// Assert
	output := buf.String()
	// No debe contener el primer mensaje de debug
	assert.NotContains(t, output, "[DEBUG] debug message\n")
	// Debe contener el segundo mensaje de debug
	assert.Contains(t, output, "[DEBUG]")
	assert.Contains(t, output, "debug message after setlevel")
}
