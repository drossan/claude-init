# Arquitectura de claude-init CLI

## Visión General

claude-init es una herramienta CLI escrita en Go que inicializa proyectos con configuración optimizada para desarrollo guiado por IA. La arquitectura se diseñó siguiendo principios de simplicidad y claridad, eliminando la complejidad de detección automática en favor de un enfoque guiado por el usuario.

## Principios de Diseño

### 1. Usuario como Fuente de Verdad
El usuario conoce mejor su stack tecnológico que cualquier algoritmo de detección. Por lo tanto, el CLI pregunta directamente al usuario en lugar de adivinar.

### 2. IA como Validador, no Adivinador
La IA se usa para validar la información proporcionada por el usuario y hacer recomendaciones, no para inferir el contexto del proyecto.

### 3. Templates sobre Generación
Es mejor copiar y personalizar templates probados en producción (claude_examples/) que generar código desde cero.

### 4. Simplicidad sobre Complejidad
Eliminamos la detección automática de lenguajes/frameworks para reducir la complejidad del mantenimiento.

## Diagrama de Arquitectura

```
┌─────────────────────────────────────────────────────────────────┐
│                         claude-init CLI                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌────────────────┐      ┌─────────────────┐                   │
│  │  cmd/          │      │  internal/      │                   │
│  │  ├── root/     │──────│  ├── survey/    │                   │
│  │  ├── init/     │      │  ├── ai/        │                   │
│  │  ├── generate/ │      │  ├── claudeexamples/              │
│  │  ├── version/  │      │  ├── config/    │                   │
│  │  └── completion│      │  ├── logger/    │                   │
│  └────────────────┘      │  └── templates/ │                   │
│                          └─────────────────┘                   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Flujo de Ejecución

### Comando `init`

```
┌──────────────────────────────────────────────────────────────────┐
│  1. Usuario ejecuta: claude-init init [path]                        │
└─────────────────────────────┬────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│  2. Validación de path                                           │
│     - Verifica que el path existe                                │
│     - Verifica que no existe .claude/ (o usa --force)            │
└─────────────────────────────┬────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│  3. Survey Interactivo (internal/survey/)                        │
│     - Ejecuta 8 preguntas predefinidas                           │
│     - Valida respuestas requeridas                               │
│     - Retorna struct Answers con toda la información             │
└─────────────────────────────┬────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│  4. Validación con IA (opcional - internal/ai/)                  │
│     - Si se configuró --ai-provider y --api-key                  │
│     - Envía respuestas a la IA                                   │
│     - Recibe información faltante y sugerencias                  │
└─────────────────────────────┬────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│  5. Recomendación de Estructura (opcional)                       │
│     - Si se configuró IA                                         │
│     - Solicita agents, commands y skills recomendados            │
│     - Si no hay IA, usa getDefaultRecommendation()               │
└─────────────────────────────┬────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│  6. Generación de Estructura (internal/claudeexamples/)          │
│     - Carga templates desde claude_examples/                     │
│     - Personaliza con Answers (placeholders)                     │
│     - Genera estructura .claude/                                 │
│       ├── agents/                                                │
│       ├── commands/                                              │
│       ├── skills/                                                │
│       ├── sessions/                                              │
│       └── development_guide.md                                   │
└─────────────────────────────┬────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│  7. Resumen                                                      │
│     - Muestra estructura generada                                │
│     - Muestra recomendaciones de IA (si aplicable)               │
└──────────────────────────────────────────────────────────────────┘
```

## Componentes Principales

### cmd/

Paquete que contiene todos los comandos del CLI.

#### cmd/root/
- **Propósito**: Comando raíz que orquesta todos los subcomandos
- **Responsabilidades**:
  - Registrar subcomandos
  - Configurar logger global
  - Manejar flags persistentes (verbose)

#### cmd/init/
- **Propósito**: Inicializar la estructura .claude/
- **Responsabilidades**:
  - Ejecutar el flujo principal del CLI
  - Coordinar survey, IA y generación
  - Manejar flags del comando init
- **Flujo**:
  1. Validar path
  2. Ejecutar survey
  3. Validar con IA (opcional)
  4. Obtener recomendaciones (opcional)
  5. Generar estructura
  6. Mostrar resumen

#### cmd/generate/
- **Propósito**: Generar configuración sin survey
- **Responsabilidades**:
  - Generar templates embebidos
  - Soportar generación parcial (agents, skills, commands, guides)

#### cmd/version/
- **Propósito**: Mostrar información de versión
- **Responsabilidades**:
  - Mostrar versión, commit, build date
  - Soportar formatos short, json, verbose

#### cmd/completion/
- **Propósito**: Generar scripts de autocompletado
- **Responsabilidades**:
  - Generar completions para bash, zsh, fish, powershell

### internal/survey/

Paquete que implementa el sistema de preguntas interactivas.

#### Componentes

- **Question**: Representa una pregunta del survey
  - `ID`: Identificador único
  - `Text`: Texto de la pregunta
  - `Type`: Tipo (input, select, multiselect, confirm, multiline)
  - `Required`: Si es obligatoria
  - `Options`: Opciones para select/multiselect

- **Answers**: Contiene las respuestas del usuario
  - `ProjectName`: Nombre del proyecto
  - `Description`: Descripción breve
  - `Language`: Lenguaje principal
  - `Framework`: Framework (opcional)
  - `Architecture`: Arquitectura deseada
  - `Database`: Base de datos (opcional)
  - `ProjectType`: Tipo de proyecto
  - `BusinessContext`: Contexto del negocio

- **Runner**: Ejecuta el survey de forma interactiva
  - Usa `github.com/AlecAivazis/survey/v2`
  - Mapea respuestas a struct Answers
  - Valida campos requeridos

#### Preguntas Predefinidas

```go
GetProjectQuestions() []*Question {
    return []*Question{
        {ID: "project_name", Text: "Nombre del proyecto:", ...},
        {ID: "description", Text: "Descripción breve:", ...},
        {ID: "language", Text: "Lenguaje principal:", Options: [...]},
        {ID: "framework", Text: "Framework (si aplica):", Options: [...]},
        {ID: "architecture", Text: "Arquitectura deseada:", Options: [...]},
        {ID: "database", Text: "Base de datos (si aplica):", Options: [...]},
        {ID: "project_type", Text: "Tipo de proyecto:", Options: [...]},
        {ID: "business_context", Text: "Contexto del negocio:", ...},
    }
}
```

### internal/ai/

Paquete que implementa clientes para interactuar con APIs de IA.

#### Interfaz AIClient

```go
type AIClient interface {
    // ValidateAnswers valida las respuestas y detecta información faltante
    ValidateAnswers(ctx context.Context, answers *survey.Answers) (*ValidationResult, error)

    // RecommendStructure recomienda qué agents, commands y skills generar
    RecommendStructure(ctx context.Context, answers *survey.Answers) (*Recommendation, error)

    // Provider retorna el nombre del proveedor
    Provider() string
}
```

#### Implementaciones

- **Claude**: Cliente para Anthropic Claude
- **OpenAI**: Cliente para OpenAI GPT
- **z.ai**: Cliente para z.ai

#### ValidationResult

Contiene el resultado de validar las respuestas:

```go
type ValidationResult struct {
    IsValid     bool     // true si las respuestas son válidas
    MissingInfo []string // información faltante detectada
    Suggestions []string // sugerencias de mejora
    Questions   []string // preguntas adicionales a hacer
}
```

#### Recommendation

Contiene la recomendación de estructura:

```go
type Recommendation struct {
    Agents      []string // agents a generar
    Commands    []string // commands a generar
    Skills      []string // skills a incluir
    Description string   // descripción de la estructura recomendada
}
```

### internal/claudeexamples/

Paquete que carga y genera templates desde `claude_examples/`.

#### Loader

Interfaz para cargar ejemplos desde el directorio:

```go
type Loader interface {
    LoadAgents() ([]*ClaudeExample, error)
    LoadCommands() ([]*ClaudeExample, error)
    LoadSkills() ([]*ClaudeExample, error)
    LoadDevelopmentGuide() (*ClaudeExample, error)
}
```

#### ClaudeExample

Representa un ejemplo de configuración:

```go
type ClaudeExample struct {
    Type     string            // "agent", "command", "skill", o "guide"
    Name     string            // nombre del archivo sin extensión
    Source   string            // contenido del archivo
    Metadata map[string]string // metadata del frontmatter
}
```

#### Generator

Interfaz para generar archivos:

```go
type Generator interface {
    Generate(outputDir string, examples []*ClaudeExample, answers *survey.Answers) error
    Customize(example *ClaudeExample, answers *survey.Answers) (string, error)
}
```

#### Personalización

El generator reemplaza placeholders en los templates:

```
{{ProjectName}}     → Nombre del proyecto
{{Description}}     → Descripción del proyecto
{{Language}}        → Lenguaje principal
{{Framework}}       → Framework
{{Architecture}}    → Arquitectura
{{Database}}        → Base de datos
{{ProjectType}}     → Tipo de proyecto
{{BusinessContext}} → Contexto del negocio
```

### internal/config/

Paquete para gestión de configuración.

#### Config

Estructura de configuración:

```go
type Config struct {
    AI       AIConfig
    Defaults DefaultsConfig
}

