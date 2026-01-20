# claude-init

> CLI para inicializar proyectos con configuración guiada por IA

[![Go Report Card](https://goreportcard.com/badge/github.com/danielrossellosanchez/claude-init)](https://goreportcard.com/report/github.com/danielrossellosanchez/claude-init)
[![Coverage](https://img.shields.io/badge/coverage-83.6%25-brightgreen)](https://github.com/danielrossellosanchez/claude-init)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://go.dev/)

**claude-init** es una herramienta CLI escrita en Go que inicializa proyectos con configuración optimizada para
desarrollo guiado por IA con Claude Code. A través de un survey interactivo, recopila información sobre tu proyecto,
opcionalmente valida con IA, y genera la estructura `.claude/` necesaria.

## Características

- **Survey Interactivo**: 8 preguntas clave sobre tu proyecto
    - Nombre, descripción, lenguaje principal
    - Framework, arquitectura, base de datos
    - Tipo de proyecto y contexto del negocio
- **Validación con IA**: Detecta información faltante automáticamente
- **Recomendaciones Inteligentes**: Sugiere agents, commands y skills basados en tu contexto
- **Templates Basados en claude_examples**: Estructura probada en producción
- **Integración con Múltiples Proveedores de IA**: Claude (Anthropic), OpenAI, z.ai
- **Alta Cobertura de Tests**: Código de alta calidad con tests comprehensivos
- **Personalizable**: Soporte para directorios de configuración custom y modo dry-run

## Instalación

### Desde código fuente

```bash
git clone https://github.com/danielrossellosanchez/claude-init.git
cd claude-init
make build
sudo mv bin/claude-init /usr/local/bin/
```

### Usando go install

```bash
go install github.com/danielrossellosanchez/claude-init@latest
```

### Desde binarios precompilados

Los binarios precompilados están disponibles en la página
de [releases](https://github.com/danielrossellosanchez/claude-init/releases) para:

- Linux (amd64)
- macOS (amd64, arm64)
- Windows (amd64)

## Uso Rápido

```bash
# Inicializar en el directorio actual (el wizard te guiará)
claude-init init

# Inicializar en un path específico
claude-init init /path/to/project

# Configurar proveedor de IA
claude-init config set --provider claude --api-key sk-ant-xxxxx

# Ver configuración actual
claude-init config show

# Mostrar versión
claude-init version

# Habilitar autocompletado (bash)
claude-init completion bash > /etc/bash_completion.d/claude-init
```

## Flujo de Trabajo

El comando `claude-init init` sigue este flujo:

1. **Configuración de IA** (si es necesario): Si no hay credenciales configuradas, el wizard te guía para configurar tu
   proveedor de IA
2. **Survey Interactivo**: Responde 8 preguntas sobre tu proyecto
3. **Validación con IA** (opcional): La IA detecta información faltante
4. **Recomendaciones** (opcional): Sugiere estructura optimizada
5. **Generación**: Crea la estructura `.claude/` desde `claude_examples/`

### Preguntas del Survey

El CLI te hará las siguientes preguntas:

1. **Nombre del proyecto**: Identificador único
2. **Descripción breve**: Resumen del propósito
3. **Lenguaje principal**: Go, Node.js, Python, Rust, etc.
4. **Framework** (opcional): Express, NestJS, Django, etc.
5. **Arquitectura deseada**: Monolito, Microservicios, Hexagonal, etc.
6. **Base de datos** (opcional): PostgreSQL, MongoDB, etc.
7. **Tipo de proyecto**: API REST, Web App, CLI, etc.
8. **Contexto del negocio**: Descripción detallada del dominio

## Comandos

### init

Inicializa el directorio `.claude/` con agents, skills, commands y guías de desarrollo.

```bash
claude-init init [path] [flags]
```

**Flags:**

- `-f, --force`: Sobrescribe archivos existentes
- `--ai-provider`: Proveedor de IA (claude, openai, zai)
- `--api-key`: API key del proveedor de IA
- `--no-ai`: No usar IA, usar valores por defecto
- `--dry-run`: Muestra qué se generaría sin crear archivos
- `--config-dir`: Directorio de configuración (default: `.claude`)

**Ejemplos:**

```bash
# Inicializar en el directorio actual
claude-init init

# Inicializar en un path específico
claude-init init /path/to/project

# Usar recomendaciones de IA
claude-init init --ai-provider claude --api-key sk-ant-xxxxx

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

Gestiona la configuración de proveedores de IA.

```bash
claude-init config [command]
```

**Subcomandos:**

- `show`: Muestra la configuración actual de IA
- `set`: Configura un proveedor de IA
- `list`: Lista los proveedores disponibles
- `unset`: Elimina la configuración de un proveedor

**Ejemplos:**

```bash
# Mostrar configuración actual
claude-init config show

# Configurar Anthropic Claude
claude-init config set --provider claude --api-key sk-ant-xxxxx

# Configurar OpenAI
claude-init config set --provider openai --api-key sk-xxxxx

# Configurar z.ai
claude-init config set --provider zai --api-key zai-xxxxx

# Listar proveedores disponibles
claude-init config list

# Eliminar configuración de un proveedor
claude-init config unset --provider claude

# Usar un path de configuración custom
claude-init config show --config-path /custom/path/config.yaml
```

**Wizard de Configuración:**

La primera vez que ejecutes `claude-init init` sin tener credenciales configuradas, el wizard te guiará para seleccionar
un proveedor e ingresar tu API key de forma interactiva.

```
============================================================
AI Configuration Required
============================================================

No AI configuration found. Let's set up your AI provider.
You can always change this later with: claude-init config set

? Would you like to configure an AI provider now? Yes
? Select an AI provider: Anthropic Claude (Recommended)
? Enter your Anthropic Claude API key: sk-ant-xxxxx

✓ AI configuration saved successfully!
  Provider: Anthropic Claude
  Config:   ~/.config/claude-init/config.yaml
```

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
provider: claude

# Configuración de proveedores
providers:
  claude:
    api_key: sk-ant-xxxxx
    base_url: https://api.anthropic.com
    model: claude-3-5-sonnet-20241022
    max_tokens: 8192
  openai:
    api_key: sk-xxxxx
    base_url: https://api.openai.com/v1
    model: gpt-4o
    max_tokens: 4096
  zai:
    api_key: zai-xxxxx
    base_url: https://api.z.ai
    model: zai-1-xxxxx
    max_tokens: 4096
```

### Gestión de Configuración

Usa el comando `config` para gestionar tus credenciales:

```bash
# Ver configuración actual
claude-init config show

# Configurar un proveedor
claude-init config set --provider claude --api-key sk-ant-xxxxx

# Listar proveedores disponibles
claude-init config list

# Eliminar configuración
claude-init config unset --provider claude
```

### Variables de Entorno (Alternativa)

También puedes usar variables de entorno (tienen menor prioridad que los flags):

```bash
export CLAUDE_INIT_AI_PROVIDER=claude
export CLAUDE_INIT_CLAUDE_API_KEY=sk-ant-xxxxx
export CLAUDE_INIT_OPENAI_API_KEY=sk-xxxxx
export CLAUDE_INIT_ZAI_API_KEY=zai-xxxxx
```

## Ejemplos

### Proyecto Go

```bash
cd /path/to/go-project
claude-init init

# Con recomendaciones de IA
claude-init init --ai-provider claude --api-key sk-ant-xxxxx
```

Estructura generada:

```
.claude/
├── agents/
│   ├── architect.md
│   ├── developer.md
│   ├── tester.md
│   └── writer.md
├── skills/
│   ├── language_go.md
│   └── skill_testing.md
├── commands/
│   ├── build.md
│   ├── test.md
│   └── run.md
├── development_guide.md
└── .gitignore
```

### Proyecto Node.js/TypeScript

```bash
cd /path/to/nodejs-project
claude-init init

# Con recomendaciones de IA
claude-init init --ai-provider claude --api-key sk-ant-xxxxx
```

### Proyecto Python

```bash
cd /path/to/python-project
claude-init init

# Con directorio de configuración custom
claude-init init --config-dir .ai-config
```

## Desarrollo

### Requisitos

- Go 1.25 o superior
- make (opcional, para usar el Makefile)

### Estructura del Proyecto

```
ia_start/
├── cmd/                    # Comandos del CLI
│   ├── root/              # Comando raíz
│   ├── init/              # Comando init
│   ├── generate/          # Comando generate
│   ├── version/           # Comando version
│   └── completion/        # Comando completion
├── internal/
│   ├── ai/                # Clientes de IA (Claude, OpenAI, z.ai)
│   ├── claudeexamples/    # Carga y generación desde claude_examples/
│   ├── config/            # Gestión de configuración
│   ├── logger/            # Utilidades de logging
│   ├── survey/            # Sistema de preguntas interactivas
│   └── templates/         # Generación de templates
│       ├── embed/         # Templates embebidos
│       ├── agents/        # Templates de agentes
│       ├── skills/        # Templates de skills
│       └── commands/      # Templates de comandos
├── tests/                 # Tests de integración
├── main.go                # Punto de entrada
├── Makefile              # Automatización de build
└── go.mod                # Módulos de Go
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
