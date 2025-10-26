# Contributing to gomap

Thank you for your interest in contributing to gomap! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- A clear and descriptive title
- Exact steps to reproduce the problem
- Expected behavior vs actual behavior
- Code samples demonstrating the issue
- Your environment (Go version, OS, etc.)

### Suggesting Features

Feature requests are welcome! Please provide:

- A clear and descriptive title
- Detailed description of the proposed feature
- Use cases and benefits
- Code examples of how you'd like to use it

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Write clear commit messages** following conventional commits:

   - `feat: add new feature`
   - `fix: resolve bug`
   - `docs: update documentation`
   - `test: add tests`
   - `refactor: code improvements`
   - `perf: performance improvements`

3. **Ensure tests pass**: Run `make test`
4. **Add tests** for new features
5. **Update documentation** as needed
6. **Run linters**: `make lint`
7. **Format code**: `make fmt`

## Development Setup

### Prerequisites

- Go 1.25 or later
- Make (optional but recommended)

### Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/gomap.git
cd gomap

# Install development tools
make install-tools

# Run tests
make test

# Run linters
make lint
```

## Code Standards

### Style Guide

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Write idiomatic Go code
- Keep functions small and focused
- Document exported functions and types

### Testing

- Write unit tests for all new code
- Maintain or improve code coverage
- Include table-driven tests when appropriate
- Test edge cases and error conditions
- Benchmark performance-critical code

Example test:

```go
func TestFeatureName(t *testing.T) {
    tests := []struct {
        name    string
        input   Input
        want    Output
        wantErr bool
    }{
        {
            name:    "valid case",
            input:   Input{...},
            want:    Output{...},
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Documentation

- Document all exported functions, types, and constants
- Use clear, concise language
- Include code examples for complex features
- Update README.md for user-facing changes

Example documentation:

```go
// Copy maps data from src to dst using the default configuration.
// The destination must be a pointer to a struct.
//
// Example:
//   var dst Destination
//   err := mapper.Copy(&dst, src)
//
// Returns an error if the mapping fails.
func Copy(dst, src interface{}, opts ...Option) error {
    // ...
}
```

## Commit Messages

Use conventional commits:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test changes
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks

Example:

```
feat(mapper): add support for custom field transformers

Adds a new option WithFieldTransformer that allows users to
apply custom transformations to field values during mapping.

Closes #123
```

## Release Process

Releases are managed by maintainers. Version numbers follow [Semantic Versioning](https://semver.org/):

- MAJOR: Breaking changes
- MINOR: New features (backwards compatible)
- PATCH: Bug fixes (backwards compatible)

## Questions?

Feel free to:

- Open an issue with the `question` label
- Start a discussion in GitHub Discussions
- Reach out to maintainers

Thank you for contributing! 🎉

---
