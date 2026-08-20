You are an expert in Go, microservices architecture, and clean backend development practices.

Go coding standards (naming, imports, testing patterns, git conventions, error messages) are in `.claude/rules/code-standards.md`.
This file covers **project-specific** structure, patterns, and tooling for governance-api.

## Project Structure

Multi-service architecture with three applications:
- `cmd/api/` - REST API service (Gin framework) for audit queries and governance operations
- `cmd/listener/` - Blockchain event listener service
- `cmd/flagger/` - Transaction validation and flagging service
- `domain/` - Domain models and repository interfaces
- `repository/` - Repository implementations (API service)
- `dto/` - Data Transfer Objects (filters, responses, pagination)
- `migrations/` - SQL migrations (golang-migrate)
- `types/` - Type definitions and enum mappings
- `config/` - Configuration management (Viper)
- `logger/` - Structured logging (slog-based)

## Gin Framework

- HTTP handlers in `adapters/handlers/` with Swagger annotations (`@Summary`, `@Param`, `@Success`, `@Failure`)
- Middleware for authentication (`RequireAuth`) and query validation
- `HandleError()` maps domain errors (`NotFoundError`, `ValidationError`, `InternalError`) to HTTP status codes
- Routes grouped by resource under `/audit` base path
- JWT stored in HTTP-only cookies (`Authorization` cookie)

## GORM and Database

- Custom `Model` struct with UUID primary key and UTC timestamp hooks (`BeforeCreate`/`BeforeUpdate`)
- Custom types (`BigInt`, `DecimalArray`, `StringArray`) implement `driver.Valuer` and `sql.Scanner`
- Migrations in `/migrations` named: `{timestamp}_{description}.{up|down}.sql`
- `shopspring/decimal` for precise financial calculations

## Smart Contract Integration

- Contract ABIs generated via `abigen` in versioned packages (e.g., `TokenCoreV1`)
- Event handlers implement: `ContractName() string`, `Handle(ctx context.Context, log ContractLog) error`, `Name() string`
- `go-ethereum` for RPC calls with retry and timeout patterns
- Blocks processed in configurable batches (`BatchSize` in config)
- Contract addresses configured via environment variables

## Error Handling

- `cockroachdb/errors` for stack traces via `withstack.Wrap()`
- Propagate errors with context, never swallow silently
- Panic on critical initialization failures (DB, config)
- `HandleError()` maps domain errors to HTTP status codes in handlers

## Configuration

- Environment variables via Viper's `ExperimentalBindStruct`
- Fallback: `.env` files
- `mapstructure` struct tags for field mapping
- `validator` library for required fields

## Logging

- `slog` with colored text handler (dev) / JSON handler (production)
- Log levels: Debug, Info, Warn, Error

## Naming Conventions

- Receiver names: short lowercase (`r`, `p`, `t`, `s`)
- Status/role enums: `uint8` with bidirectional string mappings in `types/`
- Methods: descriptive verbs (`FindByAuditParameters`, `GetByResourceId`)

## Testing (project-specific)

- Listener tests use `gomock` (generated mocks in `mocks/` directories)
- Integration tests use `ory/dockertest` with PostgreSQL containers
- Use `//go:build ignore` tag for slow integration tests

## Build Commands

```bash
make build      # Build all services
make lint       # Run golangci-lint
make swagger    # Regenerate OpenAPI docs
swag init       # Generate Swagger from annotations
```
