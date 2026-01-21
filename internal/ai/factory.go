package ai

import (
	"fmt"

	"github.com/drossan/claude-init/internal/ai/claudeapi"
	"github.com/drossan/claude-init/internal/ai/cli"
	"github.com/drossan/claude-init/internal/ai/gemini"
	"github.com/drossan/claude-init/internal/ai/groq"
	"github.com/drossan/claude-init/internal/ai/openai"
	"github.com/drossan/claude-init/internal/ai/zai"
	"github.com/drossan/claude-init/internal/config"
)

// ClientFactory crea clientes de IA según el provider.
type ClientFactory struct {
	config *config.GlobalConfig
}

// NewClientFactory crea una nueva fábrica de clientes.
func NewClientFactory() *ClientFactory {
	cfg, err := config.Load()
	if err != nil {
		// Si hay error, retornar configuración vacía
		cfg = &config.GlobalConfig{
			Providers: make(map[string]config.ProviderConfig),
		}
	}
	return &ClientFactory{
		config: cfg,
	}
}

// CreateClient crea un cliente según el provider especificado.
func (f *ClientFactory) CreateClient(provider Provider) (Client, error) {
	switch provider {
	case ProviderCLI:
		return NewCLIClient(), nil

	case ProviderClaudeAPI:
		cfg, ok := f.config.GetProviderConfig("claude-api")
		if !ok || cfg.APIKey == "" {
			return nil, fmt.Errorf("claude-api provider not configured. Please run: claude-init config --provider claude-api")
		}
		client := claudeapi.NewClient(cfg.APIKey, cfg.BaseURL, cfg.Model, cfg.MaxTokens)
		return &ClaudeAPIClient{client: client}, nil

	case ProviderOpenAI:
		cfg, ok := f.config.GetProviderConfig("openai")
		if !ok || cfg.APIKey == "" {
			return nil, fmt.Errorf("openai provider not configured. Please run: claude-init config --provider openai")
		}
		client := openai.NewClient(cfg.APIKey, cfg.BaseURL, cfg.Model, cfg.MaxTokens)
		return &OpenAIClient{client: client}, nil

	case ProviderZAI:
		cfg, ok := f.config.GetProviderConfig("zai")
		if !ok || cfg.APIKey == "" {
			return nil, fmt.Errorf("zai provider not configured. Please run: claude-init config --provider zai")
		}
		client := zai.NewClient(cfg.APIKey, cfg.BaseURL, cfg.Model, cfg.MaxTokens)
		return &ZAIClient{client: client}, nil

	case ProviderGemini:
		cfg, ok := f.config.GetProviderConfig("gemini")
		if !ok || cfg.APIKey == "" {
			return nil, fmt.Errorf("gemini provider not configured. Please run: claude-init config --provider gemini")
		}
		client := gemini.NewClient(cfg.APIKey, cfg.BaseURL, cfg.Model, cfg.MaxTokens)
		return &GeminiClient{client: client}, nil

	case ProviderGroq:
		cfg, ok := f.config.GetProviderConfig("groq")
		if !ok || cfg.APIKey == "" {
			return nil, fmt.Errorf("groq provider not configured. Please run: claude-init config --provider groq")
		}
		client := groq.NewClient(cfg.APIKey, cfg.BaseURL, cfg.Model, cfg.MaxTokens)
		return &GroqClient{client: client}, nil

	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}

// CreateClientFromString crea un cliente desde el string del provider.
func (f *ClientFactory) CreateClientFromString(providerStr string) (Client, error) {
	provider := Provider(providerStr)
	if !provider.IsValid() {
		return nil, fmt.Errorf("invalid provider: %s", providerStr)
	}
	return f.CreateClient(provider)
}

// CLIClient es un wrapper que usa cli.HybridClient internamente.
type CLIClient struct {
	wrapper *cli.HybridClient
}

// NewCLIClient crea un nuevo cliente CLI.
func NewCLIClient() *CLIClient {
	return &CLIClient{
		wrapper: cli.NewHybridClient(),
	}
}

// SendMessage envía un mensaje usando Claude CLI.
func (c *CLIClient) SendMessage(systemPrompt, userMessage string) (string, error) {
	return c.wrapper.SendMessage(systemPrompt, userMessage)
}

// SendSimpleMessage envía un mensaje sin system prompt.
func (c *CLIClient) SendSimpleMessage(message string) (string, error) {
	return c.wrapper.SendSimpleMessage(message)
}

// Provider retorna el tipo de provider.
func (c *CLIClient) Provider() Provider {
	return ProviderCLI
}

// IsAvailable verifica si Claude CLI está instalado.
func (c *CLIClient) IsAvailable() (bool, error) {
	err := c.wrapper.CheckInstalled()
	return err == nil, err
}

// Close cierra el cliente.
func (c *CLIClient) Close() error {
	return c.wrapper.Stop()
}

// ClaudeAPIClient es un wrapper para el cliente de Claude API.
type ClaudeAPIClient struct {
	client *claudeapi.Client
}

// SendMessage envía un mensaje usando Claude API.
func (c *ClaudeAPIClient) SendMessage(systemPrompt, userMessage string) (string, error) {
	return c.client.SendMessage(systemPrompt, userMessage)
}

// SendSimpleMessage envía un mensaje sin system prompt.
func (c *ClaudeAPIClient) SendSimpleMessage(message string) (string, error) {
	return c.client.SendSimpleMessage(message)
}

// Provider retorna el tipo de provider.
func (c *ClaudeAPIClient) Provider() Provider {
	return ProviderClaudeAPI
}

// IsAvailable siempre retorna true (si hay API key configurada).
func (c *ClaudeAPIClient) IsAvailable() (bool, error) {
	return true, nil
}

// Close cierra el cliente.
func (c *ClaudeAPIClient) Close() error {
	return c.client.Close()
}

// OpenAIClient es un wrapper para el cliente de OpenAI.
type OpenAIClient struct {
	client *openai.Client
}

// SendMessage envía un mensaje usando OpenAI API.
func (c *OpenAIClient) SendMessage(systemPrompt, userMessage string) (string, error) {
	return c.client.SendMessage(systemPrompt, userMessage)
}

// SendSimpleMessage envía un mensaje sin system prompt.
func (c *OpenAIClient) SendSimpleMessage(message string) (string, error) {
	return c.client.SendSimpleMessage(message)
}

// Provider retorna el tipo de provider.
func (c *OpenAIClient) Provider() Provider {
	return ProviderOpenAI
}

// IsAvailable siempre retorna true (si hay API key configurada).
func (c *OpenAIClient) IsAvailable() (bool, error) {
	return true, nil
}

// Close cierra el cliente.
func (c *OpenAIClient) Close() error {
	return c.client.Close()
}

// ZAIClient es un wrapper para el cliente de ZAI.
type ZAIClient struct {
	client *zai.Client
}

// SendMessage envía un mensaje usando ZAI API.
func (c *ZAIClient) SendMessage(systemPrompt, userMessage string) (string, error) {
	return c.client.SendMessage(systemPrompt, userMessage)
}

// SendSimpleMessage envía un mensaje sin system prompt.
func (c *ZAIClient) SendSimpleMessage(message string) (string, error) {
	return c.client.SendSimpleMessage(message)
}

// Provider retorna el tipo de provider.
func (c *ZAIClient) Provider() Provider {
	return ProviderZAI
}

// IsAvailable siempre retorna true (si hay API key configurada).
func (c *ZAIClient) IsAvailable() (bool, error) {
	return true, nil
}

// Close cierra el cliente.
func (c *ZAIClient) Close() error {
	return c.client.Close()
}

// GeminiClient es un wrapper para el cliente de Gemini.
type GeminiClient struct {
	client *gemini.Client
}

// SendMessage envía un mensaje usando Gemini API.
func (c *GeminiClient) SendMessage(systemPrompt, userMessage string) (string, error) {
	return c.client.SendMessage(systemPrompt, userMessage)
}

// SendSimpleMessage envía un mensaje sin system prompt.
func (c *GeminiClient) SendSimpleMessage(message string) (string, error) {
	return c.client.SendSimpleMessage(message)
}

// Provider retorna el tipo de provider.
func (c *GeminiClient) Provider() Provider {
	return ProviderGemini
}

// IsAvailable siempre retorna true (si hay API key configurada).
func (c *GeminiClient) IsAvailable() (bool, error) {
	return true, nil
}

// Close cierra el cliente.
func (c *GeminiClient) Close() error {
	return c.client.Close()
}

// GroqClient es un wrapper para el cliente de Groq.
type GroqClient struct {
	client *groq.Client
}

// SendMessage envía un mensaje usando Groq API.
func (c *GroqClient) SendMessage(systemPrompt, userMessage string) (string, error) {
	return c.client.SendMessage(systemPrompt, userMessage)
}

// SendSimpleMessage envía un mensaje sin system prompt.
func (c *GroqClient) SendSimpleMessage(message string) (string, error) {
	return c.client.SendSimpleMessage(message)
}

// Provider retorna el tipo de provider.
func (c *GroqClient) Provider() Provider {
	return ProviderGroq
}

// IsAvailable siempre retorna true (si hay API key configurada).
func (c *GroqClient) IsAvailable() (bool, error) {
	return true, nil
}

// Close cierra el cliente.
func (c *GroqClient) Close() error {
	return c.client.Close()
}
