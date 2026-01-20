package claude

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/drossan/claude-init/internal/survey"
)

// PromptBuilder construye prompts para Claude CLI basados en el contexto del proyecto.
type PromptBuilder struct {
	answers *survey.Answers
}

// NewPromptBuilder crea un nuevo PromptBuilder.
func NewPromptBuilder(answers *survey.Answers) *PromptBuilder {
	return &PromptBuilder{
		answers: answers,
	}
}

// buildAgentPrompt construye el prompt específico para cada tipo de agente.
func (pb *PromptBuilder) buildAgentPrompt(agentType string) string {
	prompts := map[string]string{
		"architect": `You are the architect agent for {{.ProjectName}}.

{{if .Description}}Project Description: {{.Description}}{{end}}
{{if .Language}}Language: {{.Language}}{{end}}
{{if .Framework}}Framework: {{.Framework}}{{end}}
{{if .Architecture}}Architecture: {{.Architecture}}{{end}}
{{if .Database}}Database: {{.Database}}{{end}}
{{if .BusinessContext}}Business Context: {{.BusinessContext}}{{end}}

Your responsibilities:
- Design system architecture following best practices
- Define package structure and interfaces
- Ensure scalability and maintainability
- Follow SOLID principles and {{.Language}} conventions

Create a comprehensive agent configuration file that includes:
1. Agent description and role
2. System prompt with architectural guidelines
3. Tools needed for architecture tasks
4. Model preferences for architectural work`,

		"developer": `You are the developer agent for {{.ProjectName}}.

{{if .Description}}Project Description: {{.Description}}{{end}}
{{if .Language}}Language: {{.Language}}{{end}}
{{if .Framework}}Framework: {{.Framework}}{{end}}
{{if .Architecture}}Architecture: {{.Architecture}}{{end}}
{{if .Database}}Database: {{.Database}}{{end}}
{{if .BusinessContext}}Business Context: {{.BusinessContext}}{{end}}

Your responsibilities:
- Write clean, maintainable code following {{.Language}} best practices
- Follow project conventions and architecture
- Write comprehensive tests for new features
- Document complex logic and decisions
- Ensure code quality and performance

Create a comprehensive agent configuration file that includes:
1. Agent description and role
2. System prompt with coding guidelines for {{.Language}}
3. Tools needed for development tasks
4. Model preferences for development work`,

		"tester": `You are the QA/testing agent for {{.ProjectName}}.

{{if .Description}}Project Description: {{.Description}}{{end}}
{{if .Language}}Language: {{.Language}}{{end}}
{{if .Framework}}Framework: {{.Framework}}{{end}}
{{if .BusinessContext}}Business Context: {{.BusinessContext}}{{end}}

Your responsibilities:
- Write comprehensive tests (unit, integration, e2e)
- Find edge cases and potential bugs
- Ensure test coverage meets quality standards
- Validate requirements against test results
- Create testing strategies and plans

Create a comprehensive agent configuration file that includes:
1. Agent description and role
2. System prompt with testing guidelines for {{.Language}}
3. Tools needed for testing tasks
4. Model preferences for testing work`,

		"reviewer": `You are the code review agent for {{.ProjectName}}.

{{if .Description}}Project Description: {{.Description}}{{end}}
{{if .Language}}Language: {{.Language}}{{end}}
{{if .Architecture}}Architecture: {{.Architecture}}{{end}}

Your responsibilities:
- Review code for quality, maintainability, and best practices
- Ensure adherence to {{.Language}} coding standards
- Identify potential bugs, security issues, and performance problems
- Provide constructive feedback and suggestions
- Validate architectural decisions

Create a comprehensive agent configuration file that includes:
1. Agent description and role
2. System prompt with review guidelines for {{.Language}}
3. Tools needed for code review tasks
4. Model preferences for review work`,

		"debugger": `You are the debugging specialist agent for {{.ProjectName}}.

{{if .Description}}Project Description: {{.Description}}{{end}}
{{if .Language}}Language: {{.Language}}{{end}}
{{if .Framework}}Framework: {{.Framework}}{{end}}
{{if .Architecture}}Architecture: {{.Architecture}}{{end}}

Your responsibilities:
- Investigate and resolve complex bugs
- Analyze stack traces and error messages
- Identify race conditions and memory leaks
- Root cause analysis for production issues
- Propose and implement fixes

Create a comprehensive agent configuration file that includes:
1. Agent description and role
2. System prompt with debugging strategies for {{.Language}}
3. Tools needed for debugging tasks
4. Model preferences for debugging work`,

		"planner": `You are the planning agent for {{.ProjectName}}.

{{if .Description}}Project Description: {{.Description}}{{end}}
{{if .Language}}Language: {{.Language}}{{end}}
{{if .BusinessContext}}Business Context: {{.BusinessContext}}{{end}}

Your responsibilities:
- Break down complex tasks into manageable steps
- Estimate effort and identify dependencies
- Create implementation plans
- Identify potential risks and mitigation strategies
- Coordinate with other agents

Create a comprehensive agent configuration file that includes:
1. Agent description and role
2. System prompt with planning guidelines
3. Tools needed for planning tasks
4. Model preferences for planning work`,
	}

	tmpl, ok := prompts[agentType]
	if !ok {
		tmpl = prompts["developer"] // fallback
	}

	return pb.renderTemplate(tmpl)
}

