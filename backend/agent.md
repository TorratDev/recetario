# Backend Go Agent Instructions

## Tech Stack
- **Language:** Go (latest stable version)
- **Router:** chi v5
- **ORM:** GORM v2
- **Database:** PostgreSQL 14+
- **Authentication:** JWT (github.com/golang-jwt/jwt/v5) + bcrypt
- **Validation:** go-playground/validator/v10
- **Environment:** godotenv for local development

## Project Structure
```
/cmd/api          # Application entry point
/internal
  /domain         # Business entities
  /handler        # HTTP handlers
  /service        # Business logic
  /repository     # Data access layer
  /middleware     # HTTP middleware
  /validator      # Custom validators
/pkg              # Public utilities
/migrations       # Database migrations
/config           # Configuration
```

## Architecture Principles

### Clean Architecture Boundaries
- **Handlers** → thin layer, only HTTP concerns (parsing, validation, responses)
- **Services** → business logic, orchestration, no HTTP awareness
- **Repositories** → data access only, return domain models
- **Domain** → pure business entities, no external dependencies

### Dependency Flow
- Dependencies point inward: Handler → Service → Repository → Database
- Use interfaces for testability and flexibility
- Inject dependencies via constructors

### No Global Mutable State
- Use dependency injection
- Configuration via structs passed at initialization
- Database connections via context or injected pools

## API Design

### RESTful Conventions
- `GET /resources` - list
- `GET /resources/{id}` - retrieve
- `POST /resources` - create
- `PUT /resources/{id}` - full update
- `PATCH /resources/{id}` - partial update
- `DELETE /resources/{id}` - delete

### Request/Response Format
- **Content-Type:** `application/json` only
- **Error Response Structure:**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request data",
    "details": [
      {"field": "email", "message": "must be valid email"}
    ]
  }
}
```

### Status Codes
- `200` - Success with response body
- `201` - Resource created
- `204` - Success with no content
- `400` - Validation error
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not found
- `409` - Conflict (e.g., duplicate)
- `500` - Internal server error

### Pagination
- Use query params: `?page=1&limit=20`
- Return metadata: `{"data": [...], "meta": {"total": 100, "page": 1, "limit": 20}}`

## Database Practices

### Migrations
- All schema changes via migrations (e.g., golang-migrate/migrate)
- Never use `AutoMigrate()` in production
- Up and down migrations required
- Sequential naming: `001_initial_schema.up.sql`

### Performance
- Define indexes in migrations per technical plan
- Avoid N+1 queries (use `Preload` or joins)
- Query timeout: 5-10 seconds max
- Target response time: <100ms for simple queries

### Transactions
- Use transactions for multi-step operations
- Repository methods should accept `*gorm.DB` to allow transaction passing
- Services orchestrate transaction boundaries

### Example Repository Pattern
```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uint) (*User, error)
    // Transaction-aware version
    WithTx(tx *gorm.DB) UserRepository
}
```

## Search Implementation

### Full-Text Search
- Use PostgreSQL `tsvector` and `tsquery`
- Create GIN indexes on tsvector columns
- Example:
```sql
ALTER TABLE posts ADD COLUMN search_vector tsvector;
CREATE INDEX posts_search_idx ON posts USING GIN(search_vector);
```

### Filtering
- Prefer database-side filtering with WHERE clauses
- Use parameterized queries to prevent SQL injection
- Build dynamic queries with GORM's chainable methods

## Authentication & Authorization

### JWT Implementation
- Use HS256 or RS256 (prefer RS256 for production)
- Short-lived access tokens (15 min)
- Refresh tokens stored securely (httpOnly cookies)
- Include minimal claims: `user_id`, `role`, `exp`, `iat`

### Password Handling
- bcrypt cost factor: 12-14
- Never log or return passwords
- Validate password strength on registration

### Middleware
- Authentication middleware checks JWT validity
- Authorization middleware checks user permissions
- Rate limiting for sensitive endpoints

## Error Handling

### Guidelines
- **Never silently swallow errors** - always log or return
- Wrap errors with context: `fmt.Errorf("failed to create user: %w", err)`
- Use custom error types for domain errors
- Log errors with structured logging (e.g., zerolog, zap)

### Example
```go
var ErrUserNotFound = errors.New("user not found")

func (s *UserService) GetUser(ctx context.Context, id uint) (*User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("repository error: %w", err)
    }
    return user, nil
}
```

## Code Quality

### Allowed
- Clear, explicit code over clever abstractions
- Standard library patterns
- Well-documented public APIs
- Defensive programming with validation

### Disallowed
- Heavy use of reflection (except in libraries like GORM, validator)
- Magic behavior or hidden side effects
- Silent error swallowing (check every error)
- Raw SQL unless absolutely necessary (use GORM query builder)
- Database queries in handlers

## Testing Strategy

### Unit Tests
- Test services in isolation
- Mock repository interfaces
- Use table-driven tests
- Coverage target: >80% for services

### Integration Tests
- Test handlers with real database
- Use test containers or separate test DB
- Wrap tests in transactions and rollback
- Test authentication/authorization flows

### Example
```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   *User
        wantErr bool
    }{
        {"valid user", &User{Email: "test@example.com"}, false},
        {"duplicate email", &User{Email: "existing@example.com"}, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Performance Guidelines
- Use context deadlines for external calls
- Implement connection pooling (default GORM pool is usually sufficient)
- Profile with pprof for optimization opportunities
- Cache frequently accessed, rarely changing data
- Use prepared statements (GORM does this automatically)

## Security Checklist
- Validate all user input
- Use parameterized queries (no SQL injection)
- Implement rate limiting
- CORS configuration for allowed origins
- HTTPS only in production
- Security headers (helmet-like middleware)
- Input sanitization for XSS prevention