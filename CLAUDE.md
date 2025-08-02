# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Database Operations
- Start PostgreSQL container: `docker-compose up -d`
- Stop PostgreSQL container: `docker-compose down`
- Reset database (delete all data): `docker-compose down -v`

### Application Commands
- Run the API server: `go run ./cmd/server/main.go`
- Run tests: `make test` or `go test -v ./...`
- Start database using Makefile: `make up`
- Stop database using Makefile: `make down`
- Reset database using Makefile: `make down-v`

### Environment Setup
- Copy `.env.example` to `.env` and configure DATABASE_URL and JWT_SECRET_KEY
- Database runs on port 5432, API server on port 8080
- Use `requests.http` file for API testing with REST Client extension

## Architecture Overview

This is a **Modular Monolith** built with **Clean Architecture** principles, organized into bounded contexts:

### Core Modules (Bounded Contexts)
- **Identidade** (Identity): User registration and authentication with JWT
- **Obras** (Construction Projects): CRUD for projects, stages, worker allocation
- **Pessoal** (Personnel): Employee management, time tracking, payroll
- **Suprimentos** (Supplies): Vendors, materials, quotes/budgets
- **Financeiro** (Financial): Payment processing and financial records

### Architecture Patterns
- **Clean Architecture**: Dependencies point inward to domain layer (`internal/domain`)
- **CQRS**: Separate read/write operations (complex read models like `ObraDashboard`)
- **Event-Driven Communication**: Modules communicate via internal EventBus (`internal/events`)
- **Repository Pattern**: Data access abstracted through interfaces

### Directory Structure
- `cmd/server/main.go`: Application entry point with dependency injection
- `internal/domain/`: Core business entities and interfaces (no external dependencies)
- `internal/service/`: Business logic and use cases
- `internal/handler/http/`: HTTP handlers organized by module
- `internal/infrastructure/repository/postgres/`: Database implementations
- `internal/events/`: Event definitions and handlers
- `internal/platform/bus/`: Internal event bus implementation
- `pkg/`: Shared utilities (auth, security, storage)

### Key Technologies
- **Database**: PostgreSQL with pgx driver
- **HTTP Router**: go-chi/chi/v5
- **Authentication**: JWT tokens in httpOnly cookies
- **Authorization**: Role-based access control (RBAC) with granular permissions
- **Logging**: Structured JSON logging with slog
- **Testing**: Standard Go testing with testify

### Database Schema
- Initialization scripts in `db/init/01-init.sql`
- Uses UUIDs for primary keys
- Soft deletes implemented where appropriate
- Foreign key relationships between modules

### Event-Driven Architecture
- Events defined in `internal/events/` (e.g., `OrcamentoStatusAtualizado`)
- Event handlers in module-specific `events/` subdirectories
- Asynchronous communication between bounded contexts

### Authentication & Authorization
- JWT-based authentication with middleware in `pkg/auth/`
- Granular permissions defined in `internal/authz/roles.go`
- Cookie-based token storage for security
- Permission checks on all protected routes

### Testing Strategy
- Use `go test -v ./...` to run all tests
- Test files follow `*_test.go` naming convention
- Integration tests may require running database container

### API Design
- RESTful endpoints organized by resource
- Consistent error handling and response formats
- Comprehensive examples in `restclient/*.http` files
- Health check endpoint at `/health`