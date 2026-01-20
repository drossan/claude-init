package survey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewSurvey_CreatesSurvey verifica que NewSurvey crea un survey válido.
func TestNewSurvey_CreatesSurvey(t *testing.T) {
	// Arrange & Act
	s := NewSurvey()

	// Assert
	assert.NotNil(t, s)
	assert.NotNil(t, s.Questions)
	assert.Len(t, s.Questions, 8, "debe tener 8 preguntas predefinidas")
}

// TestSurvey_Run_WithAllQuestions_CollectsAnswers verifica que Run recolecta todas las respuestas.
func TestSurvey_Run_WithAllQuestions_CollectsAnswers(t *testing.T) {
	t.Skip("requiere stdin/stdout mocking - se implementará en GREEN phase")
}

// TestSurvey_Run_WithRequiredQuestion_ReturnsError verifica que Run retorna error si falta pregunta requerida.
func TestSurvey_Run_WithRequiredQuestion_ReturnsError(t *testing.T) {
	t.Skip("requiere stdin/stdout mocking - se implementará en GREEN phase")
}

// TestSurvey_Run_WithSelect_UsesOptions verifica que Run usa las opciones para select.
func TestSurvey_Run_WithSelect_UsesOptions(t *testing.T) {
	t.Skip("requiere stdin/stdout mocking - se implementará en GREEN phase")
}

// TestValidateAnswers_Valid_ReturnsNil verifica que ValidateAnswers retorna nil para respuestas válidas.
func TestValidateAnswers_Valid_ReturnsNil(t *testing.T) {
	// Arrange
	answers := &Answers{
		ProjectName:     "test-project",
		Description:     "A test project",
		Language:        "Go",
		Framework:       "Gin",
		Architecture:    "Hexagonal",
		Database:        "PostgreSQL",
		ProjectCategory: "API REST",
		BusinessContext: "E-commerce platform for selling products",
	}

	// Act
	err := answers.Validate()

	// Assert
	assert.NoError(t, err)
}

// TestValidateAnswers_MissingRequired_ReturnsError verifica que ValidateAnswers retorna error si falta campo requerido.
func TestValidateAnswers_MissingRequired_ReturnsError(t *testing.T) {
	tests := []struct {
		name    string
		answers *Answers
	}{
		{
			name: "missing project name",
			answers: &Answers{
				Description:     "A test project",
				Language:        "Go",
				Architecture:    "Hexagonal",
				ProjectCategory: "API REST",
				BusinessContext: "E-commerce platform",
			},
		},
		{
			name: "missing description",
			answers: &Answers{
				ProjectName:     "test-project",
				Language:        "Go",
				Architecture:    "Hexagonal",
				ProjectCategory: "API REST",
				BusinessContext: "E-commerce platform",
			},
		},
		{
			name: "missing language",
			answers: &Answers{
				ProjectName:     "test-project",
				Description:     "A test project",
				Architecture:    "Hexagonal",
				ProjectCategory: "API REST",
				BusinessContext: "E-commerce platform",
			},
		},
		{
			name: "missing architecture",
			answers: &Answers{
				ProjectName:     "test-project",
				Description:     "A test project",
				Language:        "Go",
				ProjectCategory: "API REST",
				BusinessContext: "E-commerce platform",
			},
		},
		{
			name: "missing project category",
			answers: &Answers{
				ProjectName:     "test-project",
				Description:     "A test project",
				Language:        "Go",
				Architecture:    "Hexagonal",
				BusinessContext: "E-commerce platform",
			},
		},
		{
			name: "missing business context",
			answers: &Answers{
				ProjectName:     "test-project",
				Description:     "A test project",
				Language:        "Go",
				Architecture:    "Hexagonal",
				ProjectCategory: "API REST",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := tt.answers.Validate()

			// Assert
			require.Error(t, err)
		})
	}
}

// TestValidateAnswers_OptionalFields_AllowsEmpty verifica que ValidateAnswers permite campos opcionales vacíos.
func TestValidateAnswers_OptionalFields_AllowsEmpty(t *testing.T) {
	// Arrange
	answers := &Answers{
		ProjectName:     "test-project",
		Description:     "A test project",
		Language:        "Go",
		Framework:       "", // opcional
		Architecture:    "Hexagonal",
		Database:        "", // opcional
		ProjectCategory: "API REST",
		BusinessContext: "E-commerce platform",
	}

	// Act
	err := answers.Validate()

	// Assert
	assert.NoError(t, err, "framework y database son opcionales")
}

// TestGetProjectQuestions_Returns8Questions verifica que GetProjectQuestions retorna 8 preguntas.
func TestGetProjectQuestions_Returns8Questions(t *testing.T) {
	// Act
	questions := GetProjectQuestions()

	// Assert
	assert.Len(t, questions, 8, "debe retornar 8 preguntas")
}

// TestGetProjectQuestions_HasRequiredFields verifica que las preguntas tienen los campos requeridos.
func TestGetProjectQuestions_HasRequiredFields(t *testing.T) {
	// Act
	questions := GetProjectQuestions()

	// Assert
	for i, q := range questions {
		assert.NotEmpty(t, q.ID, "pregunta %d debe tener ID", i)
		assert.NotEmpty(t, q.Text, "pregunta %d debe tener Text", i)
		assert.NotEmpty(t, q.Type, "pregunta %d debe tener Type", i)
	}
}
