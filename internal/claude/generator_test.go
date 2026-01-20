package claude

import (
	"testing"
)

// TestSanitizeFilename verifies that sanitizeFilename correctly converts various naming formats to kebab-case.
func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "kebab-case already",
			input:    "code-reviewer",
			expected: "code-reviewer",
		},
		{
			name:     "PascalCase",
			input:    "CodeReviewer",
			expected: "code-reviewer",
		},
		{
			name:     "camelCase",
			input:    "codeReviewer",
			expected: "code-reviewer",
		},
		{
			name:     "spaces",
			input:    "code reviewer",
			expected: "code-reviewer",
		},
		{
			name:     "mixed case and spaces",
			input:    "MyCode Reviewer",
			expected: "my-code-reviewer",
		},
		{
			name:     "with underscores",
			input:    "code_reviewer",
			expected: "code_reviewer",
		},
		{
			name:     "PascalCase with multiple words",
			input:    "PlanningAgent",
			expected: "planning-agent",
		},
		{
			name:     "camelCase with numbers",
			input:    "api2Docs",
			expected: "api2-docs",
		},
		{
			name:     "with special characters",
			input:    "code@reviewer!",
			expected: "code-reviewer",
		},
		{
			name:     "already kebab-case with numbers",
			input:    "test-runner-2",
			expected: "test-runner-2",
		},
		{
			name:     "single word",
			input:    "architect",
			expected: "architect",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "unnamed",
		},
		{
			name:     "only special characters",
			input:    "@#$!",
			expected: "unnamed",
		},
		{
			name:     "multiple consecutive spaces",
			input:    "code   reviewer",
			expected: "code-reviewer",
		},
		{
			name:     "leading/trailing hyphens",
			input:    "-code-reviewer-",
			expected: "code-reviewer",
		},
		{
			name:     "ALL CAPS",
			input:    "CODE_REVIEWER",
			expected: "code_reviewer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeFilename(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
