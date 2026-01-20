package claude

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAnalyzer_extractJSON_BasicJSON tests extracting basic JSON from a string.
func TestAnalyzer_extractJSON_BasicJSON(t *testing.T) {
	a := &Analyzer{}
	input := `Some text before {"key":"value"} and after`

	result := a.extractJSON(input)

	assert.Equal(t, `{"key":"value"}`, result)
}

// TestAnalyzer_extractJSON_OnlyJSON tests extracting JSON when it's the only content.
func TestAnalyzer_extractJSON_OnlyJSON(t *testing.T) {
	a := &Analyzer{}
	input := `{"name":"test","language":"Go"}`

	result := a.extractJSON(input)

	assert.Equal(t, `{"name":"test","language":"Go"}`, result)
}

// TestAnalyzer_extractJSON_NestedJSON tests extracting nested JSON.
func TestAnalyzer_extractJSON_NestedJSON(t *testing.T) {
	a := &Analyzer{}
	input := `{"outer":{"inner":"value"}}`

	result := a.extractJSON(input)

	assert.Equal(t, `{"outer":{"inner":"value"}}`, result)
}

// TestAnalyzer_extractJSON_NoJSON tests handling strings without JSON.
func TestAnalyzer_extractJSON_NoJSON(t *testing.T) {
	a := &Analyzer{}
	input := `This is just plain text without any JSON`

	result := a.extractJSON(input)

	assert.Equal(t, "", result)
}

// TestAnalyzer_extractJSON_InvalidJSON tests handling strings with invalid JSON structure.
func TestAnalyzer_extractJSON_InvalidJSON(t *testing.T) {
	a := &Analyzer{}
	input := `{"key":"value"`

	result := a.extractJSON(input)

	assert.Equal(t, "", result)
}

// TestAnalyzer_extractJSON_Whitespace tests handling JSON with surrounding whitespace.
func TestAnalyzer_extractJSON_Whitespace(t *testing.T) {
	a := &Analyzer{}
	input := `
	Some text
	{"key":"value"}
	More text
	`

	result := a.extractJSON(input)

	assert.Equal(t, `{"key":"value"}`, result)
}

// TestAnalyzer_parseAnalysis_ValidJSON tests parsing valid JSON analysis.
func TestAnalyzer_parseAnalysis_ValidJSON(t *testing.T) {
	a := &Analyzer{}
	input := `{
		"name": "test-project",
		"description": "A test project",
		"language": "Go",
		"framework": "Gin",
		"architecture": "Hexagonal",
		"database": "PostgreSQL",
		"project_category": "API REST",
		"business_context": "E-commerce platform",
		"git_system": "git",
		"testing_framework": "testify"
	}`

	result, err := a.parseAnalysis(input)

	require.NoError(t, err)
	assert.Equal(t, "test-project", result.Name)
	assert.Equal(t, "A test project", result.Description)
	assert.Equal(t, "Go", result.Language)
	assert.Equal(t, "Gin", result.Framework)
	assert.Equal(t, "Hexagonal", result.Architecture)
	assert.Equal(t, "PostgreSQL", result.Database)
	assert.Equal(t, "API REST", result.ProjectCategory)
	assert.Equal(t, "E-commerce platform", result.BusinessContext)
	assert.Equal(t, "git", result.GitSystem)
	assert.Equal(t, "testify", result.TestingFramework)
}

// TestAnalyzer_parseAnalysis_ValidJSON_OptionalFields tests parsing with optional fields empty.
func TestAnalyzer_parseAnalysis_ValidJSON_OptionalFields(t *testing.T) {
	a := &Analyzer{}
	input := `{
		"name": "simple-project",
		"description": "Simple CLI tool",
		"language": "Go",
		"framework": "",
		"architecture": "Monolito",
		"database": "",
		"project_category": "CLI",
		"business_context": "Command line tool for data processing"
	}`

	result, err := a.parseAnalysis(input)

	require.NoError(t, err)
	assert.Equal(t, "simple-project", result.Name)
	assert.Equal(t, "Go", result.Language)
	assert.Equal(t, "", result.Framework)
	assert.Equal(t, "", result.Database)
	assert.Equal(t, "CLI", result.ProjectCategory)
}