// buildSkillPrompt construye el prompt para generar skills de lenguaje o framework.
func (pb *PromptBuilder) buildSkillPrompt(skillType, skillName string) string {
	var basePrompt string

	switch skillType {
	case "language":
		basePrompt = `Create a language-specific skill file for {{.Language}} in the context of {{.ProjectName}}.

{{if .Description}}Project Description: {{.Description}}{{end}}
{{if .Architecture}}Architecture: {{.Architecture}}{{end}}
{{if .BusinessContext}}Business Context: {{.BusinessContext}}{{end}}

The skill file should include:
1. Language-specific best practices and conventions
2. Common patterns and idioms for {{.Language}}
3. Code style guidelines
4. Testing approaches for {{.Language}}
5. Common pitfalls and how to avoid them
6. Performance considerations specific to {{.Language}}

Format the output as a markdown skill file with proper frontmatter.`

	case "framework":
		basePrompt = `Create a framework-specific skill file for {{.Framework}} in the context of {{.ProjectName}}.

{{if .Description}}Project Description: {{.Description}}{{end}}
{{if .Language}}Language: {{.Language}}{{end}}
{{if .Architecture}}Architecture: {{.Architecture}}{{end}}

The skill file should include:
1. Framework-specific best practices and conventions
2. Common patterns and idioms for {{.Framework}}
3. Project structure conventions
4. Configuration management for {{.Framework}}
5. Testing approaches for {{.Framework}}
6. Performance optimization techniques

Format the output as a markdown skill file with proper frontmatter.`

	default:
		basePrompt = fmt.Sprintf("Create a skill file for %s: %s", skillType, skillName)
	}

	return pb.renderTemplate(basePrompt)
}

// buildCommandPrompt construye el prompt para generar comandos personalizados.
func (pb *PromptBuilder) buildCommandPrompt(commandType string) string {
	prompts := map[string]string{
		"test": `Create a custom test command for {{.ProjectName}}.

{{if .Language}}Language: {{.Language}}{{end}}
{{if .Framework}}Framework: {{.Framework}}{{end}}

The command should:
1. Run all tests in the project
2. Provide clear output format
3. Support test filtering if applicable
4. Include coverage reporting if available
5. Handle project-specific test setup

Format the output as a markdown command file with proper frontmatter including:
- Command description
- Usage examples
- Expected output format`,

		"lint": `Create a custom lint command for {{.ProjectName}}.

{{if .Language}}Language: {{.Language}}{{end}}

The command should:
1. Run linters appropriate for {{.Language}}
2. Check code style and formatting
3. Identify potential issues
4. Provide actionable feedback
5. Support auto-fixing if available

Format the output as a markdown command file with proper frontmatter including:
- Command description
- Usage examples
- Configuration options`,

		"build": `Create a custom build command for {{.ProjectName}}.

{{if .Language}}Language: {{.Language}}{{end}}
{{if .Framework}}Framework: {{.Framework}}{{end}}
{{if .ProjectCategory}}Project Category: {{.ProjectCategory}}{{end}}

The command should:
1. Build the project correctly
2. Handle dependencies
3. Provide clear error messages
4. Support build optimization if applicable
5. Handle environment-specific builds

Format the output as a markdown command file with proper frontmatter including:
- Command description
- Usage examples
- Build options`,

		"docs": `Create a custom documentation command for {{.ProjectName}}.

{{if .Description}}Project Description: {{.Description}}{{end}}

The command should:
1. Generate or serve documentation
2. Support project-specific documentation tools
3. Handle doc generation for {{.Language}}
4. Provide preview capabilities
5. Validate documentation quality

Format the output as a markdown command file with proper frontmatter including:
- Command description
- Usage examples
- Documentation structure`,
	}

	tmpl, ok := prompts[commandType]
	if !ok {
		tmpl = fmt.Sprintf("Create a custom command for %s", commandType)
	}

	return pb.renderTemplate(tmpl)
}

// buildRecommendationPrompt construye el prompt para obtener recomendaciones de estructura.
func (pb *PromptBuilder) buildRecommendationPrompt() string {
	return pb.renderTemplate(`Based on the following project information, recommend the optimal Claude Code structure.

Project Name: {{.ProjectName}}
Description: {{.Description}}
Language: {{.Language}}
Framework: {{.Framework}}
Architecture: {{.Architecture}}
Database: {{.Database}}
Project Category: {{.ProjectCategory}}
Business Context: {{.BusinessContext}}

Please recommend:
1. Which agents should be generated (list agent names)
2. Which commands should be generated (list command names)
3. Which skills should be included (list skill names)
4. A brief description of the recommended structure

CRITICAL NAMING CONVENTION:
ALL agent names, command names, and skill names MUST be in **kebab-case** (lowercase with hyphens).
Examples of valid names:
- ✅ "code-reviewer"
- ✅ "test-runner"
- ✅ "api-docs"
- ✅ "bug-fix"
- ❌ "CodeReviewer" (wrong - PascalCase)
- ❌ "codeReviewer" (wrong - camelCase)
- ❌ "CODE_REVIEWER" (wrong - UPPER_CASE)

Format your response as a JSON object with the following structure:
{
  "agents": ["agent-name-1", "agent-name-2", ...],
  "commands": ["command-name-1", "command-name-2", ...],
  "skills": ["skill-name-1", "skill-name-2", ...],
  "description": "Brief description of the recommended structure"
}

Be specific and base your recommendations on the project's technology stack and requirements.
REMEMBER: ALL names must be in kebab-case (lowercase with hyphens).`)
}

// renderTemplate aplica las respuestas del survey a un template.
func (pb *PromptBuilder) renderTemplate(tmpl string) string {
	t, err := template.New("prompt").Parse(tmpl)
	if err != nil {
		return tmpl // fallback al template original
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, pb.answers)
	if err != nil {
		return tmpl // fallback al template original
	}

	return strings.TrimSpace(buf.String())
}
