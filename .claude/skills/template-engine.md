# Skill: Template Engine (template-engine)

## Propósito
Especialidad en generar plantillas dinámicas para la estructura `.claude/` basadas en el contexto del proyecto detectado.

## Responsabilidades
- Generar archivos `.md` con formato correcto (frontmatter con YAML)
- Crear plantillas que se adapten al stack detectado
- Mantener consistencia en la estructura de archivos generados
- Evitar duplicación de código en las plantillas

## Estructura de una Plantilla

### Frontmatter (YAML)

Todos los archivos de agents, skills y commands deben tener frontmatter:

```go
const agentTemplate = `---
name: {{.Name}}
description: {{.Description}}
tools: Read, Write, Edit, Bash, Glob
model: {{.Model}}
color: {{.Color}}
---

# {{.Title}}

## Rol
{{.Role}}

## Tu Especialidad
{{.Specialty}}
`
```

### Generador con Go Templates

```go
package templates

import (
    "bytes"
    "embed"
    "text/template"
)

//go:embed templates/*
var templateFS embed.FS

type Generator struct {
    templates map[string]*template.Template
}

func NewGenerator() (*Generator, error) {
    g := &Generator{
        templates: make(map[string]*template.Template),
    }

    // Cargar plantillas
    files, _ := templateFS.ReadDir("templates")
    for _, f := range files {
        content, _ := templateFS.ReadFile("templates/" + f.Name())
        tmpl, err := template.New(f.Name()).Parse(string(content))
        if err != nil {
            return nil, err
        }
        g.templates[f.Name()] = tmpl
    }

    return g, nil
}

type AgentData struct {
    Name        string
    Description string
    Tools       []string
    Model       string
    Color       string
    Title       string
    Role        string
    Specialty   string
}

func (g *Generator) GenerateAgent(data AgentData) (string, error) {
    tmpl, ok := g.templates["agent.md"]
    if !ok {
        return "", errors.New("agent template not found")
    }

    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }

    return buf.String(), nil
}
```

## Plantillas por Stack

### Skills por Framework

```go
type StackConfig struct {
    Frontend struct {
        Framework string // react, vue, angular, svelte
        Language  string // typescript, javascript
        State     string // redux, zustand, pinia
        CSS       string // css-modules, tailwind, scss
    }
    Backend struct {
        Language   string // go, node, python
        Framework  string // express, fastapi, gin
        Database   string // postgres, mysql, sqlite
        ORM        string // prisma, typeorm, gorm
    }
    Testing struct {
        Framework  string // vitest, jest, pytest
        E2E        string // playwright, cypress
        Coverage   bool
    }
}

func (g *Generator) GenerateSkills(config StackConfig) (map[string]string, error) {
    skills := make(map[string]string)

    // Frontend skills
    if config.Frontend.Framework != "" {
        skills[config.Frontend.Framework+".md"] = g.generateFrameworkSkill(config.Frontend.Framework)
    }

    // Backend skills
    if config.Backend.Language != "" {
        skills[config.Backend.Language+".md"] = g.generateLanguageSkill(config.Backend.Language)
    }

    // Testing skills
    if config.Testing.Framework != "" {
        skills[config.Testing.Framework+".md"] = g.generateTestingSkill(config.Testing.Framework)
    }

    return skills, nil
}
```

### Plantilla de Skill Dinámica

```go
const skillTemplate = `# Skill: {{.Name}} ({{.Slug}})

## Propósito
{{.Purpose}}

## Responsabilidades
{{range .Responsibilities}}
- {{.}}
{{end}}

{{if .Conventions}}
## Convenciones

{{range .Conventions}}
### {{.Title}}

{{.Content}}
{{end}}
{{end}}

{{if .Examples}}
## Ejemplos

{{range .Examples}}
\`\`\`{{.Lang}}
{{.Code}}
\`\`\`
{{end}}
{{end}}

{{if .AntiPatterns}}
## Anti-Patterns a Evitar

{{range .AntiPatterns}}
### {{.Title}}

\`\`\`{{.Lang}}
{{.Bad}}
\`\`\`

✅ **Bien**:
\`\`\`{{.Lang}}
{{.Good}}
\`\`\`
{{end}}
{{end}}
`

type SkillData struct {
    Name          string
    Slug          string
    Purpose       string
    Responsibilities []string
    Conventions   []Convention
    Examples      []Example
    AntiPatterns  []AntiPattern
}

type Convention struct {
    Title   string
    Content string
}

type Example struct {
    Lang string
    Code string
}

type AntiPattern struct {
    Title string
    Lang  string
    Bad   string
    Good  string
}
```

## Generación de Development Guide

```go
func (g *Generator) GenerateDevelopmentGuide(config StackConfig) (string, error) {
    sections := []string{
        g.generateTitleSection(config),
        g.generateArchitectureSection(config),
        g.generateTestingSection(config),
        g.generateCommandsSection(config),
    }

    return strings.Join(sections, "\n\n"), nil
}

func (g *Generator) generateTitleSection(config StackConfig) string {
    return fmt.Sprintf(`# Guía de Desarrollo - %s

Este documento establece los estándares y convenciones para el proyecto.

## Stack Tecnológico

**Frontend**: %s %s
**Backend**: %s %s
**Testing**: %s
`,
        config.ProjectName,
        config.Frontend.Language,
        config.Frontend.Framework,
        config.Backend.Language,
        config.Backend.Framework,
        config.Testing.Framework,
    )
}
```

## Detector de Proyectos para Contexto

```go
type ProjectContext struct {
    Name         string
    Type         string // frontend, backend, fullstack, monorepo
    Language     string
    Framework    string
    Dependencies []string
    Files        []string
}

