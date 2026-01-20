# Guía de Desarrollo del CLI - claude-init

Este documento establece los estándares, convenciones y el "contrato" de desarrollo para el proyecto **claude-init CLI**, una herramienta en Go para inicializar proyectos con configuración guiada por IA.

---

## Tabla de Contenidos

1. [Visión General del Proyecto](#visión-general-del-proyecto)
2. [Estructura del Proyecto](#estructura-del-proyecto)
3. [Configuración del Sistema](#configuración-del-sistema)
4. [Integración con APIs de IA](#integración-con-apis-de-ia)
5. [Plantillas Claude (.claude/)](#plantillas-claude-claude)
6. [Flujo de Ejecución del CLI](#flujo-de-ejecución-del-cli)
7. [Patrones de Diseño](#patrones-de-diseño)
8. [Testing](#testing)
9. [Referencias Rápidas](#referencias-rápidas)

---

## Visión General del Proyecto

**claude-init** es un CLI escrito en Go que:

1. **Detecta el tipo de proyecto** (frontend, backend, fullstack, QA, sistemas, etc.)
2. **Hace preguntas interactivas** al usuario para entender el contexto del proyecto
3. **Se conecta a una API de IA** (Claude, OpenAI, z.ai) para generar recomendaciones personalizadas
4. **Crea la estructura `.claude/`** con los archivos necesarios para desarrollo guiado por IA

### Características Principales

- **Detección automática**: Analiza archivos como `package.json`, `go.mod`, `pom.xml`, etc.
- **Configuración global**: Las credenciales de la API se almacenan en un lugar común del sistema
- **Plantillas dinámicas**: Genera skills, commands y agents basados en el stack tecnológico detectado
- **Multi-proveedor**: Soporta Claude (Anthropic), OpenAI, y z.ai

---

## Estructura del Proyecto

```
ia_start/
├── cmd/
│   └── root/
│       └── root.go              # Comando raíz de Cobra
├── internal/
│   ├── config/                  # Gestión de configuración
│   │   ├── config.go           # Estructuras de configuración
│   │   └── loader.go           # Carga de config desde archivo/variables
│   ├── detector/               # Detección de tipo de proyecto
│   │   ├── detector.go         # Lógica principal de detección
│   │   └── patterns.go         # Patrones de detección por stack
│   ├── ai/                     # Clientes de APIs de IA
│   │   ├── client.go           # Interface común
│   │   ├── claude.go           # Cliente Anthropic Claude
│   │   ├── openai.go           # Cliente OpenAI
│   │   └── zai.go              # Cliente z.ai
│   ├── templates/              # Generadores de plantillas .claude
│   │   ├── generator.go        # Lógica principal de generación
│   │   ├── skills.go           # Generador de skills
│   │   ├── commands.go         # Generador de commands
│   │   ├── agents.go           # Generador de agents
│   │   └── guide.go            # Generador de development_guide.md
│   ├── interactive/            # Prompts interactivos
│   │   └── survey.go           # Preguntas al usuario
│   └── models/                 # Modelos de datos
│       ├── project.go          # Tipo de proyecto detectado
│       └── answers.go          # Respuestas del usuario
├── .claude/                    # Plantillas para el propio proyecto
│   ├── agents/
│   ├── skills/
│   ├── commands/
│   └── development_guide.md   # Este archivo
├── go.mod
├── go.sum
└── main.go
```

---

## Configuración del Sistema

### Ubicación de la Configuración

La configuración global se almacena en:

- **Linux/macOS**: `~/.config/claude-init/config.yaml`
- **Windows**: `%APPDATA%\claude-init\config.yaml`

### Estructura del Archivo de Configuración

```yaml
# ~/.config/claude-init/config.yaml
ai:
  provider: claude  # claude | openai | zai
  claude:
    api_key: sk-ant-xxxxx
    model: claude-3-5-sonnet-20241022
    max_tokens: 8192
  openai:
    api_key: sk-xxxxx
    model: gpt-4o
    max_tokens: 4096
  zai:
    api_key: zai-xxxxx
    model: zai-1-xxxxx
    max_tokens: 4096

defaults:
  auto_detect: true
  create_gitignore: true
  overwrite_existing: false
```

### Variables de Entorno (Alternativa)

También se pueden usar variables de entorno:

```bash
export IA_START_AI_PROVIDER=claude
export IA_START_CLAUDE_API_KEY=sk-ant-xxxxx
export IA_START_OPENAI_API_KEY=sk-xxxxx
export IA_START_ZAI_API_KEY=zai-xxxxx
```

**Prioridad**: Variables de entorno > Archivo de configuración

---

## Integración con APIs de IA

### Interface Común

```go
package ai

type AIClient interface {
    // Genera recomendaciones basadas en el contexto del proyecto
    GenerateRecommendations(ctx context.Context, req RecommendationRequest) (*RecommendationResponse, error)

    // Chat interactivo para refinar respuestas
    Chat(ctx context.Context, messages []Message) (*ChatResponse, error)
}

type RecommendationRequest struct {
    ProjectType    string   // frontend, backend, fullstack, etc.
    TechStack      []string // react, vue, node, go, etc.
    UserAnswers    map[string]interface{}
    ExistingFiles  []string // Archivos .claude existentes
}

type RecommendationResponse struct {
    SuggestedSkills    []SkillConfig
    SuggestedCommands  []CommandConfig
    SuggestedAgents    []AgentConfig
    DevelopmentGuide   string
}
```

### Cliente Claude (Anthropic)

```go
type ClaudeClient struct {
    apiKey string
    model  string
    client *http.Client
}

func (c *ClaudeClient) GenerateRecommendations(ctx context.Context, req RecommendationRequest) (*RecommendationResponse, error) {
    // Implementación usando la API de Anthropic
    // POST https://api.anthropic.com/v1/messages
    // Header: x-api-key, anthropic-version, content-type
}
```

### Cliente OpenAI

```go
type OpenAIClient struct {
    apiKey string
    model  string
    client *http.Client
}

func (o *OpenAIClient) GenerateRecommendations(ctx context.Context, req RecommendationRequest) (*RecommendationResponse, error) {
    // Implementación usando la API de OpenAI
    // POST https://api.openai.com/v1/chat/completions
}
```

### Cliente z.ai

```go
type ZAIClient struct {
    apiKey string
    model  string
    client *http.Client
}

func (z *ZAIClient) GenerateRecommendations(ctx context.Context, req RecommendationRequest) (*RecommendationResponse, error) {
    // Implementación usando la API de z.ai
}
```

---

## Plantillas Claude (.claude/)

La estructura `.claude/` generada debe contener:

### agents/

**agentes base** (generados siempre):
- `orchestrator.md` - Orquestador maestro del flujo de trabajo
- `architect.md` - Especialista en arquitectura
- `developer.md` - Desarrollador principal
- `reviewer.md` - Revisor de código
- `tester.md` - Especialista en testing
- `debugger.md` - Especialista en debugging
- `writer.md` - Especialista en documentación
- `planning-agent.md` - Planificador de tareas

### skills/

**skills por stack** (generadas según detección):

**Frontend**:
- `typescript.md` - Tipado estricto con TypeScript
- `react.md` / `vue.md` / `angular.md` / `svelte.md` - Framework específico
- `css-architecture.md` - Arquitectura CSS (BEM, CSS Modules, etc.)
- `state-management.md` - Redux, Zustand, Pinia, etc.

**Backend**:
- `api-rest.md` - Diseño de APIs REST
- `node-js.md` / `go.md` / `python.md` - Lenguaje específico
- `db-expert.md` - Expertise en bases de datos
- `security.md` - Seguridad y autenticación

**Testing**:
- `tdd-champion.md` - Desarrollo guiado por tests
- `vitest.md` / `jest.md` - Framework de testing
- `e2e-testing.md` - Testing end-to-end (Playwright, Cypress)

**DevOps/Sistemas**:
- `docker.md` - Contenedores
- `k8s.md` - Kubernetes
- `ci-cd.md` - CI/CD pipelines

**Comunes**:
- `technical-writer.md` - Documentación técnica
- `code-reviewer.md` - Revisión de código
- `debug-master.md` - Debugging avanzado

### commands/

**commands base**:
- `init.md` - Inicializar/configurar el proyecto
- `new-feature.md` - Crear nueva funcionalidad
- `bug-fix.md` - Corregir errores
- `refactor.md` - Refactorizar código
- `improve-tests.md` - Mejorar tests
- `pre-flight.md` - Verificación antes de commit
- `plan-manage.md` - Gestión de planes
- `orchestrator.md` - Orquestador maestro

### Otros archivos

- `development_guide.md` - Guía de desarrollo personalizada del proyecto
- `settings.local.json` - Configuración local de permisos

---

## Flujo de Ejecución del CLI

### Comando Principal: `claude-init init`

```go
package cmd

var rootCmd = &cobra.Command{
    Use:   "claude-init",
    Short: "CLI para inicializar proyectos con configuración guiada por IA",
}

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Inicializa la configuración .claude en el proyecto actual",
    RunE: func(cmd *cobra.Command, args []string) error {
        return runInit()
    },
}
```

### Flujo Paso a Paso

```
1. Usuario ejecuta: claude-init init
   ↓
2. Detectar tipo de proyecto
   - Buscar package.json, go.mod, requirements.txt, etc.
   - Identificar stack tecnológico
   ↓
3. Hacer preguntas interactivas
   - ¿Qué tipo de proyecto es?
   - ¿Cuál es el objetivo principal?
   - ¿Qué frameworks usas?
   - ¿Necesitas skills específicas?
   ↓
4. (Opcional) Consultar API de IA
   - Enviar contexto del proyecto
   - Recibir recomendaciones personalizadas
   ↓
5. Generar estructura .claude/
   - Crear directorios necesarios
   - Generar skills según stack
   - Generar commands necesarios
   - Generar agents base
   - Crear development_guide.md personalizado
   ↓
6. Confirmar y crear archivos
   - Mostrar resumen de lo que se creará
   - Esperar confirmación del usuario
   - Escribir archivos
```

### Comando de Configuración: `claude-init config`

```go
var configCmd = &cobra.Command{
    Use:   "config",
    Short: "Configura las credenciales de la API de IA",
    RunE: func(cmd *cobra.Command, args []string) error {
        return runConfig()
    },
}
```

---

## Patrones de Diseño

### 1. Builder Pattern para Generación de Plantillas

```go
type TemplateBuilder struct {
    projectType string
    techStack   []string
    aiClient    ai.AIClient
}

func (b *TemplateBuilder) BuildSkills() ([]SkillConfig, error) {
    // Genera skills basadas en el stack
}

func (b *TemplateBuilder) BuildCommands() ([]CommandConfig, error) {
    // Genera commands necesarios
}

func (b *TemplateBuilder) BuildAgents() ([]AgentConfig, error) {
    // Genera agents base
}
```

### 2. Strategy Pattern para Detección de Proyectos

```go
type ProjectDetector interface {
    Detect() (ProjectType, error)
}

type NodeJSDetector struct{}
type GoDetector struct{}
type PythonDetector struct{}

func (d *NodeJSDetector) Detect() (ProjectType, error) {
    // Detecta proyectos Node.js
}
```

### 3. Factory Pattern para Clientes de IA

```go
func NewAIClient(provider, apiKey string) (ai.AIClient, error) {
    switch provider {
    case "claude":
        return ai.NewClaudeClient(apiKey)
    case "openai":
        return ai.NewOpenAIClient(apiKey)
    case "zai":
        return ai.NewZAIClient(apiKey)
    default:
        return nil, fmt.Errorf("proveedor no soportado: %s", provider)
    }
}
```

---

## Metodología de Desarrollo: TDD

**Este proyecto utiliza Test-Driven Development (TDD) como metodología principal.**

### Ciclo TDD (Red-Green-Refactor)

1. **Red**: Escribir un test que falle
   - Definir el comportamiento esperado
   - Ejecutar el test y verificar que falla

2. **Green**: Escribir el código mínimo para pasar el test
   - Implementar solo lo necesario
   - No escribir código de más

3. **Refactor**: Mejorar el código manteniendo los tests verdes
   - Limpiar y optimizar
   - Mantener la cobertura de tests

### Reglas de TDD en este Proyecto

1. **No escribir código de producción sin tests primero**
   - Todo nuevo código debe tener un test que lo justifique
   - Los tests definen la especificación del comportamiento

2. **Escribir tests antes de implementar**
   - Para nueva funcionalidad: tests primero
   - Para bugs: escribir test que reproduzca el bug, luego fix
   - Para refactor: tests de seguridad primero

3. **Cobertura mínima del 80%**
   - Todo paquete debe tener >=80% de cobertura
   - Los paquetes críticos (config, ai, detector) deben tener >=90%

4. **Tests como documentación**
   - Los tests deben ser legibles y documentar el comportamiento
   - Usar nombres descriptivos: `TestLogger_NewCreatesLoggerWithDefaultLevel`

### Convenciones de Testing

```go
// Estructura de un test file: package_test.go
package logger_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLogger_NewCreatesLoggerWithDefaultLevel(t *testing.T) {
    // Arrange
    expectedLevel := logger.INFOLevel

    // Act
    log := logger.New(nil, expectedLevel)

    // Assert
    assert.Equal(t, expectedLevel, log.Level())
}

func TestLogger_DebugOnlyLogsWhenDebugEnabled(t *testing.T) {
    // Arrange
    buf := &bytes.Buffer{}
    log := logger.New(buf, logger.ERRORLevel)

    // Act
    log.Debug("test message")

    // Assert
    assert.Empty(t, buf.String(), "Debug should not log when level is ERROR")
}
```

### Orden de Implementación con TDD

Para cada nueva funcionalidad:

1. **Tester Agent**: Escribir tests que fallen (Red)
2. **Developer Agent**: Implementar código mínimo (Green)
3. **Developer Agent**: Refactorizar si necesario (Refactor)
4. **Reviewer Agent**: Revisar código y tests
5. **Repetir** para la siguiente funcionalidad

---

## Testing

### Testing del CLI

Usar `testify` para assertions y mocking:

```go
package detector_test

func TestNodeJSDetector(t *testing.T) {
    detector := &NodeJSDetector{}
    projectType, err := detector.Detect()

    assert.NoError(t, err)
    assert.Equal(t, ProjectTypeBackend, projectType.Category)
}
```

### Testing de Integración

```bash
# Ejecutar tests
go test ./...

# Con cobertura
go test -cover ./...

# Con race detection
go test -race ./...

# Cobertura detallada
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Referencias Rápidas

### Comandos del CLI

```bash
# Inicializar en directorio actual
claude-init init

# Configurar API
claude-init config set provider claude
claude-init config set api-key sk-ant-xxxxx

# Ver configuración actual
claude-init config get

# Detectar proyecto sin crear archivos
claude-init detect

# Versión
claude-init version
```

### Dependencias Principales de Go

```go
require (
    github.com/spf13/cobra v1.8.0         // CLI framework
    github.com/spf13/viper v1.18.0        // Configuración
    github.com/AlecAivazis/survey/v2 v2.3.7 // Prompts interactivos
    github.com/tidwall/gjson v1.17.0      // Parsing JSON
)
```

### Archivos de Configuración

- **Go**: `go.mod` - Dependencias del módulo
- **CLI**: `~/.config/claude-init/config.yaml` - Configuración global
- **Proyecto**: `.claude/` - Directorio generado por el CLI

---

## Checklist de Calidad

- [ ] ¿El código sigue las convenciones de Go (gofmt, golint)?
- [ ] ¿Los tests tienen cobertura >80%?
- [ ] ¿La configuración se maneja correctamente?
- [ ] ¿Los errores se propagan adecuadamente?
- [ ] ¿La documentación está actualizada?
- [ ] ¿El CLI funciona en Linux, macOS y Windows?
