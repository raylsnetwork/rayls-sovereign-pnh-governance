---
name: new-endpoint
description: >
  Scaffold a new REST API endpoint in the API service. Use when the user wants
  to add a new route, endpoint, or API operation including the handler, service,
  repository, and route wiring.
argument-hint: "[resource] [operation]"
---

# New API Endpoint

You are adding a new REST API endpoint to the API service (`cmd/api/`). This follows
the hexagonal architecture: handler (primary adapter) -> service (core) -> repository (secondary adapter).

The user will provide:
- **Resource name** (e.g., `balance`, `participant`) — use `$ARGUMENTS[0]` if provided
- **Operation** (e.g., `GetByChainId`, `List`) — use `$ARGUMENTS[1]` if provided

If any of these are missing, ask the user before proceeding.

## File Map

These are the files involved, ordered from outside-in:

| Layer | File | Purpose |
|-------|------|---------|
| Route | `cmd/api/app/app.go` | Gin route registration |
| Handler | `cmd/api/adapters/handlers/<resource>_handler.go` | HTTP handler with Swagger annotations |
| Service interface | `cmd/api/core/ports.go` (Primary Ports section) | Business operation interface |
| Service impl | `cmd/api/core/<resource>_service.go` | Business logic |
| Repo interface | `cmd/api/core/ports.go` (Secondary Ports section) | Data access interface |
| Repo impl | `cmd/api/adapters/repositories/<resource>_repository.go` | GORM queries |
| DTO | `dto/<resource>.go` | Filters, response structs (if needed) |
| Domain | `domain/<resource>.go` | Domain model (if new entity) |

Use `participant_handler.go`, `participant_service.go`, and `participant_repository.go` as
the reference implementations — they are the simplest complete example.

## Step 1: Define the repository interface in `cmd/api/core/ports.go`

Add the repository interface under the **SECONDARY PORTS** section. Follow the existing
naming convention:

```go
// <Resource>Repository handles database queries for <resources>
type <Resource>Repository interface {
    // FindBy<Field> retrieves a <resource> by <field>
    FindBy<Field>(ctx context.Context, field string) (*domain.<Resource>, error)
}
```

Key patterns:
- All methods take `context.Context` as first argument
- Return `core.ErrRecordNotFound` when a record doesn't exist (not `gorm.ErrRecordNotFound`)
- Use descriptive method names: `FindBy...`, `FindWithFilters`

## Step 2: Define the service interface in `cmd/api/core/ports.go`

Add the service interface under the **PRIMARY PORTS** section:

```go
// <Resource>Service defines the business operations for <resources>
type <Resource>Service interface {
    // Get<Resource>By<Field> retrieves a single <resource> by <field>
    Get<Resource>By<Field>(ctx context.Context, field string) (*domain.<Resource>, error)
}
```

## Step 3: Implement the repository in `cmd/api/adapters/repositories/`

Create or extend `<resource>_repository.go`:

```go
package repositories

import (
    "context"
    "errors"

    "gorm.io/gorm"

    "github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/api/core"
    "github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
)

// <resource>Repository implements core.<Resource>Repository using GORM
type <resource>Repository struct {
    db *gorm.DB
}

// New<Resource>Repository creates a new GORM-based <resource> repository
func New<Resource>Repository(db *gorm.DB) core.<Resource>Repository {
    return &<resource>Repository{db: db}
}

func (r *<resource>Repository) FindBy<Field>(ctx context.Context, field string) (*domain.<Resource>, error) {
    var result domain.<Resource>

    err := r.db.WithContext(ctx).
        Where("<db_column> = ?", field).
        First(&result).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, core.ErrRecordNotFound
        }
        return nil, err
    }

    return &result, nil
}
```

Key patterns:
- Unexported struct, exported constructor returning the interface
- Constructor returns `core.<Resource>Repository` (the interface), not the concrete type
- Always use `r.db.WithContext(ctx)` for context propagation
- Map `gorm.ErrRecordNotFound` to `core.ErrRecordNotFound`

## Step 4: Implement the service in `cmd/api/core/`

Create `<resource>_service.go`:

```go
package core

import (
    "context"
    "errors"

    "github.com/raylsnetwork/rayls-sovereign-pnh-governance/domain"
    "github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
)

// <resource>Service implements the <Resource>Service interface
type <resource>Service struct {
    repo <Resource>Repository
    log  logger.Logger
}

// New<Resource>Service creates a new <resource> service
func New<Resource>Service(repo <Resource>Repository, log logger.Logger) <Resource>Service {
    return &<resource>Service{
        repo: repo,
        log:  log,
    }
}

func (s *<resource>Service) Get<Resource>By<Field>(ctx context.Context, field string) (*domain.<Resource>, error) {
    // Validate input
    if field == "" {
        return nil, NewValidationError("<field>", "<field> cannot be empty")
    }

    // Fetch from repository
    result, err := s.repo.FindBy<Field>(ctx, field)
    if err != nil {
        if errors.Is(err, ErrRecordNotFound) {
            return nil, NewNotFoundError("<resource>", field)
        }
        return nil, NewInternalError("FindBy<Field>", err)
    }

    return result, nil
}
```

