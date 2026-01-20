# claude-init

> CLI para inicializar proyectos con configuración guiada por IA

[![Go Report Card](https://goreportcard.com/badge/github.com/drossan/claude-init)](https://goreportcard.com/report/github.com/drossan/claude-init)
[![Coverage](https://img.shields.io/badge/coverage-83.6%25-brightgreen)](https://github.com/drossan/claude-init)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://go.dev/)

**claude-init** es una herramienta CLI escrita en Go que inicializa proyectos con configuración optimizada para
desarrollo guiado por IA con Claude Code. A través de un survey interactivo, recopila información sobre tu proyecto,
opcionalmente valida con IA, y genera la estructura `.claude/` necesaria.

## Características

- **Soporte para Proyectos Nuevos y Existentes**:
  - Proyectos nuevos: Survey interactivo de 8 preguntas
  - Proyectos existentes: Análisis automático de la estructura del proyecto
- **Múltiples Proveedores de IA**:
  - **Claude CLI**: Gratis con Claude Code PRO (opción por defecto)
  - **Claude API**: Anthropic Claude API
  - **OpenAI API**: GPT-4o y otros modelos
  - **Z.AI API**: Models zai
- **Análisis Automático**: Detecta lenguaje, framework, arquitectura y más
- **Recomendaciones Inteligentes**: Sugiere agents, commands y skills basados en tu contexto
- **Modo Dry Run**: Previsualiza qué se generará antes de crear archivos
- **Alta Cobertura de Tests**: Código de alta calidad con ~84% de coverage
- **Configuración Flexible**: Soporte para directorios de configuración custom

## Instalación

### Desde código fuente

```bash
git clone https://github.com/drossan/claude-init.git
cd claude-init
make build
sudo mv bin/claude-init /usr/local/bin/
```

### Usando go install

```bash
go install github.com/drossan/claude-init@latest
```

### Desde binarios precompilados

Los binarios precompilados están disponibles en la página
de [releases](https://github.com/drossan/claude-init/releases) para:

- Linux (amd64)
- macOS (amd64, arm64)
- Windows (amd64)

### Homebrew (macOS/Linux)

Puedes instalar claude-init usando [Homebrew](https://brew.sh/):

```bash
# Añadir el tap
brew tap drossan/homebrew-tools

# Instalar
brew install claude-init

# Actualizar
brew upgrade claude-init
```

El Homebrew formula se actualiza automáticamente con cada release.

## Uso Rápido

```bash
# Inicializar en el directorio actual (el wizard te guiará)
claude-init init

# Inicializar en un path específico
claude-init init /path/to/project

# Configurar proveedor de IA
claude-init config --provider claude-api

# Configurar proveedor CLI (gratis, requiere Claude Code PRO)
claude-init config --provider cli

# Mostrar versión
claude-init version

# Habilitar autocompletado (bash)
claude-init completion bash > /etc/bash_completion.d/claude-init
```

## Flujo de Trabajo

El comando `claude-init init` sigue este flujo:

1. **Selección de Proveedor de IA**: Elige entre Claude CLI (gratis con Claude Code PRO) o APIs de IA
2. **Configuración** (si es necesario): Si eliges una API, el wizard te guía para ingresar tu API key
3. **Origen del Proyecto**: Indica si es un proyecto nuevo o existente
4. **Survey Interactivo**: Responde preguntas sobre tu proyecto (para proyectos nuevos) o confirma el análisis automático (para proyectos existentes)
5. **Generación**: Crea la estructura `.claude/` con agents, skills y commands personalizados

### Preguntas del Survey

El CLI te hará las siguientes preguntas:

1. **Nombre del proyecto**: Identificador único
2. **Descripción breve**: Resumen del propósito
3. **Lenguaje principal**: Go, Node.js, Python, Rust, etc.
4. **Framework** (opcional): Express, NestJS, Django, Gin, etc.
5. **Arquitectura deseada**: Monolito, Microservicios, Hexagonal, Clean, DDD, etc.
6. **Base de datos** (opcional): PostgreSQL, MongoDB, MySQL, etc.
7. **Categoría del proyecto**: API REST, Web App, CLI, Library, etc.
8. **Contexto del negocio**: Descripción detallada del dominio

**Para proyectos existentes**, el CLI también puede:
- Analizar automáticamente la estructura del proyecto
- Detectar el lenguaje, framework y arquitectura
- Preguntar por directorios de documentación adicionales

## Comandos

### init

Inicializa el directorio `.claude/` con agents, skills, commands y guías de desarrollo.

```bash
claude-init init [path] [flags]
```

**Flags:**

- `-f, --force`: Sobrescribe archivos existentes
- `--dry-run`: Muestra qué se generaría sin crear archivos
- `--config-dir`: Directorio de configuración (default: `.claude`)

**Ejemplos:**

```bash
# Inicializar en el directorio actual (wizard interactivo)
claude-init init

# Inicializar en un path específico
claude-init init /path/to/project

# Dry run para ver qué se generaría
claude-init init --dry-run

# Directorio de configuración custom
claude-init init --config-dir .ai-config
```

**Qué hace:**

1. Ejecuta un survey interactivo con 8 preguntas
2. Opcionalmente valida las respuestas con IA
3. Opcionalmente obtiene recomendaciones de estructura
4. Genera la estructura `.claude/`:
    - `agents/`: Configuraciones de agentes
    - `skills/`: Skills específicas del lenguaje/framework
    - `commands/`: Comandos personalizados
    - `development_guide.md`: Guía de desarrollo del proyecto
    - `.gitignore`: Configurado para ignorar archivos sensibles

### config

Configura los proveedores de IA (Claude CLI, Claude API, OpenAI, Z.AI).

```bash
claude-init config [flags]
```

**Flags:**

- `-p, --provider`: Proveedor a configurar (cli, claude-api, openai, zai)

**Proveedores Disponibles:**

- `cli`: Claude CLI (gratis con Claude Code PRO) - **Opción por defecto**
- `claude-api`: Anthropic Claude API (requiere API key)
- `openai`: OpenAI API (requiere API key)
- `zai`: Z.AI API (requiere API key)

**Ejemplos:**

```bash
# Configurar usando wizard interactivo (recomendado)
claude-init config

# Configurar Claude CLI directamente
claude-init config --provider cli

# Configurar Claude API
claude-init config --provider claude-api

# Configurar OpenAI
claude-init config --provider openai

# Configurar Z.AI
claude-init config --provider zai
```

**Wizard de Configuración:**

El comando `config` iniciará un wizard interactivo que te guiará paso a paso:

1. **Selección de proveedor**: Elige entre Claude CLI, Claude API, OpenAI o Z.AI
2. **API Key** (si aplica): Ingresa tu API key de forma segura
3. **Configuración avanzada** (opcional): Base URL, modelo, max tokens

La configuración se guarda en `~/.config/claude-init/config.yaml` (macOS/Linux) o `%APPDATA%\claude-init\config.yaml` (Windows).

### generate

Genera archivos de configuración para Claude Code.

```bash
claude-init generate [project-path] [flags]
```

**Flags:**

- `-f, --force`: Sobrescribe archivos existentes
- `--dry-run`: Muestra qué se generaría sin crear archivos
- `--config-dir`: Directorio de configuración (default: `.claude/`)
- `--output-dir`: Directorio de salida (default: `<project>/.claude/`)
- `--only-agents`: Genera solo los agentes
- `--only-skills`: Genera solo las skills
- `--only-commands`: Genera solo los comandos
- `--only-guides`: Genera solo las guías

**Ejemplos:**

```bash
# Generar toda la configuración
claude-init generate

# Generar solo agentes y skills
claude-init generate --only-agents --only-skills

# Dry run
claude-init generate --dry-run

# Directorio de salida custom
claude-init generate --output-dir /custom/path
```

### version

Muestra información de la versión.

```bash
claude-init version [flags]
```

**Flags:**

- `-s, --short`: Muestra solo el número de versión
- `-j, --json`: Output en formato JSON
- `-v, --verbose`: Muestra información detallada (commit, build date)

**Ejemplos:**

```bash
$ claude-init version
claude-init version 0.1.0

$ claude-init version --short
0.1.0

$ claude-init version --verbose
claude-init version 0.1.0
commit: abc123
built at: 2026-01-17
```

### completion

Genera scripts de autocompletado para shells.

```bash
claude-init completion [bash|zsh|fish|powershell]
```

**Ejemplos:**

```bash
# Bash (Linux)
claude-init completion bash > /etc/bash_completion.d/claude-init

# Bash (macOS)
claude-init completion bash > /usr/local/etc/bash_completion.d/claude-init

# Zsh
claude-init completion zsh > "${fpath[1]}/_claude-init"

# Fish
claude-init completion fish > ~/.config/fish/completions/claude-init.fish

# PowerShell
claude-init completion powershell > claude-init.ps1
```

## Configuración

La configuración de IA se almacena en `~/.config/claude-init/config.yaml`:

```yaml
# Proveedor por defecto
provider: cli

# Configuración de proveedores
providers:
  cli:
    # Claude CLI no requiere API key, solo Claude Code PRO
  claude-api:
    api_key: sk-ant-xxxxx
    base_url: https://api.anthropic.com/v1/messages
    model: claude-sonnet-4-20250514
    max_tokens: 8192
  openai:
    api_key: sk-xxxxx
    base_url: https://api.openai.com/v1
    model: gpt-4o
    max_tokens: 4096
  zai:
    api_key: zai-xxxxx
    base_url: https://api.z.ai/v1
    model: zai-large
    max_tokens: 4096
```

### Notas Importantes

- **Claude CLI** (`cli`): Es la opción por defecto y gratuita si tienes Claude Code PRO. No requiere API key.
- **APIs de IA**: Requieren una suscripción activa y API key válida.
- **Configuración interactiva**: Usa `claude-init config` para configurar cualquier proveedor.

## Ejemplos

### Proyecto Go

```bash
cd /path/to/go-project
claude-init init
```

Estructura generada:

```
.claude/
├── agents/
│   ├── architect.md
│   ├── developer.md
│   ├── tester.md
│   └── reviewer.md
├── skills/
│   ├── go.md
│   └── testing.md
├── commands/
│   ├── build.md
│   ├── test.md
│   └── lint.md
├── project.yaml
└── development_guide.md
```

### Proyecto Node.js/TypeScript

```bash
cd /path/to/nodejs-project
claude-init init
```

### Proyecto Python

```bash
cd /path/to/python-project
claude-init init --config-dir .ai-config
```

## Desarrollo

### Requisitos

- Go 1.25 o superior
- make (opcional, para usar el Makefile)
- Claude Code PRO (si usas el proveedor `cli`) o suscripción a una API de IA (claude-api, openai, zai)

### Estructura del Proyecto

```
ia_start/
├── cmd/                    # Comandos del CLI
│   ├── root/              # Comando raíz
│   ├── init/              # Comando init
│   ├── generate/          # Comando generate
│   ├── config/            # Comando config
│   ├── version/           # Comando version
│   └── completion/        # Comando completion
├── internal/
│   ├── ai/                # Clientes de IA (Claude CLI, Claude API, OpenAI, Z.AI)
│   ├── claude/            # Analizador de proyectos y generador de contenido
│   ├── config/            # Gestión de configuración
│   ├── logger/            # Utilidades de logging
│   └── survey/            # Sistema de preguntas interactivas
├── main.go                # Punto de entrada
├── Makefile              # Automatización de build
└── go.mod                # Módulos de Go (Go 1.25)
```

### Comandos de Make

```bash
# Build
make build              # Build a bin/
make build-local        # Build al directorio actual
make build-all          # Build para todas las plataformas
make clean              # Limpiar artefactos de build
make clean-all          # Limpiar todos los artefactos (build + dist)

# Release
make release            # Crear release (build-all + archives + checksums)
make release-checksums  # Generar checksums para release existente
make verify-checksums   # Verificar todos los checksums

# Testing
make test               # Ejecutar todos los tests
make test-short         # Ejecutar tests cortos
make test-race          # Ejecutar con race detector
make test-cover         # Ejecutar con coverage (genera coverage.html)
make test-integration   # Ejecutar tests de integración
make test-all           # Ejecutar todos los tests

# Linting
make lint               # Ejecutar linters
make lint-fix           # Ejecutar linters con auto-fix
make fmt                # Formatear código
make vet                # Ejecutar go vet

# Benchmarks
make benchmark          # Ejecutar benchmarks
make benchmark-cpu      # CPU profiling
make benchmark-mem      # Memory profiling

# Dependencias
make deps               # Instalar dependencias
make deps-update        # Actualizar dependencias
make deps-verify        # Verificar dependencias

# Instalación
make install            # Instalar a GOBIN
make install-tools      # Instalar herramientas de desarrollo

# Utilidad
make run                # Build y ejecutar
make check              # Ejecutar todos los checks (fmt, vet, lint, test)
make ci                 # Ejecutar checks de CI
make help               # Mostrar comandos disponibles
```

## Proceso de Release

Para crear releases, ver [docs/RELEASE.md](docs/RELEASE.md) para instrucciones detalladas.

Resumen rápido:

```bash
# Release automatizado (recomendado)
./scripts/release.sh v0.1.0

# Release manual
make release VERSION=v0.1.0
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

Cuando haces push de un tag, GitHub Actions automáticamente:

1. Ejecuta todos los tests
2. Build para todas las plataformas (Linux, macOS, Windows)
3. Crea GitHub Release
4. Sube artefactos de release con checksums

### Testing

```bash
# Ejecutar todos los tests
go test ./...

# Ejecutar con coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Ejecutar con race detector
go test -race ./...

# Ejecutar tests de paquetes específicos
go test ./internal/survey/...
```

### Calidad de Código

El proyecto usa múltiples linters y herramientas:

- **golangci-lint**: Linting comprehensivo (19+ linters configurados)
- **go fmt**: Formateo de código
- **go vet**: Análisis estático
- **Test coverage**: 83.6% de coverage global

## Contribuyendo

¡Las contribuciones son bienvenidas! Por favor ver [CONTRIBUTING.md](CONTRIBUTING.md) para guías.

## Arquitectura

Para documentación detallada de arquitectura, ver [ARCHITECTURE.md](ARCHITECTURE.md).

## Changelog

Ver [CHANGELOG.md](CHANGELOG.md) para historial de versiones.

## Licencia

MIT License - ver [LICENSE](LICENSE) para detalles.

## Autor

Daniel Rosello Sánchez

## Reconocimientos

Construido con:

- [Cobra](https://github.com/spf13/cobra) - Framework de CLI
- [Survey](https://github.com/AlecAivazis/survey/v2) - Prompts interactivos
- [Anthropic Claude](https://www.anthropic.com/) - API de IA
- [OpenAI](https://openai.com/) - API de IA
