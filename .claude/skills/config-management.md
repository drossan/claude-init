# Skill: Config Management (config-management)

## Propósito
Especialidad en gestionar configuraciones de aplicaciones CLI en Go, incluyendo archivos YAML, variables de entorno y configuración global del sistema.

## Responsabilidades
- Leer configuraciones desde archivos YAML
- Leer configuraciones desde variables de entorno
- Combinar múltiples fuentes de configuración con prioridad
- Validar configuraciones antes de usarlas
- Crear configuraciones por defecto sensatas

## Estructura de Configuración

### Struct de Configuración

```go
package config

import (
    "os"
    "time"
)

// Config representa la configuración completa del CLI
type Config struct {
    AI      AIConfig      `yaml:"ai"`
    Defaults Defaults     `yaml:"defaults"`
}

// AIConfig configura el proveedor de IA
type AIConfig struct {
    Provider string          `yaml:"provider" env:"IA_START_AI_PROVIDER"` // claude, openai, zai
    Claude   ProviderConfig  `yaml:"claude"`
    OpenAI   ProviderConfig  `yaml:"openai"`
    ZAI      ProviderConfig  `yaml:"zai"`
}

// ProviderConfig es la configuración de un proveedor específico
type ProviderConfig struct {
    APIKey    string        `yaml:"api_key" env:"IA_START_{{.Provider}}_API_KEY"`
    Model     string        `yaml:"model" env:"IA_START_{{.Provider}}_MODEL"`
    MaxTokens int           `yaml:"max_tokens" env:"IA_START_{{.Provider}}_MAX_TOKENS"`
    Timeout   time.Duration `yaml:"timeout" env:"IA_START_{{.Provider}}_TIMEOUT"`
}

// Defaults configura los valores por defecto
type Defaults struct {
    AutoDetect         bool `yaml:"auto_detect" env:"IA_START_AUTO_DETECT"`
    CreateGitignore    bool `yaml:"create_gitignore" env:"IA_START_CREATE_GITIGNORE"`
    OverwriteExisting  bool `yaml:"overwrite_existing" env:"IA_START_OVERWRITE"`
}
```

## Lectura de Configuración

### Ubicación de Archivos

```go
package config

import (
    "os"
    "path/filepath"
)

// GetConfigPath retorna la ruta al archivo de configuración
func GetConfigPath() (string, error) {
    // Usar XDG_CONFIG_DIR o ~/.config
    configDir := os.Getenv("XDG_CONFIG_HOME")
    if configDir == "" {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return "", err
        }
        configDir = filepath.Join(homeDir, ".config")
    }

    return filepath.Join(configDir, "claude-init", "config.yaml"), nil
}

// GetConfigDir retorna el directorio de configuración
func GetConfigDir() (string, error) {
    configPath, err := GetConfigPath()
    if err != nil {
        return "", err
    }
    return filepath.Dir(configPath), nil
}
```

### Loader de Configuración

```go
package config

import (
    "fmt"
    "os"

    "github.com/spf13/viper"
)

type Loader struct {
    configPath string
    viper      *viper.Viper
}

func NewLoader() (*Loader, error) {
    configPath, err := GetConfigPath()
    if err != nil {
        return nil, err
    }

    v := viper.New()
    v.SetConfigFile(configPath)
    v.SetConfigType("yaml")

    // Configurar variables de entorno
    v.SetEnvPrefix("IA_START")
    v.AutomaticEnv()

    return &Loader{
        configPath: configPath,
        viper:      v,
    }, nil
}

func (l *Loader) Load() (*Config, error) {
    config := &Config{}

    // Leer archivo si existe
    if _, err := os.Stat(l.configPath); err == nil {
        if err := l.viper.ReadInConfig(); err != nil {
            return nil, fmt.Errorf("failed to read config: %w", err)
        }
    }

    // Unmarshal
    if err := l.viper.Unmarshal(config); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }

    // Aplicar defaults
    l.applyDefaults(config)

    // Validar
    if err := config.Validate(); err != nil {
        return nil, err
    }

    return config, nil
}

func (l *Loader) applyDefaults(config *Config) {
    // Defaults para AI
    if config.AI.Provider == "" {
        config.AI.Provider = "claude"
    }

    // Defaults para Claude
    if config.AI.Claude.Model == "" {
        config.AI.Claude.Model = "claude-3-5-sonnet-20241022"
    }
    if config.AI.Claude.MaxTokens == 0 {
        config.AI.Claude.MaxTokens = 8192
    }
    if config.AI.Claude.Timeout == 0 {
        config.AI.Claude.Timeout = 30 * time.Second
    }

    // Defaults generales
    config.Defaults.AutoDetect = true
    config.Defaults.CreateGitignore = true
}
```

### Validación de Configuración

