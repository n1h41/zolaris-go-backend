# Zolaris Backend - Agent Guidelines

## Build/Test Commands
- `go build -o bin/app .` - Build the application
- `go test ./...` - Run all tests
- `go test ./internal/utils -v` - Run tests for specific package
- `make start-dev` - Start development environment with Docker Compose
- `make migrate-up` - Run database migrations up

## Project Structure
- `main.go` - Application entry point with Gin router setup
- `internal/` - Core business logic (services, repositories, middleware)
- `api/handlers/` - HTTP request handlers
- `internal/transport/dto/` - Data transfer objects
- `internal/domain/` - Domain entities

## Code Style Guidelines
- Use Go 1.24+ syntax and idioms
- Import order: standard library, third-party, local (grouped with blank lines)
- Package naming: lowercase, single word when possible
- Interface naming: append "Interface" suffix (e.g., `UserRepositoryInterface`)
- Error handling: wrap errors with context using `fmt.Errorf("message: %w", err)`
- Logging: use `log.Printf()` for structured logging with context
- Constructor pattern: `NewServiceName()` functions for dependency injection
- Comments: document exported functions with proper Go doc format
- Testing: use `testing` package, test files end with `_test.go`