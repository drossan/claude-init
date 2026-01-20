# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- **AI template loading bug**: Removed dependency on `clade_skeleton/` filesystem files
  - `LoadTemplate()` functions in prompts packages attempted to read non-embedded files
  - Eliminated `LoadTemplate()` methods from AgentPromptGenerator, SkillPromptGenerator, and CommandPromptGenerator
  - AI generation now uses embedded default templates directly, removing warnings during compilation
  - Removed unused `os` and `filepath` imports from prompt generators
- **Duplicate .claude directory bug**: Fixed path handling in AIGenerator.GenerateAll
  - When `outputPath` already ended with `.claude`, the code was creating `.claude/.claude/` structure
  - Added logic to detect if `outputPath` already ends with `.claude` and use it directly
  - This fixes the issue where files were generated in nested `.claude/.claude/` directories

### Added
- **AI-powered content generation**: Generate `.claude/` files using AI instead of static templates
  - AIGenerator uses AI clients to generate contextualized content for agents, skills, commands, and guides
  - ProjectScanner analyzes project structure to provide additional context to AI
  - Prompt generators create detailed prompts with project-specific information
  - Supports Claude, OpenAI, and z.ai providers for content generation
  - Fallback to local templates when `--no-ai` flag is used
- **`generate` command AI flags**: New flags for AI-powered generation
  - `--ai-provider`: Specify AI provider (claude, openai, zai)
  - `--api-key`: Provide API key for the selected provider
- **Extended AIClient interface**: Added `GenerateContent()` method
  - All AI clients (Claude, OpenAI, z.ai) implement free-text content generation
  - Enables generating markdown content for agents, skills, commands, and guides
- **Enhanced TemplateContext**: Added new fields for better AI generation
  - `Architecture`: Project architecture pattern (hexagonal, clean, etc.)
  - `Database`: Database technology used
  - `BusinessContext`: Business context and objectives
- **Generator options pattern**: Functional options for Generator configuration
  - `WithAI(client)`: Enable AI-powered generation
  - `NewGenerator()` accepts options for flexible configuration
- **Project scanning capabilities**: Analyze project structure for AI context
  - Detects programming language from package files (go.mod, package.json, etc.)
  - Detects frameworks from directory structure (pages/, components/, cypress/)
  - Generates project structure information for AI prompts
- **Comprehensive test coverage**: Tests for AIGenerator and ProjectScanner
  - Unit tests for all AI generation methods
  - Tests for project scanning and detection
  - Mock AI client for testing without API calls

### Changed
- **`init` command AI integration**: Now uses AI for content generation when available
  - Generates agents, skills, commands with AI instead of static templates
  - Falls back to local templates if AI is not configured or fails
  - Passes Architecture, Database, and BusinessContext to AI for better context
- **`generate` command enhancement**: Supports AI-powered generation
  - Uses AIGenerator when `--ai-provider` and `--api-key` are provided
  - Falls back to local templates when AI is not available
  - Reads project configuration to include Architecture, Database, BusinessContext

### Fixed
- **TemplateContext missing fields**: Added Architecture, Database, BusinessContext fields
  - These fields were needed for AI generation but were missing from the context
  - Added corresponding With* functions for configuration
- **Generator routing**: Fixed GenerateAll() to properly route to AIGenerator or local templates
  - Now correctly checks `useAI` flag before delegating to appropriate generator

### Added
- **`init` command now generates complete project-specific structure**:
  - Generates both example files (from `claude_examples/`) AND project-specific templates
  - Creates language-specific skills (e.g., `language_javascript-cypress.md` for Cypress projects)
  - Generates framework-specific agents, commands, and development guide
  - Single command now handles complete initialization - no need to run `generate` separately
- **`generate` command enhancement**: Now reads project configuration to generate appropriate documentation
  - Reads `.claude/project.yaml` created by `init` command
  - Generates language-specific files based on actual project type (e.g., `language_typescript-cypress.md`)
  - Uses project description and business context for customized documentation
- **Project configuration persistence**: `init` command now saves project configuration to `.claude/project.yaml`
  - Stores all survey answers: name, description, language, framework, architecture, database, type, business context
  - Includes creation timestamp
  - Enables `generate` command to create accurate, project-specific documentation
- **AI Configuration Wizard**: Interactive wizard to guide users through AI provider setup
  - Automatically prompts on first run if no credentials are configured
  - Provider selection using interactive menu (no typos possible)
  - API key validation per provider (prefix and length checks)
  - **Optional model selection**: Choose which AI model to use per provider
    - Predefined model lists for each provider (Claude, OpenAI, z.ai)
    - Custom model input for any model not in the list
    - Press Enter to use the recommended default model
  - Saves configuration to `~/.config/claude-init/config.yaml`
