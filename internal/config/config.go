package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// configPathFunc permite mockear GetConfigPath en tests.
var configPathFunc = defaultGetConfigPath

// GlobalConfig representa la configuración global del CLI.
type GlobalConfig struct {
	Provider  string                    `yaml:"provider"`
	Providers map[string]ProviderConfig `yaml:"providers"`
}

// ProviderConfig contiene la configuración de un provider específico.
type ProviderConfig struct {
	APIKey    string `yaml:"api_key"`
	BaseURL   string `yaml:"base_url,omitempty"`
	Model     string `yaml:"model,omitempty"`
	MaxTokens int    `yaml:"max_tokens,omitempty"`
}

// defaultGetConfigPath implementa la lógica por defecto para obtener el path de configuración.
func defaultGetConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("error getting current user: %w", err)
	}

	var configDir string
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		configDir = filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "claude-init")
	} else if os.Getenv("APPDATA") != "" {
		// Windows
		configDir = filepath.Join(os.Getenv("APPDATA"), "claude-init")
	} else {
		// macOS y Linux
		configDir = filepath.Join(usr.HomeDir, ".config", "claude-init")
	}

	// Asegurar que el directorio existe
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("error creating config directory: %w", err)
	}

	return filepath.Join(configDir, "config.yaml"), nil
}

// GetConfigPath retorna el path del archivo de configuración global.
// Sigue el estándar XDG Base Directory Specification.
// macOS: ~/.config/claude-init/config.yaml
// Linux: ~/.config/claude-init/config.yaml
// Windows: %APPDATA%\claude-init\config.yaml
func GetConfigPath() (string, error) {
	return configPathFunc()
}

// Load carga la configuración global desde disco.
func Load() (*GlobalConfig, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Si el archivo no existe, retornar configuración vacía
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &GlobalConfig{
			Providers: make(map[string]ProviderConfig),
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config GlobalConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Inicializar map si es nil
	if config.Providers == nil {
		config.Providers = make(map[string]ProviderConfig)
	}

	return &config, nil
}

// Save guarda la configuración global en disco.
func (c *GlobalConfig) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Inicializar map si es nil
	if c.Providers == nil {
		c.Providers = make(map[string]ProviderConfig)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	// Crear archivo con permisos restrictivos (solo lectura para el usuario)
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// SetProviderConfig configura un provider específico.
func (c *GlobalConfig) SetProviderConfig(provider string, config ProviderConfig) {
	if c.Providers == nil {
		c.Providers = make(map[string]ProviderConfig)
	}
	c.Providers[provider] = config
}

// GetProviderConfig retorna la configuración de un provider específico.
func (c *GlobalConfig) GetProviderConfig(provider string) (ProviderConfig, bool) {
	if c.Providers == nil {
		return ProviderConfig{}, false
	}
	config, exists := c.Providers[provider]
	return config, exists
}

// IsProviderConfigured retorna true si un provider tiene API key configurada.
func (c *GlobalConfig) IsProviderConfigured(provider string) bool {
	config, exists := c.GetProviderConfig(provider)
	return exists && strings.TrimSpace(config.APIKey) != ""
}

// SetDefaultProvider establece el provider por defecto.
func (c *GlobalConfig) SetDefaultProvider(provider string) {
	c.Provider = provider
}

// GetDefaultProvider retorna el provider por defecto, o "cli" si no está configurado.
func (c *GlobalConfig) GetDefaultProvider() string {
	if c.Provider == "" {
		return "cli"
	}
	return c.Provider
}

// HasAnyProviderConfigured retorna true si al menos un provider tiene API key.
func (c *GlobalConfig) HasAnyProviderConfigured() bool {
	if c.Providers == nil {
		return false
	}
	for _, config := range c.Providers {
		if strings.TrimSpace(config.APIKey) != "" {
			return true
		}
	}
	return false
}

// EnsureConfigured verifica que el provider por defecto existe y está configurado.
// Si no hay provider configurado, retorna error.
func EnsureConfigured() (*GlobalConfig, error) {
	config, err := Load()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	provider := config.GetDefaultProvider()

	// Si es CLI, no necesita API key
	if provider == "cli" {
		return config, nil
	}

	// Si es un provider de API, verificar que tiene API key
	if !config.IsProviderConfigured(provider) {
		return nil, fmt.Errorf("provider '%s' not configured. Please run: claude-init config --provider %s", provider, provider)
	}

	return config, nil
}
