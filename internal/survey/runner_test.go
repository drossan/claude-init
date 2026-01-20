package survey

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRunner_NewRunner_CreatesRunner verifica que NewRunner crea un runner válido.
func TestRunner_NewRunner_CreatesRunner(t *testing.T) {
	// Arrange
	questions := GetProjectQuestions()

	// Act
	r := NewRunner(questions)

	// Assert
	assert.NotNil(t, r)
	assert.NotNil(t, r.Survey)
	assert.Equal(t, questions, r.Survey.Questions)
}

// TestRunner_Run_WithMockInput_CollectsAnswers verifica que Run recolecta respuestas con input mockeado.
func TestRunner_Run_WithMockInput_CollectsAnswers(t *testing.T) {
	t.Skip("requiere implementación de stdin mocking")
}

// TestRunner_Run_WithEmptyRequiredField_ReturnsError verifica que Run retorna error si falta campo requerido.
func TestRunner_Run_WithEmptyRequiredField_ReturnsError(t *testing.T) {
	t.Skip("requiere implementación de stdin mocking")
}

// TestAnswers_ToMap_ConvertsToMap verifica que ToMap convierte Answers a map[string]string.
func TestAnswers_ToMap_ConvertsToMap(t *testing.T) {
	// Arrange
	answers := &Answers{
		ProjectOrigin:   "new",
		ProjectName:     "test-project",
		Description:     "A test project",
		Language:        "Go",
		Framework:       "Gin",
		Architecture:    "Hexagonal",
		Database:        "PostgreSQL",
		ProjectCategory: "API REST",
		BusinessContext: "E-commerce platform",
	}

	// Act
	result := answers.ToMap()

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "new", result["project_origin"])
	assert.Equal(t, "test-project", result["project_name"])
	assert.Equal(t, "A test project", result["description"])
	assert.Equal(t, "Go", result["language"])
	assert.Equal(t, "Gin", result["framework"])
	assert.Equal(t, "Hexagonal", result["architecture"])
	assert.Equal(t, "PostgreSQL", result["database"])
	assert.Equal(t, "API REST", result["project_category"])
	assert.Equal(t, "E-commerce platform", result["business_context"])
}

// TestAnswers_FromMap_LoadsFromMap verifica que FromMap carga Answers desde map[string]string.
func TestAnswers_FromMap_LoadsFromMap(t *testing.T) {
	// Arrange
	input := map[string]string{
		"project_origin":   "new",
		"project_name":     "test-project",
		"description":      "A test project",
		"language":         "Go",
		"framework":        "Gin",
		"architecture":     "Hexagonal",
		"database":         "PostgreSQL",
		"project_category": "API REST",
		"business_context": "E-commerce platform",
	}

	// Act
	answers := &Answers{}
	answers.FromMap(input)

	// Assert
	assert.Equal(t, "new", answers.ProjectOrigin)
	assert.Equal(t, "test-project", answers.ProjectName)
	assert.Equal(t, "A test project", answers.Description)
	assert.Equal(t, "Go", answers.Language)
	assert.Equal(t, "Gin", answers.Framework)
	assert.Equal(t, "Hexagonal", answers.Architecture)
	assert.Equal(t, "PostgreSQL", answers.Database)
	assert.Equal(t, "API REST", answers.ProjectCategory)
	assert.Equal(t, "E-commerce platform", answers.BusinessContext)
}

// TestAnswers_FromMap_WithEmptyMap_CreatesEmptyAnswers verifica que FromMap con map vacío crea Answers vacíos.
func TestAnswers_FromMap_WithEmptyMap_CreatesEmptyAnswers(t *testing.T) {
	// Arrange
	input := map[string]string{}

	// Act
	answers := &Answers{}
	answers.FromMap(input)

	// Assert
	assert.Empty(t, answers.ProjectOrigin)
	assert.Empty(t, answers.ProjectName)
	assert.Empty(t, answers.Description)
	assert.Empty(t, answers.Language)
	assert.Empty(t, answers.Framework)
	assert.Empty(t, answers.Architecture)
	assert.Empty(t, answers.Database)
	assert.Empty(t, answers.ProjectCategory)
	assert.Empty(t, answers.BusinessContext)
}