- **`config` command**: New command to manage AI provider configuration
  - `config show`: Display current configuration with masked API keys
  - `config set`: Configure a provider (claude, openai, zai)
  - `config list`: List all available providers with status
  - `config unset`: Remove provider configuration
- **Multi-provider support**: Configure multiple AI providers simultaneously
  - Each provider can have its own API key and model
  - One provider can be marked as default
  - Predefined URLs and models for each provider
- **Provider validation**: API key format validation for all providers
  - Claude: `sk-ant-` prefix, 40+ characters
  - OpenAI: `sk-` prefix, 20+ characters
  - z.ai: UUID format (no prefix), 40+ characters
- **Model management**: Predefined model lists per provider
  - Claude: claude-3-5-sonnet, claude-3-5-haiku, claude-3-opus, claude-3-sonnet
  - OpenAI: gpt-4o, gpt-4o-mini, gpt-4-turbo, gpt-3.5-turbo
  - z.ai: glm-4.6, glm-4, glm-3-turbo
- **ConfigManager**: Centralized configuration management
  - Thread-safe YAML load/save operations
  - Provider CRUD operations
  - Default provider management

### Changed
- **`init` command flow**: Now validates path before running AI wizard
- **AI config location**: AI provider configuration now saved globally per user
  - Previous: `.claude/config.yaml` (local to project)
  - Correct: `~/.config/claude-init/config.yaml` (global user config)
  - Benefits: Credentials shared across all projects, single setup for all work
- **Survey questions now use text input**: Changed from restrictive selects to open text inputs
  - Language, Framework, Architecture, Database, and Project Type are now free text
  - Allows any custom values without being limited to predefined options
  - Optional fields (Framework, Database) clearly marked as "(opcional, presiona Enter para omitir)"
- Improved UX: Provider selection uses interactive menus instead of text input
- Updated YAML configuration format:
  ```yaml
  provider: claude  # default provider
  providers:        # nested provider configs
    claude:
      api_key: sk-ant-xxxxx
      base_url: https://api.anthropic.com
      model: claude-3-5-sonnet-20241022
      max_tokens: 8192
  ```

### Fixed
- **OpenAI client implementation**: Implemented full OpenAI API client
  - Problem: OpenAI client was not implemented ("OpenAI client not yet implemented" error)
  - Solution: Implemented `ValidateAnswers()` and `Recommendation()` methods with full API integration
  - Added HTTP client for making requests to OpenAI API endpoint
  - Correct API endpoint: `https://api.openai.com/v1/chat/completions`
  - Bearer token authentication using API key
  - Response parsing compatible with OpenAI format
- **OpenAI debug logging**: Added comprehensive logging for OpenAI API requests
  - Logs show model, URL, prompt preview (first 500 chars)
  - Logs show request body, response status, and response body
  - Logs show token usage (prompt_tokens, completion_tokens)
  - Logger is set via `SetLogger()` method on `OpenAIClient`
  - Debug logs only shown with `--verbose` flag
  - Useful for debugging API calls and testing
- **Survey select questions missing defaults**: Fixed survey error when running `init` command
  - Error: "default value "" not found in options" when answering language question
  - Root cause: Select questions had empty Default values
  - Solution: Changed survey questions from restrictive Selects to open text inputs
- **generate command nil pointer**: Fixed panic when running `claude-init generate`
  - Error: `panic: runtime error: invalid memory address or nil pointer dereference`
  - Root cause: Logger was not initialized in `runGenerate` function
  - Solution: Added nil check and logger initialization at the start of `runGenerate`
- **Multiline input UX**: Changed from external editor back to standard input for reliability
  - Problem: External editor (vim/nano) could hang or confuse users on exit
  - Solution: Changed `survey.Editor` to standard `survey.Input` with minimum length validation (20 characters)
  - Benefits: Simple, predictable interface that works consistently across all terminals
  - Users can paste text or write a short description without any complexity
  - Clear validation message when description is too short
- **`init` command not generating project-specific files**: Fixed bug where language/framework skills were not generated
  - Problem: Command combined language and framework into single string (e.g., `javascript-cypress`) but templates expected separate values
  - Root cause: `determineProjectType()` returned combined string like `javascript-cypress` but templates looked for `language_javascript.tmpl` and `framework_cypress.tmpl`
  - Solution: Separated language and framework handling in both `init` and `generate` commands
  - `TemplateContext` now receives separate `Language` and `Framework` values
  - Added `normalizeLanguage()` and `normalizeFramework()` functions
  - Created `framework_cypress.tmpl` template for Cypress projects
  - JavaScript projects now use `nodejs` template (language mapping: `javascript` → `nodejs`)
