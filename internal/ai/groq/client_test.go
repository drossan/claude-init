package groq

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
			model:     "llama-3.3-70b-versatile",
			maxTokens: 8192,
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
			if tt.baseURL == "" && client.baseURL != "https://api.groq.com/openai/v1" {
				t.Errorf("expected default baseURL, got %q", client.baseURL)
			}
			if tt.baseURL != "" && client.baseURL != tt.baseURL {
				t.Errorf("expected baseURL %q, got %q", tt.baseURL, client.baseURL)
			}

			if tt.model == "" && client.model != "llama-3.3-70b-versatile" {
				t.Errorf("expected default model llama-3.3-70b-versatile, got %q", client.model)
			}
			if tt.model != "" && client.model != tt.model {
				t.Errorf("expected model %q, got %q", tt.model, client.model)
			}

			if tt.maxTokens == 0 && client.maxTokens != 32768 {
				t.Errorf("expected default maxTokens 32768, got %d", client.maxTokens)
			}
			if tt.maxTokens != 0 && client.maxTokens != tt.maxTokens {
				t.Errorf("expected maxTokens %d, got %d", tt.maxTokens, client.maxTokens)
			}
		})
	}
}

func TestClient_SendSimpleMessage(t *testing.T) {
	client := NewClient("test-key", "", "", 0)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	_ = client
}

func TestClient_Close(t *testing.T) {
	client := NewClient("test-key", "", "", 0)

	err := client.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}
}
