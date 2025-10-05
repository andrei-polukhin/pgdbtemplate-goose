# Contributing to `pgdbtemplate-goose`

We welcome contributions to `pgdbtemplate-goose`! This document provides guidelines
for contributing to the project.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Create a new branch for your feature or bugfix
4. Make your changes
5. Add or update tests as necessary
6. Ensure all tests pass
7. Submit a pull request

## Development Environment

### Prerequisites

- Go 1.21 or later
- PostgreSQL 10 or later (for testing)
- Git

### Setup

```bash
# Clone your fork
git clone https://github.com/yourusername/pgdbtemplate-goose.git
cd pgdbtemplate-goose

# Install dependencies
go mod tidy

# Set up PostgreSQL connection string
export POSTGRES_CONNECTION_STRING="postgres://postgres:password@localhost:5432/postgres?sslmode=disable"

# Run tests
go test -v ./...

# Run tests with race detection
go test -race -v ./...

# Check code formatting
go fmt ./...

# Run linter
go vet ./...
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `go vet` before committing
- Add comments for exported functions and types
- Keep functions focused and testable

## Testing

### Running Tests

```bash
# Set PostgreSQL connection
export POSTGRES_CONNECTION_STRING="postgres://postgres:password@localhost:5432/postgres?sslmode=disable"

# Run all tests
go test -v ./...

# Run with coverage
go test -v -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Run with race detector
go test -race -v ./...
```

### Writing Tests

- Use `quicktest` for assertions
- Create parallel tests when possible
- Use table-driven tests for multiple scenarios
- Clean up resources in defer blocks
- Test error conditions

Example:

```go
func TestMigrationRunner(t *testing.T) {
    t.Parallel()
    c := qt.New(t)
    
    c.Run("test case", func(c *qt.C) {
        c.Parallel()
        
        // Test implementation
        result := doSomething()
        c.Assert(result, qt.Equals, expected)
    })
}
```

## Pull Request Process

1. **Update Documentation**: Update README.md if you change functionality
2. **Add Tests**: Include tests for new features or bug fixes
3. **Run All Tests**: Ensure all tests pass before submitting
4. **Update CHANGELOG**: Add entry to CHANGELOG.md (if exists)
5. **Follow Commit Convention**: Use clear, descriptive commit messages

### Commit Message Format

```
<type>: <short summary>

<optional detailed description>

<optional footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test additions or changes
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks

Example:
```
feat: add support for custom goose options

Add WithGooseOptions function to allow users to pass
custom goose provider options for advanced configurations.

Closes #42
```

## Code Review

All submissions require review. We use GitHub pull requests for this purpose.
Reviewers will look for:

- Code quality and maintainability
- Test coverage
- Documentation completeness
- Adherence to Go best practices
- Security considerations

## Reporting Bugs

When reporting bugs, please include:

- Go version (`go version`)
- PostgreSQL version
- Operating system
- Steps to reproduce
- Expected behavior
- Actual behavior
- Any error messages or logs

Use the [bug report template](https://github.com/andrei-polukhin/pgdbtemplate-goose/issues/new?template=bug_report.md)
when available.

## Suggesting Enhancements

We welcome feature suggestions! Please:

- Check existing issues first
- Describe the use case clearly
- Explain why this enhancement would be useful
- Consider backward compatibility

## Documentation

- Update README.md for user-facing changes
- Add godoc comments for exported symbols
- Update examples if behavior changes
- Keep security considerations current

## Getting Help

- **Questions**: Open a [GitHub Discussion](https://github.com/andrei-polukhin/pgdbtemplate-goose/discussions)
- **Bugs**: Create an [Issue](https://github.com/andrei-polukhin/pgdbtemplate-goose/issues)
- **Security**: See [SECURITY.md](SECURITY.md)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Code of Conduct

We are committed to providing a welcoming and inclusive environment.
Please be respectful and professional in all interactions.

Key principles:
- Be respectful and considerate
- Welcome newcomers
- Focus on constructive feedback
- Accept criticism gracefully
- Prioritize community well-being

## Recognition

Contributors will be acknowledged in:
- Git commit history
- GitHub contributors page
- CONTRIBUTORS.md (for significant contributions)

Thank you for contributing to `pgdbtemplate-goose`! 🎉
