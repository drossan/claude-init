// Package survey proporciona un sistema de preguntas interactivas para recopilar
// información del usuario sobre su proyecto.
package survey

import (
	"fmt"
	"strings"
)

// QuestionType representa el tipo de pregunta.
type QuestionType string

const (
	// QuestionTypeInput es una pregunta de texto libre.
	QuestionTypeInput QuestionType = "input"
	// QuestionTypeSelect es una pregunta de selección única.
	QuestionTypeSelect QuestionType = "select"
	// QuestionTypeMultiSelect es una pregunta de selección múltiple.
	QuestionTypeMultiSelect QuestionType = "multiselect"
	// QuestionTypeConfirm es una pregunta de confirmación (sí/no).
	QuestionTypeConfirm QuestionType = "confirm"
	// QuestionTypeMultiline es una pregunta de texto multilinea.
	QuestionTypeMultiline QuestionType = "multiline"
)

// Question representa una pregunta del survey.
type Question struct {
	ID          string       // Identificador único de la pregunta
	Text        string       // Texto de la pregunta
	Type        QuestionType // Tipo de pregunta
	Required    bool         // Si es true, la respuesta es obligatoria
	Options     []string     // Opciones para select/multiselect
	Default     string       // Valor por defecto
	Placeholder string       // Placeholder para input
}

// Answers contiene las respuestas del usuario.
type Answers struct {
	ProjectOrigin     string   // Origen del proyecto: "new" o "existing"
	ProjectName       string   // Nombre del proyecto
	Description       string   // Descripción breve
	Language          string   // Lenguaje principal
	Framework         string   // Framework (si aplica)
	Architecture      string   // Arquitectura
	Database          string   // Base de datos (si aplica)
	ProjectCategory   string   // Categoría del proyecto (API REST, Web App, CLI, Library, etc.)
	BusinessContext   string   // Contexto del negocio
	AIProvider        string   // Provider de IA: "cli", "claude-api", "openai", "zai"
	DocumentationDirs []string // Directorios de documentación adicionales (para proyectos existentes)
}

// Survey representa un conjunto de preguntas.
type Survey struct {
	Questions []*Question
}

// NewSurvey crea un nuevo survey con las preguntas predefinidas.
func NewSurvey() *Survey {
	return &Survey{
		Questions: GetProjectQuestions(),
	}
}

// GetProjectQuestions retorna las 8 preguntas predefinidas para el proyecto.
func GetProjectQuestions() []*Question {
	return []*Question{
		{
			ID:       "project_name",
			Text:     "Nombre del proyecto:",
			Type:     QuestionTypeInput,
			Required: true,
		},
		{
			ID:       "description",
			Text:     "Descripción breve del proyecto:",
			Type:     QuestionTypeInput,
			Required: true,
		},
		{
			ID:       "language",
			Text:     "Lenguaje principal:",
			Type:     QuestionTypeInput,
			Default:  "",
			Required: true,
		},
		{
			ID:       "framework",
			Text:     "Framework (opcional, presiona Enter para omitir):",
			Type:     QuestionTypeInput,
			Default:  "",
			Required: false,
		},
		{
			ID:       "architecture",
			Text:     "Arquitectura deseada:",
			Type:     QuestionTypeInput,
			Default:  "",
			Required: true,
		},
		{
			ID:       "database",
			Text:     "Base de datos (opcional, presiona Enter para omitir):",
			Type:     QuestionTypeInput,
			Default:  "",
			Required: false,
		},
		{
			ID:       "project_category",
			Text:     "Categoría del proyecto (ej: API REST, Web App, CLI, Library):",
			Type:     QuestionTypeInput,
			Default:  "",
			Required: true,
		},
		{
			ID:          "business_context",
			Text:        "Contexto del negocio (descripción detallada):",
			Type:        QuestionTypeMultiline,
			Required:    true,
			Placeholder: "Escribe una descripción detallada del contexto del negocio, objetivos y requisitos...",
		},
	}
}

// Validate verifica que todas las respuestas requeridas estén presentes.
func (a *Answers) Validate() error {
	// Limpiar espacios en blanco
	projectName := strings.TrimSpace(a.ProjectName)
	description := strings.TrimSpace(a.Description)
	language := strings.TrimSpace(a.Language)
	architecture := strings.TrimSpace(a.Architecture)
	projectCategory := strings.TrimSpace(a.ProjectCategory)
	businessContext := strings.TrimSpace(a.BusinessContext)

	// Validar campos requeridos
	if projectName == "" {
		return fmt.Errorf("project name is required")
	}

	if description == "" {
		return fmt.Errorf("description is required")
	}

	if language == "" {
		return fmt.Errorf("language is required")
	}

	if architecture == "" {
		return fmt.Errorf("architecture is required")
	}

	if projectCategory == "" {
		return fmt.Errorf("project category is required")
	}

	if businessContext == "" {
		return fmt.Errorf("business context is required")
	}

	return nil
}

// ToMap convierte Answers a un map[string]string.
func (a *Answers) ToMap() map[string]string {
	return map[string]string{
		"project_origin":   a.ProjectOrigin,
		"project_name":     a.ProjectName,
		"description":      a.Description,
		"language":         a.Language,
		"framework":        a.Framework,
		"architecture":     a.Architecture,
		"database":         a.Database,
		"project_category": a.ProjectCategory,
		"business_context": a.BusinessContext,
	}
}

// FromMap carga Answers desde un map[string]string.
func (a *Answers) FromMap(m map[string]string) {
	a.ProjectOrigin = m["project_origin"]
	a.ProjectName = m["project_name"]
	a.Description = m["description"]
	a.Language = m["language"]
	a.Framework = m["framework"]
	a.Architecture = m["architecture"]
	a.Database = m["database"]
	a.ProjectCategory = m["project_category"]
	a.BusinessContext = m["business_context"]
}
