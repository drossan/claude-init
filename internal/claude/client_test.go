package claude

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCLIWrapper(t *testing.T) {
	wrapper := NewCLIWrapper()

	assert.NotNil(t, wrapper)
	assert.Equal(t, 2*time.Minute, wrapper.timeout)
}

func TestCLIWrapper_SetTimeout(t *testing.T) {
	wrapper := NewCLIWrapper()
	expectedTimeout := 5 * time.Second
	wrapper.SetTimeout(expectedTimeout)

	assert.Equal(t, expectedTimeout, wrapper.timeout)
}

// Nota: Los tests de CheckInstalled, GetVersion y SendMessage requieren
// que Claude CLI esté instalado, por lo que son tests de integración.
// Se pueden ejecutar manualmente con:
//   go test -v ./internal/claude -run TestCLIWrapper_Integration

func TestCLIWrapper_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	wrapper := NewCLIWrapper()

	// Test CheckInstalled
	err := wrapper.CheckInstalled()
	if err != nil {
		t.Skipf("Claude CLI not installed: %v", err)
	}

	// Test GetVersion
	version, err := wrapper.GetVersion()
	assert.NoError(t, err)
	assert.NotEmpty(t, version)
	t.Logf("Claude CLI version: %s", version)

	// Test SendSimpleMessage (solo si tenemos una cuenta válida)
	// Este test puede fallar si no hay autenticación configurada
	response, err := wrapper.SendSimpleMessage("Say 'Hello, World!'")
	if err != nil {
		t.Logf("SendSimpleMessage failed (expected if not authenticated): %v", err)
	} else {
		assert.NotEmpty(t, response)
		t.Logf("Response: %s", response)
	}
}