- **`init` command files generated in wrong directory**: Fixed bug where agents/commands files were created in `.claude/` root instead of subdirectories
  - Problem: Files like `architect.md`, `bug-fix.md`, etc. were generated directly in `.claude/` instead of `.claude/agents/` and `.claude/commands/`
  - Root cause: `AgentGenerator.GenerateAll()` and `CommandGenerator.GenerateAll()` wrote files directly to `outputDir` without creating subdirectories
  - Solution: Updated both generators to create `agents/` and `commands/` subdirectories within `outputDir`
  - Updated tests to verify files are created in correct subdirectories
- **z.ai client implementation**: Implemented full z.ai API client
  - Problem: z.ai client was not implemented ("z.ai client not yet implemented" error)
  - Solution: Implemented `ValidateAnswers()` and `RecommendStructure()` methods with full API integration
  - Added HTTP client for making requests to z.ai API endpoint
  - Correct API endpoint: `https://api.z.ai/api/paas/v4/chat/completions`
  - Uses GLM-4.6 model by default with 8192 max tokens
  - Bearer token authentication using API key
  - Response parsing compatible with Anthropic Claude format
- **z.ai debug logging**: Added comprehensive logging for z.ai API requests
  - Logs show model, URL, prompt preview (first 500 chars)
  - Logs show request body, response status, and response body
  - Logs show token usage (input/output tokens)
  - Logger is set via `SetLogger()` method on ZAIClient
  - Debug logs only shown with `--verbose` flag
  - Useful for debugging API calls and testing
- **z.ai API configuration**: Fixed incorrect API URL
  - Previous: `https://api.z.ai/api/anthropic` (incorrect)
  - Correct: `https://api.z.ai/api/paas/v4/` (from official documentation)
  - API key format: UUID with dot separator (e.g., `450e9e2546fb477baa158270c38dc12f.KoQUhnnKBAzSwL5L`)
  - Default model: `glm-4.6`
  - Documentation: https://docs.z.ai/guides/develop/http/introduction
- Fixed logger format string errors for Go 1.25 compatibility in init command
- Corrected step numbering in init command flow comments

## [0.2.0] - 2026-01-18

### Added
- Go 1.25 support
- Initial release of claude-init (renamed from ia-start)

### Changed
- **BREAKING CHANGE**: Project renamed from `ia-start` to `claude-init`
  - Module path changed from `github.com/danielrossellosanchez/ia-start` to `github.com/danielrossellosanchez/claude-init`
  - Binary renamed from `ia-start` to `claude-init`
  - All CLI commands now use `claude-init` prefix
  - Environment variables renamed: `IA_START_*` → `CLAUDE_INIT_*`
  - Config directory changed: `~/.config/ia-start/` → `~/.config/claude-init/`
- Go version updated from 1.23 to 1.25
- GitHub Actions workflows updated to Go 1.25
- Fixed logger tests to comply with Go 1.25 format string rules

## [Unreleased]

### Added
- Initial implementation of claude-init CLI
- Project type detection for 10+ languages:
  - Go, Node.js, TypeScript, Python, Rust, PHP, Ruby, Java, C#, C++
- Framework detection for 6+ frameworks:
  - NestJS, Next.js, Express, React, Django, FastAPI
- AI integration with multiple providers:
  - Anthropic Claude API
  - OpenAI API
  - z.ai API
- Template engine with 30+ embedded files:
  - Agent templates (architect, developer, tester, writer)
  - Skill templates for languages and frameworks
  - Command templates (build, test, run, etc.)
- Six CLI commands:
  - `init`: Initialize .claude/ directory
  - `detect`: Detect project type without creating files
  - `generate`: Generate configuration files
  - `version`: Show version information
  - `completion`: Generate shell completion scripts
- Configuration management:
  - YAML configuration file support
  - Environment variable support
  - Command-line flag overrides
- Interactive prompts for user input
- Dry-run mode for all commands
- Custom config directory support
- Comprehensive test suite with 83.6% coverage
- 27 benchmark tests for performance validation
- 19 linters configured via golangci-lint
- Makefile with 30+ commands for build, test, lint, and more

### Changed
- Improved error messages across all commands
- Enhanced project detection accuracy

### Fixed
- Fixed template generation for edge cases
- Fixed race conditions in concurrent operations

## [0.1.0] - 2026-01-17

### Added
- First official release of claude-init CLI
- Support for Go, Node.js, and Python projects
- Claude API integration
- Basic template generation
- `init` and `detect` commands
- Configuration file support

[Unreleased]: https://github.com/danielrossellosanchez/claude-init/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/danielrossellosanchez/claude-init/releases/tag/v0.2.0
[0.1.0]: https://github.com/danielrossellosanchez/claude-init/releases/tag/v0.1.0
