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
- Check API compilation: `go build ./cmd/server/main.go`

### Environment Setup
- Copy `.env.example` to `.env` and configure DATABASE_URL and JWT_SECRET_KEY
- Database runs on port 5432, API server on port 8081 (configurable via PORT env var)
- Use `restclient/*.http` files for API testing with REST Client extension
- Server logs are structured JSON format for production readiness

## Architecture Overview

This is a **Modular Monolith** built with **Clean Architecture** principles, organized into bounded contexts:

### Core Modules (Bounded Contexts)
- **Identidade** (Identity): User registration, authentication with JWT, and role-based authorization
- **Obras** (Construction Projects): Project management with financial control, payment schedules, and resource allocation
- **Pessoal** (Personnel): Employee management, time tracking, quinzenal payroll, and payment approval
- **Suprimentos** (Supplies): Vendor management, product catalog, budget quotes with approval workflow
- **Financeiro** (Financial): Complete financial management with accounts receivable, accounts payable, and cash flow

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

## Financial Module Details

The Financial module is a comprehensive system for construction company financial management:

### Key Entities
- **ContaReceber** (Accounts Receivable): Revenue from construction projects
- **ContaPagar** (Accounts Payable): Payments to suppliers and service providers
- **CronogramaRecebimento** (Payment Schedule): Staged payment plans for projects
- **ParcelaContaPagar** (Payment Installments): Installment support for supplier payments

### Financial Features
- **Cash Flow Management**: Real-time tracking of money in/out based on actual transactions
- **Automated Account Creation**: Accounts payable automatically created when budgets are approved
- **Project Revenue Tracking**: Payment schedules by project stages
- **Supplier Payment Control**: Complete payable accounts with installment support
- **Event-Driven Integration**: Financial movements trigger events for other modules

### Cash Flow Calculation
- **Inflows**: Actual received amounts from `contas_receber` with status `RECEBIDO`
- **Outflows**: Paid amounts from `contas_pagar` with status `PAGO` + employee payments
- **Dashboard API**: `/dashboard/fluxo-caixa` provides consolidated financial view

### Integration Points
- **Approved Budgets** → Automatically create accounts payable
- **Project Payment Schedules** → Create accounts receivable
- **Employee Payments** → Register financial outflows
- **Financial Movements** → Update cash flow dashboard

### Database Tables
- `contas_receber`: Accounts receivable with project linkage
- `contas_pagar`: Accounts payable with supplier and budget linkage
- `parcelas_conta_pagar`: Installment support for payables
- `cronograma_recebimentos`: Payment schedules by project stages
- Enhanced `obras` table with financial contract fields

## Development Guidelines

### Code Organization
- Follow Clean Architecture principles with clear layer separation
- Domain entities contain business logic and validation
- Services handle use cases and orchestration
- Handlers manage HTTP concerns only
- Repositories abstract data access

### Event-Driven Communication
- Modules communicate via events published to internal EventBus
- Events are defined in `internal/events/` with typed payloads
- Event handlers are in module-specific `events/` directories
- Example: `OrcamentoStatusAtualizado` event creates accounts payable automatically

### Adding New Features
1. Define domain entities with business methods in `internal/domain/[module]/`
2. Create service layer with DTOs in `internal/service/[module]/`
3. Implement repository interfaces in `internal/infrastructure/repository/postgres/`
4. Add HTTP handlers in `internal/handler/http/[module]/`
5. Wire dependencies in `cmd/server/main.go`
6. Add database migrations to `db/migrations/`
7. Create test files in `restclient/[module].http`

### Database Migrations
- Migration files in `db/migrations/` with format `XXX-description.sql`
- Apply manually: `docker exec -i masterconstrutora-db psql -U user -d mastercostrutora_db < migration.sql`
- Always use UUIDs for primary keys
- Include proper indexes for performance
- Use foreign key constraints for data integrity

### Testing Strategy
- Unit tests for domain entities and services
- Integration tests may require database container
- HTTP endpoint tests using `restclient/*.http` files
- Run tests with `go test -v ./...`

### Documentation
- Module documentation in `docs/MODULO_*.md`
- API documentation includes examples and response formats
- Architecture documentation in `docs/ARCHITECTURE.md`
- Authentication details in `docs/AUTH.md`