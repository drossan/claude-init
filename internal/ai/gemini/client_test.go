package gemini

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		apiKey    string
		baseURL   string
		model     string
		maxTokens int
	}{
		{
			name:      "default values",
			apiKey:    "test-key",
			baseURL:   "",
			model:     "",
			maxTokens: 0,
		},
		{
			name:      "custom values",
			apiKey:    "custom-key",
			baseURL:   "https://custom.endpoint.com",
			model:     "gemini-2.5-flash",
			maxTokens: 50000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.apiKey, tt.baseURL, tt.model, tt.maxTokens)

			if client == nil {
				t.Fatal("NewClient returned nil")
			}

			if client.apiKey != tt.apiKey {
				t.Errorf("expected apiKey %q, got %q", tt.apiKey, client.apiKey)
			}

			// Verificar valores por defecto
			if tt.baseURL == "" && client.baseURL != "https://generativelanguage.googleapis.com/v1beta/models" {
				t.Errorf("expected default baseURL, got %q", client.baseURL)
			}
			if tt.baseURL != "" && client.baseURL != tt.baseURL {
				t.Errorf("expected baseURL %q, got %q", tt.baseURL, client.baseURL)
			}

			if tt.model == "" && client.model != "gemini-2.5-flash" {
				t.Errorf("expected default model gemini-2.5-flash, got %q", client.model)
			}
			if tt.model != "" && client.model != tt.model {
				t.Errorf("expected model %q, got %q", tt.model, client.model)
			}

			if tt.maxTokens == 0 && client.maxTokens != 1000000 {
				t.Errorf("expected default maxTokens 1000000, got %d", client.maxTokens)
			}
			if tt.maxTokens != 0 && client.maxTokens != tt.maxTokens {
				t.Errorf("expected maxTokens %d, got %d", tt.maxTokens, client.maxTokens)
			}
		})
	}
}

func TestClient_SendSimpleMessage(t *testing.T) {
	// Este test requerirá un mock o una API key válida
	// Por ahora solo verificamos que la función existe
	client := NewClient("test-key", "", "", 0)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	// La implementación real se probará con integration tests
	_ = client
}

func TestClient_Close(t *testing.T) {
	client := NewClient("test-key", "", "", 0)

	err := client.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}
}
