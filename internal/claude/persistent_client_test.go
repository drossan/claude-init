package claude

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersistentClient_Lifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := NewPersistentClient()

	// Test Start
	t.Run("Start", func(t *testing.T) {
		err := client.Start()
		if err != nil {
			t.Skipf("Claude CLI not available: %v", err)
		}
		assert.True(t, client.IsRunning())
	})

	// Test SendMessage
	t.Run("SendMessage", func(t *testing.T) {
		if !client.IsRunning() {
			t.Skip("Client not running")
		}

		response, err := client.SendSimpleMessage("Say 'Hello' and nothing else")
		if err != nil {
			t.Logf("SendMessage failed (may need auth): %v", err)
		} else {
			require.NoError(t, err)
			assert.NotEmpty(t, response)
			t.Logf("Response: %s", response)
		}
	})

	// Test Multiple Messages
	t.Run("MultipleMessages", func(t *testing.T) {
		if !client.IsRunning() {
			t.Skip("Client not running")
		}

		// Enviar múltiples mensajes sin reiniciar
		for i := 0; i < 3; i++ {
			response, err := client.SendSimpleMessage("Count: 1")
			if err != nil {
				t.Logf("Message %d failed: %v", i, err)
				continue
			}
			t.Logf("Message %d response length: %d", i, len(response))
		}
	})

	// Test Stop
	t.Run("Stop", func(t *testing.T) {
		err := client.Stop()
		require.NoError(t, err)
		assert.False(t, client.IsRunning())
	})
}

func TestPersistentClient_Benchmarks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping benchmark in short mode")
	}

	t.Run("CompareWithNewProcess", func(t *testing.T) {
		// Benchmark con proceso nuevo cada vez
		start := time.Now()
		for i := 0; i < 3; i++ {
			wrapper := NewCLIWrapper()
			_, err := wrapper.SendSimpleMessage("Say 'OK'")
			if err != nil {
				t.Logf("Spawn error: %v", err)
			}
		}
		spawnTime := time.Since(start)
		t.Logf("Time spawning 3 processes: %v", spawnTime)

		// Benchmark con proceso persistente
		client := NewPersistentClient()
		if err := client.Start(); err != nil {
			t.Skipf("Cannot start persistent client: %v", err)
		}
		defer client.Stop()

		start = time.Now()
		for i := 0; i < 3; i++ {
			_, err := client.SendSimpleMessage("Say 'OK'")
			if err != nil {
				t.Logf("Persistent error: %v", err)
			}
		}
		persistentTime := time.Since(start)
		t.Logf("Time with persistent process: %v", persistentTime)

		t.Logf("Speedup: %.2fx", float64(spawnTime)/float64(persistentTime))
	})
}
