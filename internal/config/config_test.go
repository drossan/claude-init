package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	require.NoError(t, err)
	assert.Contains(t, path, "claude-init")
	assert.Contains(t, path, "config.yaml")
}

func TestGlobalConfig_SaveAndLoad(t *testing.T) {
	// Crear directorio temporal para tests
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Guardar el path original para restaurarlo después
	originalGetConfigPath := GetConfigPath
	GetConfigPath = func() (string, error) {
		return configPath, nil
	}
	defer func() { GetConfigPath = originalGetConfigPath }()

	t.Run("save and load valid config", func(t *testing.T) {
		config := &GlobalConfig{
			Provider: "openai",
			Providers: map[string]ProviderConfig{
				"openai": {
					APIKey:    "sk-test-key",
					BaseURL:   "https://api.openai.com/v1",
					Model:     "gpt-4o",
					MaxTokens: 4096,
				},
				"claude-api": {
					APIKey:    "sk-ant-test-key",
					Model:     "claude-sonnet-4-20250514",
					MaxTokens: 8192,
				},
			},
		}

		err := config.Save()
		require.NoError(t, err)

		// Verificar que el archivo existe
		_, err = os.Stat(configPath)
		require.NoError(t, err)

		// Cargar la configuración
		loaded, err := Load()
		require.NoError(t, err)
		assert.Equal(t, "openai", loaded.Provider)
		assert.Equal(t, "sk-test-key", loaded.Providers["openai"].APIKey)
		assert.Equal(t, "https://api.openai.com/v1", loaded.Providers["openai"].BaseURL)
		assert.Equal(t, "gpt-4o", loaded.Providers["openai"].Model)
		assert.Equal(t, 4096, loaded.Providers["openai"].MaxTokens)
		assert.Equal(t, "sk-ant-test-key", loaded.Providers["claude-api"].APIKey)
	})

	t.Run("load when file does not exist", func(t *testing.T) {
		// Borrar el archivo si existe
		os.Remove(configPath)

		loaded, err := Load()
		require.NoError(t, err)
		assert.Equal(t, "", loaded.Provider)
		assert.NotNil(t, loaded.Providers)
	})

	t.Run("save creates file with correct permissions", func(t *testing.T) {
		config := &GlobalConfig{
			Provider: "openai",
			Providers: map[string]ProviderConfig{
				"openai": {
					APIKey: "sk-test-key",
				},
			},
		}

		err := config.Save()
		require.NoError(t, err)

		// Verificar permisos del archivo
		info, err := os.Stat(configPath)
		require.NoError(t, err)

		// En Unix, 0600 = rw-------
		assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
	})
}

func TestGlobalConfig_IsProviderConfigured(t *testing.T) {
	config := &GlobalConfig{
		Providers: map[string]ProviderConfig{
			"openai": {
				APIKey: "sk-test-key",
			},
			"claude-api": {
				APIKey: "",
			},
		},
	}

	t.Run("configured provider", func(t *testing.T) {
		assert.True(t, config.IsProviderConfigured("openai"))
	})

	t.Run("empty API key", func(t *testing.T) {
		assert.False(t, config.IsProviderConfigured("claude-api"))
	})

	t.Run("non-existent provider", func(t *testing.T) {
		assert.False(t, config.IsProviderConfigured("zai"))
	})

	t.Run("whitespace only API key", func(t *testing.T) {
		configWithWhitespace := &GlobalConfig{
			Providers: map[string]ProviderConfig{
				"openai": {
					APIKey: "   ",
				},
			},
		}
		assert.False(t, configWithWhitespace.IsProviderConfigured("openai"))
	})
}

func TestGlobalConfig_ProviderDefaults(t *testing.T) {
	t.Run("empty provider defaults to cli", func(t *testing.T) {
		config := &GlobalConfig{}
		assert.Equal(t, "cli", config.GetDefaultProvider())
	})

	t.Run("configured provider", func(t *testing.T) {
		config := &GlobalConfig{Provider: "openai"}
		assert.Equal(t, "openai", config.GetDefaultProvider())
	})
}

func TestGlobalConfig_HasAnyProviderConfigured(t *testing.T) {
	t.Run("no providers configured", func(t *testing.T) {
		config := &GlobalConfig{}
		assert.False(t, config.HasAnyProviderConfigured())
	})

	t.Run("one provider configured", func(t *testing.T) {
		config := &GlobalConfig{
			Providers: map[string]ProviderConfig{
				"openai": {
					APIKey: "sk-test-key",
				},
			},
		}
		assert.True(t, config.HasAnyProviderConfigured())
	})

	t.Run("all providers empty", func(t *testing.T) {
		config := &GlobalConfig{
			Providers: map[string]ProviderConfig{
				"openai":     {APIKey: ""},
				"claude-api": {APIKey: ""},
			},
		}
		assert.False(t, config.HasAnyProviderConfigured())
	})
}