func DetectProjectContext(path string) (ProjectContext, error) {
    ctx := ProjectContext{
        Name: filepath.Base(path),
    }

    // Detectar tipo
    if fileExists(path, "package.json") {
        return detectNodeContext(path)
    }
    if fileExists(path, "go.mod") {
        return detectGoContext(path)
    }
    if fileExists(path, "requirements.txt") {
        return detectPythonContext(path)
    }

    return ctx, nil
}

func detectNodeContext(path string) (ProjectContext, error) {
    // Leer package.json
    content, _ := os.ReadFile(filepath.Join(path, "package.json"))

    var pkg struct {
        Name          string            `json:"name"`
        Dependencies  map[string]string `json:"dependencies"`
        DevDependencies map[string]string `json:"devDependencies"`
    }
    json.Unmarshal(content, &pkg)

    ctx := ProjectContext{
        Name:     pkg.Name,
        Type:     "backend",
        Language: "typescript",
    }

    // Detectar framework
    for dep := range pkg.Dependencies {
        switch {
        case dep == "react":
            ctx.Type = "frontend"
            ctx.Framework = "react"
        case dep == "express":
            ctx.Type = "backend"
            ctx.Framework = "express"
        case dep == "vue":
            ctx.Type = "frontend"
            ctx.Framework = "vue"
        }
    }

    return ctx, nil
}
```

## Escritura de Archivos

```go
type Writer struct {
    targetDir string
    overwrite bool
}

func NewWriter(targetDir string, overwrite bool) *Writer {
    return &Writer{
        targetDir: targetDir,
        overwrite: overwrite,
    }
}

func (w *Writer) WriteAgents(agents map[string]string) error {
    agentDir := filepath.Join(w.targetDir, ".claude", "agents")
    if err := os.MkdirAll(agentDir, 0755); err != nil {
        return err
    }

    for name, content := range agents {
        path := filepath.Join(agentDir, name)

        // Verificar si existe
        if !w.overwrite && fileExists(path) {
            fmt.Printf("Skipping %s (already exists)\n", name)
            continue
        }

        if err := os.WriteFile(path, []byte(content), 0644); err != nil {
            return fmt.Errorf("failed to write %s: %w", name, err)
        }
        fmt.Printf("Created %s\n", name)
    }

    return nil
}

func (w *Writer) WriteSkills(skills map[string]string) error {
    skillDir := filepath.Join(w.targetDir, ".claude", "skills")
    return w.writeFiles(skillDir, skills)
}

func (w *Writer) WriteCommands(commands map[string]string) error {
    cmdDir := filepath.Join(w.targetDir, ".claude", "commands")
    return w.writeFiles(cmdDir, commands)
}

func (w *Writer) WriteDevelopmentGuide(content string) error {
    path := filepath.Join(w.targetDir, ".claude", "development_guide.md")

    if !w.overwrite && fileExists(path) {
        fmt.Println("Skipping development_guide.md (already exists)")
        return nil
    }

    return os.WriteFile(path, []byte(content), 0644)
}

func (w *Writer) writeFiles(dir string, files map[string]string) error {
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    for name, content := range files {
        path := filepath.Join(dir, name)

        if !w.overwrite && fileExists(path) {
            fmt.Printf("Skipping %s (already exists)\n", name)
            continue
        }

        if err := os.WriteFile(path, []byte(content), 0644); err != nil {
            return fmt.Errorf("failed to write %s: %w", name, err)
        }
        fmt.Printf("Created %s\n", name)
    }

    return nil
}

func fileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}
```

## Plantillas Predefinidas

### Template de Command

```go
const commandTemplate = `---
name: {{.Name}}
description: {{.Description}}
usage: "{{.Usage}}"
---

# Comando: {{.Title}}

Este comando {{.Description}}.

## Ejemplo de Uso

\`\`\`bash
claude-init {{.Name}} [opciones]
\`\`\`

## Opciones

{{range .Flags}}
--{{.Name}} ({{.Short}}): {{.Description}}
{{end}}

## Salida

{{.Output}}
`
```

### Template de Development Guide

```go
const devGuideTemplate = `# Guía de Desarrollo - {{.ProjectName}}

## Stack Tecnológico

{{- if .Frontend}}
- **Frontend**: {{.Frontend.Language}} con {{.Frontend.Framework}}
{{- end}}
{{- if .Backend}}
- **Backend**: {{.Backend.Language}} con {{.Backend.Framework}}
{{- end}}

## Estructura del Proyecto

\`\`\`
{{.ProjectName}}/
{{range .Structure}}
├── {{.}}
{{end}}
\`\`\`

## Comandos Disponibles

{{range .Commands}}
- \`{{.Name}}\`: {{.Description}}
{{end}}

## Testing

\`\`\`bash
{{.TestCommand}}
\`\`\`
`
```

## Checklist de Template Engine

- [ ] Las plantillas tienen frontmatter YAML válido
- [ ] Las plantillas usan Go templates (`text/template`)
- [ ] Las plantillas están incrustadas con `embed`
- [ ] Los nombres de archivo siguen convenciones
- [ ] El contenido se valida antes de escribir
- [ ] Se respetan los archivos existentes (si no es overwrite)
- [ ] Los directorios se crean con permisos correctos
- [ ] Los templates manejan campos opcionales