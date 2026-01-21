package config

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/drossan/claude-init/internal/config"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Configure AI provider settings",
	Long: `Configure AI provider settings for claude-init CLI.

This command allows you to set up API keys for different AI providers.
The configuration is stored in ~/.config/claude-init/config.yaml (macOS/Linux)
or %APPDATA%\claude-init\config.yaml (Windows).

Available providers:
- cli: Claude CLI (free with Claude Code PRO)
- openai: OpenAI API (requires API key)
- gemini: Google Gemini API (free tier available)
- groq: Groq API (free tier available)
- claude-api: Anthropic Claude API (requires API key)
- zai: Z.AI API (requires API key)`,
	RunE: runConfig,
}

var providerFlag string

func init() {
	Cmd.Flags().StringVarP(&providerFlag, "provider", "p", "", "Provider to configure (cli, openai, gemini, groq, claude-api, zai)")
}

func runConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Seleccionar provider
	provider := providerFlag
	if provider == "" {
		provider, err = askProvider()
		if err != nil {
			return fmt.Errorf("error getting provider: %w", err)
		}
	}

	// Si es CLI, no necesita configuración
	if provider == "cli" {
		cfg.SetDefaultProvider("cli")
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}
		fmt.Println("✓ Provider set to: Claude CLI")
		fmt.Println("  No API key needed")
		return nil
	}

	// Para providers de API, solicitar API key
	apiKey, err := askAPIKey(provider)
	if err != nil {
		return fmt.Errorf("error getting API key: %w", err)
	}

	// Valores por defecto según provider
	defaults := getDefaultsForProvider(provider)

	// Opcionalmente preguntar por más configuración
	baseURL := defaults.baseURL
	model := defaults.model
	maxTokens := defaults.maxTokens

	allowMoreConfig := false
	promptMore := &survey.Confirm{
		Message: "Do you want to configure advanced options (base URL, model, max tokens)?",
		Default: false,
	}
	if err := survey.AskOne(promptMore, &allowMoreConfig); err == nil && allowMoreConfig {
		baseURL, _ = askBaseURL(provider, defaults.baseURL)
		model, _ = askModel(provider, defaults.model)
		maxTokens, _ = askMaxTokens(provider, defaults.maxTokens)
	}

	// Guardar configuración
	cfg.SetDefaultProvider(provider)
	cfg.SetProviderConfig(provider, config.ProviderConfig{
		APIKey:    apiKey,
		BaseURL:   baseURL,
		Model:     model,
		MaxTokens: maxTokens,
	})

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("error saving config: %w", err)
	}

	fmt.Printf("✓ Provider configured: %s\n", provider)
	fmt.Printf("  Config file: %s\n", getConfigPathDisplay())
	if provider != "cli" {
		fmt.Println("  You can now use claude-init commands with this provider")
	}
	return nil
}

func askProvider() (string, error) {
	var provider string
	prompt := &survey.Select{
		Message: "Select AI provider to configure:",
		Options: []string{
			"Claude CLI (Free with Claude Code PRO)",
			"OpenAI API",
			"Google Gemini API (Free tier available)",
			"Groq API (Free tier available)",
			// "Claude API (Anthropic API)",
			// "Z.AI API",
		},
		Default: "Claude CLI (Free with Claude Code PRO)",
	}

	if err := survey.AskOne(prompt, &provider); err != nil {
		return "", err
	}

	// Mapear la selección al valor interno
	switch provider {
	case "Claude CLI (Free with Claude Code PRO)":
		return "cli", nil
	case "OpenAI API":
		return "openai", nil
	case "Google Gemini API (Free tier available)":
		return "gemini", nil
	case "Groq API (Free tier available)":
		return "groq", nil
	// Comentado temporalmente - se usará más adelante
	// case "Claude API (Anthropic API)":
	// 	return "claude-api", nil
	// case "Z.AI API":
	// 	return "zai", nil
	default:
		return "cli", nil
	}
}

func askAPIKey(provider string) (string, error) {
	var apiKey string
	prompt := &survey.Input{
		Message: fmt.Sprintf("Enter %s API key:", strings.ToUpper(provider)),
		Help:    "You can get your API key from the provider's dashboard",
	}

	if err := survey.AskOne(prompt, &apiKey, survey.WithValidator(survey.Required)); err != nil {
		return "", err
	}

	return strings.TrimSpace(apiKey), nil
}

func askBaseURL(provider, defaultURL string) (string, error) {
	var baseURL string
	prompt := &survey.Input{
		Message: "Base URL (press Enter for default):",
		Default: defaultURL,
	}

	if err := survey.AskOne(prompt, &baseURL); err != nil {
		return "", err
	}

	if strings.TrimSpace(baseURL) == "" {
		return defaultURL, nil
	}

	return strings.TrimSpace(baseURL), nil
}

func askModel(provider, defaultModel string) (string, error) {
	var model string
	prompt := &survey.Input{
		Message: "Model (press Enter for default):",
		Default: defaultModel,
	}

	if err := survey.AskOne(prompt, &model); err != nil {
		return "", err
	}

	if strings.TrimSpace(model) == "" {
		return defaultModel, nil
	}

	return strings.TrimSpace(model), nil
}

func askMaxTokens(provider string, defaultMaxTokens int) (int, error) {
	var maxTokensStr string
	prompt := &survey.Input{
		Message: "Max tokens (press Enter for default):",
		Default: fmt.Sprintf("%d", defaultMaxTokens),
	}

	if err := survey.AskOne(prompt, &maxTokensStr); err != nil {
		return 0, err
	}

	if strings.TrimSpace(maxTokensStr) == "" {
		return defaultMaxTokens, nil
	}

	var maxTokens int
	if _, err := fmt.Sscanf(maxTokensStr, "%d", &maxTokens); err != nil {
		return 0, fmt.Errorf("invalid number: %w", err)
	}

	return maxTokens, nil
}

type providerDefaults struct {
	baseURL   string
	model     string
	maxTokens int
}

func getDefaultsForProvider(provider string) providerDefaults {
	switch provider {
	case "claude-api":
		return providerDefaults{
			baseURL:   "https://api.anthropic.com/v1/messages",
			model:     "claude-opus-4",
			maxTokens: 200000, // Claude Opus 4 - 200K tokens
		}
	case "openai":
		return providerDefaults{
			baseURL:   "https://api.openai.com/v1",
			model:     "gpt-4o-mini",
			maxTokens: 16384, // GPT-4o-mini soporta 16K completion tokens
		}
	case "zai":
		return providerDefaults{
			baseURL:   "https://api.z.ai/v1",
			model:     "glm-4.7",
			maxTokens: 204800, // GLM-4.7 tiene 204,800 tokens de contexto
		}
	case "gemini":
		return providerDefaults{
			baseURL:   "https://generativelanguage.googleapis.com/v1beta/models",
			model:     "gemini-2.5-flash",
			maxTokens: 1000000, // Gemini 2.5 Flash tiene 1M tokens
		}
	case "groq":
		return providerDefaults{
			baseURL:   "https://api.groq.com/openai/v1",
			model:     "llama-3.3-70b-versatile",
			maxTokens: 32768, // Groq soporta 32K context window
		}
	default:
		return providerDefaults{
			baseURL:   "",
			model:     "",
			maxTokens: 4096,
		}
	}
}

func getConfigPathDisplay() string {
	if path, err := config.GetConfigPath(); err == nil {
		return path
	}
	return "~/.config/claude-init/config.yaml"
}
