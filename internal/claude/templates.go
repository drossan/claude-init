package claude

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/drossan/claude-init/internal/survey"
)

// TemplateLoader maneja la carga de templates base desde claude_examples/.
type TemplateLoader struct {
	templatesPath string
}

// NewTemplateLoader crea un nuevo TemplateLoader.
func NewTemplateLoader() *TemplateLoader {
	// Buscar claude_examples en varios directorios posibles
	paths := []string{
		"./claude_examples",
		"../claude_examples",
		"../../claude_examples",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return &TemplateLoader{templatesPath: path}
		}
	}

	// Si no se encuentra, retornar con path vacío (se manejará como error al cargar)
	return &TemplateLoader{templatesPath: ""}
}

// Template contiene el contenido de un template.
type Template struct {
	Name    string
	Type    string // "agent", "command", "skill"
	Content string
}

// LoadTemplate carga un template específico por tipo y nombre.
func (tl *TemplateLoader) LoadTemplate(templateType, name string) (*Template, error) {
	if tl.templatesPath == "" {
		return nil, fmt.Errorf("templates path not found")
	}

	var templatePath string
	switch templateType {
	case "agent":
		templatePath = filepath.Join(tl.templatesPath, "agents", name+".md")
	case "command":
		templatePath = filepath.Join(tl.templatesPath, "commands", name+".md")
	case "skill":
		templatePath = filepath.Join(tl.templatesPath, "skills", name+".md")
	default:
		return nil, fmt.Errorf("unknown template type: %s", templateType)
	}

	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	return &Template{
		Name:    name,
		Type:    templateType,
		Content: string(content),
	}, nil
}

// AdaptTemplate adapta un template al proyecto actual reemplazando placeholders.
func (tl *TemplateLoader) AdaptTemplate(template *Template, answers interface{}) string {
	content := template.Content

	// Convertir answers a la estructura esperada
	ans, ok := answers.(*survey.Answers)
	if !ok {
		return content
	}

	// Reemplazos específicos según el tipo de template
	if template.Type == "agent" {
		content = tl.adaptAgentTemplate(content, ans)
	} else if template.Type == "command" {
		content = tl.adaptCommandTemplate(content, ans)
	} else if template.Type == "skill" {
		content = tl.adaptSkillTemplate(content, ans)
	}

	return content
}

// adaptAgentTemplate adapta un template de agent al proyecto.
func (tl *TemplateLoader) adaptAgentTemplate(content string, ans *survey.Answers) string {
	// Nombre del proyecto base que viene en los templates
	baseProjectName := "Griddo API"

	// Reemplazar el nombre del proyecto
	content = strings.ReplaceAll(content, baseProjectName, ans.ProjectName)
	content = strings.ReplaceAll(content, "Griddo API", ans.ProjectName)
	content = strings.ReplaceAll(content, "Griddo", ans.Framework)

	// Reemplazar tecnologías específicas según el lenguaje/framework
	content = tl.replaceTechnologySpecifics(content, ans)

	// Actualizar descripción en frontmatter si existe
	if strings.Contains(content, "description:") {
		// Buscar y reemplazar descripción que contiene "Griddo API"
		re := regexp.MustCompile(`description: Especialista[^\n]*para Griddo API[^.]*\.`)
		newDesc := fmt.Sprintf("description: Especialista en diseño de arquitectura para %s. Responsable de definir la estructura de módulos, capas y la interacción entre componentes.", ans.ProjectName)
		content = re.ReplaceAllString(content, newDesc)
	}

	return content
}

// adaptCommandTemplate adapta un template de command al proyecto.
func (tl *TemplateLoader) adaptCommandTemplate(content string, ans *survey.Answers) string {
	// Reemplazar el nombre del proyecto
	content = strings.ReplaceAll(content, "Griddo API", ans.ProjectName)
	content = strings.ReplaceAll(content, "Griddo", ans.Framework)

	// Reemplazar tecnologías específicas
	content = tl.replaceTechnologySpecifics(content, ans)

	return content
}

// adaptSkillTemplate adapta un template de skill al proyecto.
func (tl *TemplateLoader) adaptSkillTemplate(content string, ans *survey.Answers) string {
	// Reemplazar el nombre del proyecto
	content = strings.ReplaceAll(content, "Griddo API", ans.ProjectName)
	content = strings.ReplaceAll(content, "Griddo", ans.Framework)

	// Reemplazar tecnologías específicas
	content = tl.replaceTechnologySpecifics(content, ans)

	return content
}