// TestAnalyzer_parseAnalysis_MissingRequiredField tests error when required field is missing.
func TestAnalyzer_parseAnalysis_MissingRequiredField(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError string
	}{
		{
			name: "missing name",
			input: `{
				"description": "Test",
				"language": "Go",
				"architecture": "Clean",
				"project_category": "API",
				"business_context": "Test"
			}`,
			wantError: "missing required field: name",
		},
		{
			name: "missing language",
			input: `{
				"name": "test",
				"description": "Test",
				"architecture": "Clean",
				"project_category": "API",
				"business_context": "Test"
			}`,
			wantError: "missing required field: language",
		},
		{
			name: "missing architecture",
			input: `{
				"name": "test",
				"description": "Test",
				"language": "Go",
				"project_category": "API",
				"business_context": "Test"
			}`,
			wantError: "missing required field: architecture",
		},
		{
			name: "missing project_category",
			input: `{
				"name": "test",
				"description": "Test",
				"language": "Go",
				"architecture": "Clean",
				"business_context": "Test"
			}`,
			wantError: "missing required field: project_category",
		},
		{
			name: "missing business_context",
			input: `{
				"name": "test",
				"description": "Test",
				"language": "Go",
				"architecture": "Clean",
				"project_category": "API"
			}`,
			wantError: "missing required field: business_context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Analyzer{}
			_, err := a.parseAnalysis(tt.input)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantError)
		})
	}
}

// TestAnalyzer_parseAnalysis_InvalidJSON tests error when JSON is invalid.
func TestAnalyzer_parseAnalysis_InvalidJSON(t *testing.T) {
	a := &Analyzer{}
	input := `This is not valid JSON at all`

	_, err := a.parseAnalysis(input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no JSON found in response")
}

// TestAnalyzer_parseAnalysis_MalformedJSON tests error when JSON is malformed.
func TestAnalyzer_parseAnalysis_MalformedJSON(t *testing.T) {
	a := &Analyzer{}
	input := `{"name":"test","invalid": }`

	_, err := a.parseAnalysis(input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}

// TestNewAnalyzer_CreatesAnalyzer tests NewAnalyzer creates a valid Analyzer.
func TestNewAnalyzer_CreatesAnalyzer(t *testing.T) {
	projectPath := "/path/to/project"
	a := NewAnalyzer(projectPath)

	assert.NotNil(t, a)
	assert.Equal(t, projectPath, a.projectPath)
}

// TestAnalyzer_SetLogger tests SetLogger sets the logger.
func TestAnalyzer_SetLogger(t *testing.T) {
	a := NewAnalyzer("/test")
	mockLogger := &mockLogger{}

	a.SetLogger(mockLogger)

	assert.Equal(t, mockLogger, a.logger)
}

// TestAnalyzer_buildAnalysisPrompt_ReturnsPrompt tests buildAnalysisPrompt returns valid prompt.
func TestAnalyzer_buildAnalysisPrompt_ReturnsPrompt(t *testing.T) {
	a := NewAnalyzer("/test/project")

	prompt := a.buildAnalysisPrompt()

	assert.Contains(t, prompt, "/test/project")
	assert.Contains(t, prompt, "JSON")
	assert.Contains(t, prompt, "name")
	assert.Contains(t, prompt, "language")
	assert.Contains(t, prompt, "architecture")
	assert.Contains(t, prompt, "project_category")
	assert.Contains(t, prompt, "business_context")
}

// mockLogger is a simple mock implementation for testing.
type mockLogger struct {
	debugMessages []string
	infoMessages  []string
	warnMessages  []string
	errorMessages []string
}

func (m *mockLogger) Debug(format string, args ...interface{}) {
	m.debugMessages = append(m.debugMessages, format)
}

func (m *mockLogger) Info(format string, args ...interface{}) {
	m.infoMessages = append(m.infoMessages, format)
}

func (m *mockLogger) Warn(format string, args ...interface{}) {
	m.warnMessages = append(m.warnMessages, format)
}

func (m *mockLogger) Error(format string, args ...interface{}) {
	m.errorMessages = append(m.errorMessages, format)
}
