# Contributing to LaunchDate Backend

Thank you for your interest in contributing to LaunchDate Backend! This document provides guidelines and instructions for contributing.

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [Development Setup](#development-setup)
4. [Making Changes](#making-changes)
5. [Testing](#testing)
6. [Submitting Changes](#submitting-changes)
7. [Code Style](#code-style)
8. [Commit Messages](#commit-messages)

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/launchdate-backend.git
   cd launchdate-backend
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/vamosdalian/launchdate-backend.git
   ```

## Development Setup

### Prerequisites

- Go 1.21 or later
- PostgreSQL 15
- Redis 7
- Docker (recommended)

### Quick Setup

1. Copy environment file:
   ```bash
   cp .env.example .env
   ```

2. Start dependencies:
   ```bash
   docker run -d --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:15-alpine
   docker run -d --name redis -p 6379:6379 redis:7-alpine
   ```

3. Create the database:
   ```bash
   docker exec -it postgres psql -U postgres -c "CREATE DATABASE launchdate;"
   ```

4. Run migrations:
   ```bash
   make migrate-up
   ```

5. Install dependencies:
   ```bash
   make deps
   ```

6. Run the application:
   ```bash
   make run
   ```

## Making Changes

### Creating a Branch

Create a feature branch from `main`:

```bash
git checkout main
git pull upstream main
git checkout -b feature/your-feature-name
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Test additions or changes

### Development Workflow

1. Make your changes
2. Add tests for new functionality
3. Run tests: `make test`
4. Run linter: `make lint`
5. Format code: `make fmt`
6. Commit your changes
7. Push to your fork
8. Create a Pull Request

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific package tests
go test -v ./internal/api/...
```

### Linting

We use `golangci-lint` for code quality:

```bash
# Run linter
make lint

# Auto-fix issues (when possible)
golangci-lint run --fix
```

## Testing

### Writing Tests

- Place test files next to the code being tested
- Name test files with `_test.go` suffix
- Follow table-driven test pattern when appropriate
- Mock external dependencies (database, cache, etc.)

Example test:

```go
func TestLaunchService_CreateLaunch(t *testing.T) {
	tests := []struct {
		name    string
		req     *models.CreateLaunchRequest
		wantErr bool
	}{
		{
			name: "valid launch",
			req: &models.CreateLaunchRequest{
				Title:      "Test Launch",
				LaunchDate: time.Now(),
			},
			wantErr: false,
		},
		// Add more test cases...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test implementation
		})
	}
}
```

### Test Coverage

We aim for at least 80% test coverage. Check coverage:

```bash
make coverage
```

## Submitting Changes

### Pull Request Process

1. Update documentation if needed
2. Add tests for new features
3. Ensure all tests pass
4. Ensure linter passes
5. Update CHANGELOG.md (if applicable)
6. Create a Pull Request with a clear description

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests added/updated
- [ ] All tests passing
- [ ] Linter passing

## Related Issues
Closes #123
```

### Review Process

1. At least one approval required
2. All CI checks must pass
3. No merge conflicts
4. Up to date with main branch

## Code Style

### Go Style Guidelines

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments):

- Use `gofmt` for formatting
- Follow Go naming conventions
- Write clear, self-documenting code
- Add comments for exported functions
- Keep functions small and focused

### Project Structure

```
cmd/          - Application entry points
internal/     - Private application code
  api/        - HTTP handlers
  config/     - Configuration
  database/   - Database connections
  middleware/ - HTTP middleware
  models/     - Domain models
  repository/ - Data access layer
  service/    - Business logic
migrations/   - Database migrations
docs/         - Documentation
```

### Best Practices

1. **Error Handling**: Always handle errors explicitly
   ```go
   if err != nil {
       return fmt.Errorf("operation failed: %w", err)
   }
   ```

2. **Context**: Pass context for cancellation
   ```go
   func (s *Service) DoWork(ctx context.Context) error {
       // Use ctx for cancellation
   }
   ```

3. **Interfaces**: Use small, focused interfaces
   ```go
   type Reader interface {
       Read(ctx context.Context, id int64) (*Model, error)
   }
   ```

4. **Constants**: Use constants for magic values
   ```go
   const (
       DefaultLimit = 50
       MaxLimit     = 100
   )
   ```

## Commit Messages

### Format

```
<type>: <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `chore`: Build process or auxiliary tool changes

### Examples

```
feat: add user authentication endpoint

Implement JWT-based authentication for API endpoints.
Includes middleware for token validation.

Closes #45
```

```
fix: resolve database connection pool leak

Connection pool was not properly releasing connections
on error conditions. Added proper cleanup logic.

Fixes #67
```

## API Changes

When modifying the API:

1. Update OpenAPI specification (`docs/openapi.yaml`)
2. Update API documentation (`docs/API.md`)
3. Maintain backward compatibility when possible
4. Document breaking changes clearly
5. Update version number if needed

## Database Migrations

When adding migrations:

1. Create both up and down migrations:
   ```bash
   make migrate-create NAME=add_user_table
   ```

2. Test both directions:
   ```bash
   make migrate-up
   make migrate-down
   make migrate-up
   ```

3. Include migration in PR description
4. Consider data migration for existing deployments

## Documentation

Keep documentation up to date:

- **README.md**: Getting started, basic usage
- **docs/API.md**: API endpoints and examples
- **docs/DEPLOYMENT.md**: Deployment instructions
- **Code comments**: For complex logic
- **OpenAPI spec**: API contract

## Questions?

- Open an issue for bugs or feature requests
- Start a discussion for questions
- Check existing issues before creating new ones

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

Thank you for contributing to LaunchDate Backend! 🚀