Key patterns:
- Unexported struct, exported constructor returning the interface
- Constructor returns `<Resource>Service` (the interface), not the concrete type
- Depends on repository interface, not concrete type
- Use domain errors: `NewValidationError`, `NewNotFoundError`, `NewInternalError`
- Never return raw DB errors — always wrap with the appropriate domain error type

## Step 5: Create the HTTP handler in `cmd/api/adapters/handlers/`

Create or extend `<resource>_handler.go`:

```go
package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/raylsnetwork/rayls-sovereign-pnh-governance/cmd/api/core"
    "github.com/raylsnetwork/rayls-sovereign-pnh-governance/logger"
)

// <Resource>Handler handles HTTP requests for <resource> operations
type <Resource>Handler struct {
    service core.<Resource>Service
    log     logger.Logger
}

// New<Resource>Handler creates a new <resource> handler
func New<Resource>Handler(service core.<Resource>Service, log logger.Logger) *<Resource>Handler {
    return &<Resource>Handler{
        service: service,
        log:     log,
    }
}

// Get<Resource>By<Field> godoc
// @Summary      Get a <Resource> by <Field>
// @Description  Retrieves a <resource> by its <field>.
// @Tags         Audit
// @Param        <field>  path  string  true  "The <field> of the <resource>"
// @Success      200  {object}  domain.<Resource>
// @Failure      400  {object}  map[string]string  "Invalid <field>"
// @Failure      404  {object}  map[string]string  "Not found"
// @Failure      500  {object}  map[string]string  "Internal error"
// @Router       /audit/<resources>/{<field>} [get]
func (h *<Resource>Handler) Get<Resource>By<Field>(c *gin.Context) {
    field := c.Param("<field>")

    result, err := h.service.Get<Resource>By<Field>(c.Request.Context(), field)
    if err != nil {
        HandleError(c, h.log, err)
        return
    }

    c.JSON(http.StatusOK, result)
}
```

Key patterns:
- Depends on service interface, not concrete type
- Always use `c.Request.Context()` when calling the service (not `c` itself)
- Use `HandleError(c, h.log, err)` for error responses — it maps domain errors to HTTP status codes
- Include Swagger annotations (`@Summary`, `@Param`, `@Success`, `@Failure`, `@Router`)
- For query parameter filters, use `c.ShouldBindQuery(&filters)` with a DTO struct
- For list endpoints with filters, apply `middleware.ValidateQueryParams(dto.Filters{})` at route level

## Step 6: Wire everything in `cmd/api/app/app.go`

### 6a. Create the repository:
```go
<resource>Repo := repositories.New<Resource>Repository(dbClient)
```

### 6b. Create the service:
```go
<resource>Service := core.New<Resource>Service(<resource>Repo, log)
```

### 6c. Create the handler:
```go
<resource>Handler := handlers.New<Resource>Handler(<resource>Service, log)
```

### 6d. Register the route:

For audit routes, add under the `auditGroup`:
```go
<resource>Routes := auditGroup.Group("/<resources>")
{
    <resource>Routes.GET("/:<field>", <resource>Handler.Get<Resource>By<Field>)
}
```

For list endpoints with query validation middleware:
```go
<resource>Routes.GET(
    "",
    middleware.ValidateQueryParams(dto.<Resource>ListFilters{}),
    <resource>Handler.Get<Resource>List,
)
```

For auth-protected routes, add outside `auditGroup` with the `authMiddleware`:
```go
router.GET("/<path>", authMiddleware, <resource>Handler.Get<Resource>)
```

## Step 7: Write tests

Create `cmd/api/core/<resource>_service_test.go`. Follow the existing test pattern:

- Use fake repositories from `cmd/api/testutil/fakes.go`
- Use `testutil.StubLogger{}` for the logger dependency
- Use `testify/assert` and `testify/require`
- Test cases: happy path, not found, validation errors, empty results
- Use descriptive test names: `TestResourceService_Operation_ExpectedBehavior`

If the endpoint needs a new fake repository, add it to `cmd/api/testutil/fakes.go`.

## Step 8: Regenerate Swagger docs

```bash
make swagger
```

This runs `swag init` and regenerates `cmd/api/docs/`.

## Step 9: Verify

```bash
make build
make lint
make api-coverage
```

## Checklist

- [ ] Repository interface in `core/ports.go` (Secondary Ports)
- [ ] Service interface in `core/ports.go` (Primary Ports)
- [ ] Repository implementation in `adapters/repositories/`
- [ ] Service implementation in `core/`
- [ ] HTTP handler in `adapters/handlers/` with Swagger annotations
- [ ] Wiring in `app/app.go` (repo -> service -> handler -> route)
- [ ] Service tests in `core/`
- [ ] Fake repository in `testutil/fakes.go` (if new resource)
- [ ] Swagger docs regenerated (`make swagger`)
- [ ] Build passes (`make build`)
- [ ] Lint passes (`make lint`)