type AIConfig struct {
    Provider string
    Claude   ClaudeConfig
    OpenAI   OpenAIConfig
    ZAI      ZAIConfig
}

type DefaultsConfig struct {
    CreateGitignore     bool
    OverwriteExisting   bool
}
```

#### Ubicación

La configuración se almacena en `~/.config/claude-init/config.yaml`.

### internal/logger/

Paquete para logging estructurado.

#### Logger

```go
type Logger struct {
    out      io.Writer
    minLevel Level
}

func (l *Logger) Debug(format string, args ...interface{})
func (l *Logger) Info(format string, args ...interface{})
func (l *Logger) Warn(format string, args ...interface{})
func (l *Logger) Error(format string, args ...interface{})
```

#### Niveles

- DEBUG: Información detallada para debugging
- INFO: Información general
- WARN: Advertencias
- ERROR: Errores

### internal/templates/

Paquete para generación de templates embebidos.

#### TemplateContext

Contexto para la generación:

```go
type TemplateContext struct {
    ProjectName string
    Language    string
    Path        string
    Description string
    Framework   string
    // ... más campos
}
```

#### Generator

Genera diferentes tipos de templates:

- Agents: Architect, Developer, Tester, Writer
- Skills: Language-specific, framework-specific
- Commands: Build, test, run, etc.
- Guides: Development guide

## Estructura de Directorios

```
.claude/
├── agents/              # Configuraciones de agentes de IA
│   ├── architect.md    # Agente especializado en arquitectura
│   ├── developer.md    # Agente especializado en desarrollo
│   ├── tester.md       # Agente especializado en testing
│   └── writer.md       # Agente especializado en documentación
├── commands/            # Comandos personalizados para Claude Code
│   ├── build.md        # Comando para build del proyecto
│   ├── test.md         # Comando para ejecutar tests
│   └── run.md          # Comando para ejecutar el proyecto
├── skills/              # Skills específicas del proyecto
│   ├── language_go.md  # Skill específica de Go
│   └── testing.md      # Skill de testing
├── sessions/            # Sesiones de Claude Code (creadas por el usuario)
└── development_guide.md # Guía de desarrollo del proyecto
```

## Claude Examples/

El directorio `claude_examples/` contiene templates base que se copian y personalizan:

```
claude_examples/
├── agents/              # Templates de agentes
│   ├── architect.md
│   ├── developer.md
│   ├── tester.md
│   └── writer.md
├── commands/            # Templates de comandos
│   ├── build.md
│   ├── test.md
│   └── run.md
├── skills/              # Templates de skills
│   ├── go.md
│   ├── nodejs.md
│   └── python.md
└── development_guide.md # Template de guía de desarrollo
```

## Integración con Claude Code

La estructura generada está optimizada para usarse con [Claude Code](https://claude.ai/code):

### Agents

Los agents definen roles especializados que Claude puede adoptar:

- **Architect**: Especialista en arquitectura de software
- **Developer**: Especialista en desarrollo e implementación
- **Tester**: Especialista en testing y calidad
- **Writer**: Especialista en documentación

### Commands

Los commands definen comandos personalizados que Claude puede ejecutar:

- **build**: Compila el proyecto
- **test**: Ejecuta los tests
- **run**: Ejecuta el proyecto

### Skills

Las skills definen conocimientos especializados:

- **Language-specific**: Go, Node.js, Python, etc.
- **Framework-specific**: Express, Django, Gin, etc.
- **Domain-specific**: Testing, arquitectura, etc.

### Development Guide

La guía de desarrollo contiene:

- Descripción del proyecto
- Stack tecnológico
- Patrones de arquitectura
- Convenciones de código
- Comandos útiles
- Información de contacto

## Testing

### Estrategia de Testing

El proyecto sigue una estrategia de testing comprehensiva:

1. **Unit Tests**: Tests de componentes individuales
2. **Integration Tests**: Tests de flujos completos
3. **Table-Driven Tests**: Tests con múltiples casos

### Cobertura

El objetivo es mantener >80% de cobertura global.

### Ejecución

```bash
# Todos los tests
go test ./...

