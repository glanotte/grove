# Development Guide

This document outlines the development practices, tools, and workflows for the Grove project.

## Go Best Practices

### Testing

#### Test Structure
- **Unit Tests**: Test individual functions/methods in isolation
- **Integration Tests**: Test component interactions
- **End-to-End Tests**: Test full workflows
- **Benchmark Tests**: Performance testing

#### Test Organization
```
pkg/
├── worktree/
│   ├── manager.go
│   ├── manager_test.go      # Unit tests
│   ├── config.go
│   └── config_test.go       # Unit tests
internal/
├── testutil/
│   └── testutil.go          # Test utilities and helpers
```

#### Test Naming Conventions
- Test functions: `TestFunctionName`
- Benchmark functions: `BenchmarkFunctionName`
- Table-driven tests: Use descriptive test names
- Test helpers: Start with lowercase (not exported)

#### Table-Driven Tests
```go
func TestSanitizeBranchName(t *testing.T) {
    tests := []struct {
        name       string
        branchName string
        want       string
    }{
        {"simple", "main", "main"},
        {"with slash", "feature/auth", "feature-auth"},
        {"with underscore", "feature_auth", "feature-auth"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := sanitizeBranchName(tt.branchName)
            if got != tt.want {
                t.Errorf("sanitizeBranchName() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

#### Test Utilities
- Use `t.TempDir()` for temporary directories
- Use `t.Cleanup()` for cleanup functions
- Create test helpers in `internal/testutil`
- Use mocks for external dependencies

### Linting and Code Quality

#### golangci-lint Configuration
We use golangci-lint with a comprehensive configuration:

**Enabled Linters:**
- `errcheck` - Check for unchecked errors
- `gosimple` - Simplify code
- `govet` - Go vet
- `staticcheck` - Static analysis
- `gosec` - Security issues
- `misspell` - Spelling errors
- `gocritic` - Code review comments
- `revive` - Fast linter

**Key Settings:**
- Line length limit: 120 characters
- Cyclomatic complexity: 15
- Enable all security checks
- Exclude test files from some checks

#### Running Linters
```bash
# Run all linters
make lint

# Auto-fix issues
make lint-fix

# Run security scanner
make security

# Run all quality checks
make quality
```

### Code Coverage

#### Coverage Targets
- **Minimum**: 70% overall coverage
- **Target**: 80%+ for core packages
- **Exclude**: Test files, generated code, examples

#### Running Coverage
```bash
# Generate coverage report
make test-coverage

# View HTML report
open coverage.html

# Check coverage threshold
go tool cover -func=coverage.out
```

### Build and Development

#### Make Targets
```bash
# Development
make deps           # Download dependencies
make build          # Build binary
make test           # Run tests
make test-race      # Run tests with race detector
make test-bench     # Run benchmarks
make fmt            # Format code
make lint           # Run linters
make vet            # Run go vet

# Quality assurance
make quality        # Run all quality checks
make pre-commit     # Pre-commit checks
make ci             # CI pipeline simulation

# Cross-compilation
make build-all      # Build for all platforms
```

#### Development Workflow
1. **Before starting work:**
   ```bash
   make deps
   make pre-commit
   ```

2. **During development:**
   ```bash
   # Run tests frequently
   make test
   
   # Check for race conditions
   make test-race
   
   # Format and lint
   make fmt
   make lint-fix
   ```

3. **Before committing:**
   ```bash
   make pre-commit
   ```

### CI/CD Pipeline

#### GitHub Actions Workflows

**CI Pipeline (`.github/workflows/ci.yml`):**
- **Test Matrix**: Go 1.20, 1.21 on Ubuntu, Windows, macOS
- **Linting**: golangci-lint + gosec
- **Coverage**: Codecov integration
- **Build**: Cross-compilation
- **Docker**: Container image building

**Release Pipeline (`.github/workflows/release.yml`):**
- **Triggered by**: Git tags (v*)
- **Builds**: Cross-platform binaries
- **Releases**: GitHub releases with assets
- **Homebrew**: Formula updates
- **Docker**: Multi-platform images

#### Pre-commit Hooks
We use pre-commit hooks for:
- Code formatting (`go fmt`)
- Import organization (`goimports`)
- Linting (`golangci-lint`)
- Security scanning (`gosec`)
- Shell script linting (`shellcheck`)
- Dockerfile linting (`hadolint`)

Install pre-commit hooks:
```bash
pip install pre-commit
pre-commit install
```

### Testing Strategies

#### Unit Testing
- Test individual functions in isolation
- Use dependency injection for testability
- Mock external dependencies
- Test edge cases and error conditions

#### Integration Testing
- Test component interactions
- Use real dependencies where possible
- Test configuration loading
- Test template processing

#### End-to-End Testing
- Test complete workflows
- Use temporary directories and Git repos
- Test CLI commands
- Verify file generation

#### Performance Testing
- Benchmark critical paths
- Monitor memory allocations
- Test with large datasets
- Profile CPU usage

### Security Practices

#### Code Security
- Use `gosec` for security scanning
- Avoid hardcoded secrets
- Validate user inputs
- Use secure file permissions

#### Dependencies
- Regularly update dependencies
- Use `go mod tidy`
- Scan for vulnerabilities
- Pin critical dependencies

#### Container Security
- Use minimal base images
- Run as non-root user
- Scan container images
- Keep images updated

### Documentation

#### Code Documentation
- Document all public APIs
- Use examples in documentation
- Keep comments up to date
- Document design decisions

#### README Files
- Project overview
- Quick start guide
- Installation instructions
- Usage examples

#### ADRs (Architecture Decision Records)
- Document significant decisions
- Include context and rationale
- Store in `docs/adr/`

### Debugging and Troubleshooting

#### Common Issues
1. **Build failures**: Check Go version, dependencies
2. **Test failures**: Check environment, cleanup
3. **Linting errors**: Run `make lint-fix`
4. **Coverage drops**: Add tests for new code

#### Debugging Tools
- `go test -v` - Verbose test output
- `go test -race` - Race condition detection
- `go test -bench` - Benchmarking
- `go tool pprof` - Profiling
- `dlv` - Go debugger

#### Logging
- Use structured logging
- Log at appropriate levels
- Include context in logs
- Avoid logging sensitive data

### Performance Optimization

#### Profiling
```bash
# CPU profiling
go test -cpuprofile cpu.prof -bench .

# Memory profiling
go test -memprofile mem.prof -bench .

# Analyze profiles
go tool pprof cpu.prof
```

#### Optimization Guidelines
- Profile before optimizing
- Optimize hot paths first
- Consider memory allocations
- Use appropriate data structures
- Benchmark changes

### Release Process

#### Versioning
- Use semantic versioning (SemVer)
- Tag releases with `v` prefix
- Update CHANGELOG.md
- Create release notes

#### Release Checklist
1. Update version in code
2. Update CHANGELOG.md
3. Run full test suite
4. Create and push tag
5. Verify CI/CD pipeline
6. Test released artifacts
7. Update documentation

This development guide ensures consistency, quality, and maintainability across the project.