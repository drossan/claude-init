# Contributing to claude-init

Thank you for considering contributing to claude-init! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Commit Messages](#commit-messages)
- [Pull Requests](#pull-requests)

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

### Prerequisites

- Go 1.25 or higher
- Git
- make (recommended, but not required)

### Setting Up Your Development Environment

1. Fork the repository on GitHub
2. Clone your fork:

```bash
git clone https://github.com/YOUR_USERNAME/claude-init.git
cd claude-init
```

3. Install dependencies:

```bash
make deps
```

4. Install development tools:

```bash
make install-tools
```

5. Verify your setup:

```bash
make check
```

## Development Workflow

### 1. Create a Branch

Create a new branch for your contribution:

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Test improvements
- `chore/` - Maintenance tasks

### 2. Make Your Changes

Follow the coding standards outlined below.

### 3. Test Your Changes

Run the test suite:

```bash
# Run all tests
make test

# Run with coverage
make test-cover

# Run with race detector
make test-race

# Run integration tests
make test-integration
```

### 4. Lint Your Code

Run linters and formatters:

```bash
# Format code
make fmt

# Run linters
make lint

# Run static analysis
make vet
```

### 5. Commit Your Changes

Follow the commit message guidelines below.

```bash
git add .
git commit -m "feat: add new feature"
```

### 6. Push to Your Fork

```bash
git push origin feature/your-feature-name
```

### 7. Create a Pull Request

Go to the original repository and create a pull request from your fork.

## Coding Standards

### Go Conventions

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting (run `make fmt`)
- Follow standard Go package layouts
- Use meaningful variable and function names

### TDD Approach

This project follows Test-Driven Development (TDD):

1. Write tests first
2. Run tests (they should fail)
3. Write minimal code to pass tests
4. Refactor if needed
5. Repeat

Example:

```go
// First, write the test
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}

// Then, write the implementation
func Add(a, b int) int {
    return a + b
}
```

### Documentation

All exported functions must have godoc comments:

```go
// Package detector provides functionality to detect project types
// and technology stacks by analyzing project files and directory structure.
package detector

// Detect analyzes the project at the given path and returns
// a ProjectInfo struct containing the detected project type,
// technology stack, and other metadata.
//
// If the path does not exist or cannot be read, an error is returned.
func (d *Detector) Detect(path string) (ProjectInfo, error) {
    // ...
}
```

### Error Handling

- Always handle errors
- Use descriptive error messages
- Wrap errors with context:

```go
if err != nil {
    return fmt.Errorf("failed to detect project: %w", err)
}
```

### Code Organization

- Keep functions small and focused
- Use interfaces for abstraction
- Avoid deep nesting
- Use table-driven tests for multiple cases

## Testing

### Writing Tests

- Write tests for all new functionality
- Aim for >80% code coverage
- Use table-driven tests for multiple test cases
- Mock external dependencies

Example of table-driven tests:

```go
func TestDetectProjectType(t *testing.T) {
    tests := []struct {
        name    string
        path    string
        want    string
        wantErr bool
    }{
        {
            name:    "Go project",
            path:    "testdata/go-project",
            want:    "go",
            wantErr: false,
        },
        {
            name:    "Node.js project",
            path:    "testdata/nodejs-project",
            want:    "nodejs",
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := DetectProjectType(tt.path)
            if (err != nil) != tt.wantErr {
                t.Errorf("DetectProjectType() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("DetectProjectType() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Test Organization

- Unit tests: `*_test.go` files in the same package
- Integration tests: `tests/integration/` directory
- Test data: `testdata/` directories

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/detector/...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run with race detector
go test -race ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/) specification:

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `chore`: Maintenance tasks
- `perf`: Performance improvements

### Examples

```
feat(init): add support for custom config directory

Add a --config-dir flag to the init command to allow users to
specify a custom configuration directory name instead of .claude/

Closes #123
```

```
fix(detect): correct framework detection for Next.js

The detector was incorrectly identifying Next.js projects as
plain React projects. This fix checks for next.config.js
to properly identify Next.js.

Fixes #456
```

```
docs: update installation instructions

Add instructions for installing from pre-built binaries
and update the go install command.
```

## Pull Requests

### PR Guidelines

1. **Title**: Use a clear, descriptive title following conventional commits
2. **Description**: Explain what changes you made and why
3. **Linked Issues**: Reference related issues (e.g., "Closes #123")
4. **Tests**: Include tests for new functionality
5. **Docs**: Update documentation if needed

### PR Checklist

- [ ] Tests pass locally (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] Linters pass (`make lint`)
- [ ] Coverage is maintained or improved
- [ ] Documentation is updated
- [ ] Commit messages follow conventions

### Review Process

1. Automated checks must pass
2. At least one maintainer approval required
3. Address review comments
4. Squash commits if needed
5. Merge when approved

## Getting Help

- Open an issue for bugs or feature requests
- Start a discussion for questions
- Check existing issues and discussions

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
