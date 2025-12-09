# Claude AI Guidelines for Markdown Viewer/Editor

## Project Overview

This is a native desktop markdown viewer and editor application built with Go and the Fyne GUI toolkit.

## Development Guidelines

### Testing Requirements

**CRITICAL**: All business logic MUST have corresponding tests.

1. **Test Coverage**:
   - Write tests for all functions that handle file operations (reading, writing, creating, deleting)
   - Write tests for markdown parsing and preview generation
   - Write tests for file list management and filtering
   - Use table-driven tests for functions with multiple scenarios
   - Aim for >80% code coverage

2. **Test File Organization**:
   - Place tests in `*_test.go` files alongside the code they test
   - Use descriptive test names: `TestFunctionName_Scenario_ExpectedResult`
   - Example: `TestLoadFile_ValidMarkdown_ReturnsContent`

3. **Running Tests**:
   ```bash
   # Run all tests
   go test ./...

   # Run tests with coverage
   go test -cover ./...

   # Generate coverage report
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
   ```

### Linting Requirements

**CRITICAL**: All code MUST pass linting before pushing.

1. **Linting Tools**:
   - Use `golangci-lint` for comprehensive linting
   - Configuration in `.golangci.yml`

2. **Running Linters**:
   ```bash
   # Install golangci-lint
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

   # Run linter
   golangci-lint run

   # Run with auto-fix
   golangci-lint run --fix
   ```

3. **Standard Formatting**:
   ```bash
   # Format all Go files
   go fmt ./...

   # Vet for common issues
   go vet ./...
   ```

### Git Workflow

**CRITICAL**: Use Git for all revision history. Never commit untested code.

1. **Branch Strategy**:
   - `main` branch is protected and always passing tests
   - Create feature branches: `feature/description`
   - Create bugfix branches: `fix/description`

2. **Commit Messages**:
   - Use conventional commits format:
     ```
     type(scope): subject

     body (optional)

     footer (optional)
     ```
   - Types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`
   - Examples:
     - `feat(editor): add auto-save functionality`
     - `fix(preview): correct markdown rendering for code blocks`
     - `test(files): add tests for file loading edge cases`

3. **Pre-Push Checklist**:
   ```bash
   # 1. Format code
   go fmt ./...

   # 2. Run linter
   golangci-lint run

   # 3. Run tests
   go test ./...

   # 4. Build to verify
   go build -o markdown-viewer-editor main.go
   ```

### GitHub Actions CI/CD

The project uses GitHub Actions for continuous integration.

1. **Automated Checks** (runs on every push and PR):
   - Go formatting check (`go fmt`)
   - Linting (`golangci-lint`)
   - Tests (`go test`)
   - Build verification

2. **Workflow File**: `.github/workflows/ci.yml`

3. **Pull Request Requirements**:
   - All CI checks must pass (green)
   - Code review approval required
   - No merge until tests pass

### Code Quality Standards

1. **Error Handling**:
   - Always check and handle errors
   - Provide meaningful error messages
   - Use `fmt.Errorf` to wrap errors with context

2. **Documentation**:
   - Add GoDoc comments for all exported functions and types
   - Keep README.md up to date
   - Document complex logic inline

3. **Code Organization**:
   - Keep functions small and focused (ideally <50 lines)
   - Extract business logic from GUI code
   - Use meaningful variable and function names

### Before Committing Checklist

```bash
# Run this before every commit
./pre-commit-check.sh
```

Or manually:
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] All tests pass (`go test ./...`)
- [ ] New code has tests
- [ ] Documentation updated if needed
- [ ] Commit message follows conventional format

## Testing Examples

### Example Test Structure

```go
package main

import (
    "testing"
    "os"
    "path/filepath"
)

func TestLoadFile_ValidMarkdown_ReturnsContent(t *testing.T) {
    // Arrange
    tmpDir := t.TempDir()
    testFile := filepath.Join(tmpDir, "test.md")
    expectedContent := "# Test\n\nContent"
    os.WriteFile(testFile, []byte(expectedContent), 0644)

    // Act
    content, err := loadFileContent(testFile)

    // Assert
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if content != expectedContent {
        t.Errorf("got %q, want %q", content, expectedContent)
    }
}

func TestLoadFile_NonExistent_ReturnsError(t *testing.T) {
    // Act
    _, err := loadFileContent("/nonexistent/file.md")

    // Assert
    if err == nil {
        t.Error("expected error for nonexistent file, got nil")
    }
}
```

### Table-Driven Tests

```go
func TestMarkdownExtension(t *testing.T) {
    tests := []struct {
        name     string
        filename string
        want     bool
    }{
        {"valid md", "test.md", true},
        {"valid MD", "test.MD", true},
        {"valid markdown", "test.markdown", true},
        {"invalid txt", "test.txt", false},
        {"no extension", "test", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := isMarkdownFile(tt.filename)
            if got != tt.want {
                t.Errorf("isMarkdownFile(%q) = %v, want %v",
                    tt.filename, got, tt.want)
            }
        })
    }
}
```

## Continuous Improvement

When adding new features:
1. Write tests FIRST (TDD approach preferred)
2. Implement the feature
3. Ensure tests pass
4. Run linter and fix issues
5. Update documentation
6. Create pull request
7. Wait for CI to pass
8. Get code review
9. Merge to main

## AI Agent Guidelines

When working on this codebase:
- **NEVER** skip writing tests for business logic
- **ALWAYS** run `go fmt` before committing
- **ALWAYS** run `golangci-lint run` before pushing
- **ALWAYS** ensure `go test ./...` passes
- Use Git commit messages that follow conventional commits
- Ask for clarification if requirements are unclear
- Suggest improvements to test coverage when relevant

## Resources

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Testing](https://golang.org/pkg/testing/)
- [Fyne Documentation](https://developer.fyne.io/)
- [Conventional Commits](https://www.conventionalcommits.org/)