# Con coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Con race detector
go test -race ./...
```

## Rendimiento

### Consideraciones

- **Survey**: Ejecución síncrona (bloqueante por diseño)
- **IA API**: Llamadas asíncronas con context timeout
- **I/O**: Uso de buffers para lectura/escritura de archivos

### Optimizaciones

- Carga de templates en memoria
- Reutilización de conexiones HTTP
- Paralelización de generación de archivos

## Seguridad

### API Keys

- Las API keys se almacenan en `~/.config/claude-init/config.yaml`
- Se pueden sobrescribir con variables de entorno
- Nunca se loguean o muestran en output

### Validación de Input

- Todas las respuestas del usuario se validan
- Los paths se sanitizan antes de usar
- Se previene path traversal

### Permisos

- Solo crea archivos en el directorio del proyecto
- Usa permisos 0644 para archivos, 0755 para directorios
- Respeta archivos existentes (requiere --force)

## Mantenibilidad

### Principios

1. **Código Limpio**: Código legible y auto-documentado
2. **Separación de Responsabilidades**: Cada paquete tiene una responsabilidad clara
3. **Interfaces**: Uso de interfaces para desacoplar componentes
4. **Tests**: Tests comprehensivos para facilitar refactorizaciones

### Documentación

- **Godoc**: Documentación de código en comentarios
- **README.md**: Guía de usuario
- **ARCHITECTURE.md**: Este documento
- **CHANGELOG.md**: Historial de cambios

## Extensibilidad

### Agregar Nuevos Proveedores de IA

Para agregar un nuevo proveedor de IA:

1. Implementar la interfaz `AIClient` en `internal/ai/`
2. Agregar configuración en `internal/config/`
3. Actualizar documentación

### Agregar Nuevas Preguntas al Survey

Para agregar nuevas preguntas:

1. Agregar campo a `Answers` en `internal/survey/survey.go`
2. Agregar pregunta en `GetProjectQuestions()`
3. Actualizar método `setAnswer()` en `runner.go`
4. Actualizar placeholders en `internal/claudeexamples/`

### Agregar Nuevos Templates

Para agregar nuevos templates:

1. Crear archivo en `claude_examples/`
2. Agregar metadata en frontmatter si es necesario
3. Actualizar `loader.go` si es un nuevo tipo

## Referencias

- [Cobra](https://github.com/spf13/cobra): Framework de CLI
- [Survey](https://github.com/AlecAivazis/survey/v2): Prompts interactivos
- [Claude Code](https://claude.ai/code): Asistente de desarrollo con IA
- [Effective Go](https://go.dev/doc/effective_go): Guía de estilo de Go