// replaceTechnologySpecifics reemplaza referencias a tecnologías específicas
// según el lenguaje y framework del proyecto.
func (tl *TemplateLoader) replaceTechnologySpecifics(content string, ans *survey.Answers) string {
	switch strings.ToLower(ans.Language) {
	case "go", "golang":
		content = strings.ReplaceAll(content, "TypeScript", "Go")
		content = strings.ReplaceAll(content, "typescript", "go")
		content = strings.ReplaceAll(content, "npm run", "go run")
		content = strings.ReplaceAll(content, "Express", "Gin" /* o el framework que sea */)
		content = strings.ReplaceAll(content, "Zod", "validator" /* o lib de validación Go */)
		content = strings.ReplaceAll(content, "TypeORM", "GORM" /* o ORM Go */)

	case "typescript", "javascript", "js", "ts":
		// Ya está en TypeScript, mantener referencias
		if ans.Framework != "" {
			// Podríamos agregar referencias específicas del framework
			frameworkLower := strings.ToLower(ans.Framework)
			if strings.Contains(frameworkLower, "react") {
				content = strings.ReplaceAll(content, "Express", "React")
			} else if strings.Contains(frameworkLower, "next") {
				content = strings.ReplaceAll(content, "Express", "Next.js")
			} else if strings.Contains(frameworkLower, "vue") {
				content = strings.ReplaceAll(content, "Express", "Vue")
			} else if strings.Contains(frameworkLower, "angular") {
				content = strings.ReplaceAll(content, "Express", "Angular")
			}
		}

	case "python":
		content = strings.ReplaceAll(content, "TypeScript", "Python")
		content = strings.ReplaceAll(content, "typescript", "python")
		content = strings.ReplaceAll(content, "npm run", "python")
		content = strings.ReplaceAll(content, "Express", "Flask" /* o FastAPI */)
		content = strings.ReplaceAll(content, "Zod", "Pydantic" /* o lib de validación Python */)
		content = strings.ReplaceAll(content, "TypeORM", "SQLAlchemy" /* o ORM Python */)
		content = strings.ReplaceAll(content, ".ts", ".py")

	case "rust":
		content = strings.ReplaceAll(content, "TypeScript", "Rust")
		content = strings.ReplaceAll(content, "typescript", "rust")
		content = strings.ReplaceAll(content, "npm run", "cargo run")
		content = strings.ReplaceAll(content, "Express", "Actix-web" /* o framework Rust */)
		content = strings.ReplaceAll(content, ".ts", ".rs")
	}

	return content
}

// HasTemplate verifica si existe un template para el tipo y nombre dados.
func (tl *TemplateLoader) HasTemplate(templateType, name string) bool {
	if tl.templatesPath == "" {
		return false
	}

	var templatePath string
	switch templateType {
	case "agent":
		templatePath = filepath.Join(tl.templatesPath, "agents", name+".md")
	case "command":
		templatePath = filepath.Join(tl.templatesPath, "commands", name+".md")
	case "skill":
		templatePath = filepath.Join(tl.templatesPath, "skills", name+".md")
	default:
		return false
	}

	if _, err := os.Stat(templatePath); err != nil {
		return false
	}

	return true
}

// GetTemplatePath retorna el path donde se encuentran los templates.
func (tl *TemplateLoader) GetTemplatePath() string {
	return tl.templatesPath
}

// NormalizeSkillName normaliza un nombre de skill para buscar el template.
// Por ejemplo, "typescript" → "typescript-expert", "code-reviewer" → "code-reviewer"
func NormalizeSkillName(skillName string) string {
	// Mapeo de nombres simples a nombres de templates
	mapping := map[string]string{
		"typescript":       "typescript-expert",
		"javascript":       "javascript-expert",
		"go":               "go-expert",
		"python":           "python-expert",
		"code-reviewer":    "code-reviewer",
		"technical-writer": "technical-writer",
		"debug-master":     "debug-master",
	}

	if normalized, ok := mapping[strings.ToLower(skillName)]; ok {
		return normalized
	}

	return skillName
}

// GetAvailableTemplates retorna una lista de todos los templates disponibles.
func (tl *TemplateLoader) GetAvailableTemplates() map[string][]string {
	if tl.templatesPath == "" {
		return map[string][]string{}
	}

	result := map[string][]string{
		"agents":   {},
		"commands": {},
		"skills":   {},
	}

	// Leer agents
	if agentsDir := filepath.Join(tl.templatesPath, "agents"); dirExists(agentsDir) {
		files, _ := filepath.Glob(filepath.Join(agentsDir, "*.md"))
		for _, f := range files {
			name := strings.TrimSuffix(filepath.Base(f), ".md")
			result["agents"] = append(result["agents"], name)
		}
	}

	// Leer commands
	if commandsDir := filepath.Join(tl.templatesPath, "commands"); dirExists(commandsDir) {
		files, _ := filepath.Glob(filepath.Join(commandsDir, "*.md"))
		for _, f := range files {
			name := strings.TrimSuffix(filepath.Base(f), ".md")
			result["commands"] = append(result["commands"], name)
		}
	}

	// Leer skills
	if skillsDir := filepath.Join(tl.templatesPath, "skills"); dirExists(skillsDir) {
		files, _ := filepath.Glob(filepath.Join(skillsDir, "*.md"))
		for _, f := range files {
			name := strings.TrimSuffix(filepath.Base(f), ".md")
			result["skills"] = append(result["skills"], name)
		}
	}

	return result
}

// dirExists verifica si un directorio existe.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