```go
func (c *Config) Validate() error {
    // Validar provider
    switch c.AI.Provider {
    case "claude", "openai", "zai":
        // OK
    default:
        return fmt.Errorf("invalid provider: %s (must be claude, openai, or zai)", c.AI.Provider)
    }

    // Validar que el provider seleccionado tenga API key
    providerConfig, err := c.AI.GetProviderConfig()
    if err != nil {
        return err
    }

    if providerConfig.APIKey == "" {
        return fmt.Errorf("API key not set for provider %s", c.AI.Provider)
    }

    return nil
}

func (a AIConfig) GetProviderConfig() (ProviderConfig, error) {
    switch a.Provider {
    case "claude":
        return a.Claude, nil
    case "openai":
        return a.OpenAI, nil
    case "zai":
        return a.ZAI, nil
    default:
        return ProviderConfig{}, fmt.Errorf("unknown provider: %s", a.Provider)
    }
}
```

## Escritura de Configuración

```go
func (l *Loader) Save(config *Config) error {
    // Validar antes de guardar
    if err := config.Validate(); err != nil {
        return err
    }

    // Crear directorio si no existe
    configDir := filepath.Dir(l.configPath)
    if err := os.MkdirAll(configDir, 0755); err != nil {
        return fmt.Errorf("failed to create config dir: %w", err)
    }

    // Configurar viper para escribir
    l.viper.Set("ai", config.AI)
    l.viper.Set("defaults", config.Defaults)

    // Escribir archivo
    if err := l.viper.WriteConfigAs(l.configPath); err != nil {
        return fmt.Errorf("failed to write config: %w", err)
    }

    return nil
}
```

## Comandos de Configuración

### Comando Config Set

```go
var configSetCmd = &cobra.Command{
    Use:   "set <key> <value>",
    Short: "Establece un valor de configuración",
    Args:  cobra.ExactArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
        key := args[0]
        value := args[1]

        loader, _ := NewLoader()
        config, _ := loader.Load()

        // Establecer valor
        switch key {
        case "provider":
            config.AI.Provider = value
        case "api-key":
            provider := config.AI.Provider
            switch provider {
            case "claude":
                config.AI.Claude.APIKey = value
            case "openai":
                config.AI.OpenAI.APIKey = value
            case "zai":
                config.AI.ZAI.APIKey = value
            }
        case "model":
            provider := config.AI.Provider
            switch provider {
            case "claude":
                config.AI.Claude.Model = value
            case "openai":
                config.AI.OpenAI.Model = value
            case "zai":
                config.AI.ZAI.Model = value
            }
        default:
            return fmt.Errorf("unknown key: %s", key)
        }

        // Guardar
        return loader.Save(config)
    },
}
```

### Comando Config Get

```go
var configGetCmd = &cobra.Command{
    Use:   "get [key]",
    Short: "Muestra la configuración actual",
    RunE: func(cmd *cobra.Command, args []string) error {
        loader, _ := NewLoader()
        config, _ := loader.Load()

        if len(args) == 0 {
            // Mostrar toda la configuración
            return printConfig(config)
        }

        // Mostrar valor específico
        key := args[0]
        value, err := getConfigValue(config, key)
        if err != nil {
            return err
        }

        fmt.Println(value)
        return nil
    },
}

func printConfig(config *Config) error {
    data, err := yaml.Marshal(config)
    if err != nil {
        return err
    }
    fmt.Println(string(data))
    return nil
}
```

### Comando Config Unset

```go
var configUnsetCmd = &cobra.Command{
    Use:   "unset <key>",
    Short: "Elimina un valor de configuración (vuelve al default)",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        key := args[0]

        loader, _ := NewLoader()
        config, _ := loader.Load()

        // Eliminar valor (se usará el default)
        switch key {
        case "api-key":
            provider := config.AI.Provider
            switch provider {
            case "claude":
                config.AI.Claude.APIKey = ""
            case "openai":
                config.AI.OpenAI.APIKey = ""
            case "zai":
                config.AI.ZAI.APIKey = ""
            }
        default:
            return fmt.Errorf("cannot unset key: %s", key)
        }

        return loader.Save(config)
    },
}
```

## Prioridad de Configuración

```
Variables de Entorno > Archivo de Configuración > Defaults
```

```go
func (l *Loader) LoadWithOverrides() (*Config, error) {
    // Primero cargar desde archivo
    config, err := l.Load()
    if err != nil {
        return nil, err
    }

    // Luego sobrescribir con variables de entorno
    if provider := os.Getenv("IA_START_AI_PROVIDER"); provider != "" {
        config.AI.Provider = provider
    }

    // API key específica del provider
    switch config.AI.Provider {
    case "claude":
        if key := os.Getenv("IA_START_CLAUDE_API_KEY"); key != "" {
            config.AI.Claude.APIKey = key
        }
    case "openai":
        if key := os.Getenv("IA_START_OPENAI_API_KEY"); key != "" {
            config.AI.OpenAI.APIKey = key
        }
    case "zai":
        if key := os.Getenv("IA_START_ZAI_API_KEY"); key != "" {
            config.AI.ZAI.APIKey = key
        }
    }

    return config, nil
}
```

## Checklist de Config Management

- [ ] Las configuraciones tienen valores por defecto sensatos
- [ ] Las variables de entorno tienen prefijo consistente
- [ ] Los archivos se crean en la ubicación correcta (XDG)
- [ ] Las configuraciones se validan antes de usar
- [ ] Los errores de lectura son claros
- [ ] La prioridad de configuración está documentada
- [ ] Los secrets (API keys) no se loguean
- [ ] Los archivos de configuración tienen permisos correctos (0600)
