# Code Standards

Go coding standards for this project. Project-specific structure and patterns are in `CLAUDE.md`.

## Principles

- **Explicit over implicit** - Be clear about types, return values, and error conditions
- **Simple over clever** - Readable code beats clever one-liners
- **Fail fast** - Validate inputs early and fail with clear error messages
- **Test what matters** - Focus tests on behavior, not implementation details

## Architecture

- **Clean Architecture**: handlers → services → repositories → domain models
- **Interface-driven development** with explicit dependency injection
- **Composition over inheritance**; small, purpose-specific interfaces
- Public functions accept interfaces, not concrete types
- Compile-time interface checks: `var _ Interface = (*Struct)(nil)`

## Go Rules

- Always handle errors (never use `_` to ignore)
- Wrap errors with context: `fmt.Errorf("failed to do something: %w", err)`
- Use `context.Context` as the first parameter
- Avoid global state; use constructor functions for DI
- Defer closing resources to avoid leaks
- Guard shared goroutine state with channels or sync primitives

## Naming

| Type | Convention | Example |
|------|------------|---------|
| Files | snake_case | `user_service.go`, `token_handler.go` |
| Directories | lowercase, no underscores | `handlers/`, `repositories/` |
| Constants | PascalCase | `MaxRetries`, `DefaultTimeout` |
| Boolean variables | Prefix with is/has/can/should | `isActive`, `hasPermission` |

## Import Order

1. Standard library
2. External dependencies
3. Internal packages

```go
import (
    "context"
    "fmt"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/core"
    "github.com/raylsnetwork/rayls-privacy-pnh-governance-api/domain"
)
```

## Testing

- One test function per scenario (not table-driven)
- Name: `TestServiceName_MethodName_Scenario`
- One-line comment at top describing expected behavior
- Blank lines between setup, action, assertions (implicit AAA)
- `testutil` fakes for repositories, `testutil.StubLogger{}` for logging
- `testify/require` for fatal checks, `testify/assert` for the rest

```go
func TestTokenService_GetTokenByResourceId_ReturnsTokenWithBalances(t *testing.T) {
    // Querying a token by resourceId returns the token with its balances
    repo := testutil.NewFakeTokenRepository()
    repo.Tokens = []domain.TokenWithBalancesAndFreezeState{
        buildToken("abcd", "Token A", "TKA"),
    }

    svc := NewTokenService(repo, &testutil.StubLogger{})

    result, err := svc.GetTokenByResourceId(context.Background(), "abcd")

    require.NoError(t, err)
    assert.Equal(t, "abcd", result.ResourceId)
    assert.Equal(t, "Token A", result.Name)
    assert.Equal(t, "TKA", result.Symbol)
}
```

## Error Messages

```
Bad:  "Invalid input"
Good: "resourceId must be a valid hex string"

Bad:  "Connection failed"
Good: "Failed to connect to database at localhost:5432: connection refused"
```

## Comments

**Do comment:** complex business logic ("why"), public API docs (GoDoc), workarounds with issue links.
**Don't comment:** obvious code, every function, commented-out code (delete it).

## Git Commits

```
<type>: <short description>

<optional body explaining why>
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

## Git Safety

Never commit: `.env`, `*.pem`, `*.key`, `credentials.json`, `vendor/`, `build/`, `.DS_Store`, files > 500KB.

## Linting

- Config: `.golangci.yml`
- Key linters: `errcheck`, `govet`, `staticcheck`, `gosec`, `gocritic`
- Race detection: `go test -race ./...`

## Changelog

| Date | Change |
|------|--------|
| 2026-03-05 | Deduplicate: move generic Go standards here, keep only project-specific in CLAUDE.md |
| 2026-02-18 | Merge go-standards.md into this file, delete go-standards.md |
| 2026-02-13 | Remove non-Go content, align with project conventions |
| 2026-01-22 | Add config files, init commands, and deviation detection |
| 2026-01-18 | Initial standards document |
