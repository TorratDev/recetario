# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Multi-platform recipe management app: a Go backend (`backend/`) serving both a JSON API (`/api/*`) and a server-rendered HTMX web UI (`/`), plus an early-stage native Android client (`mobile/`, mostly scaffold today).

## Commands

Run from `backend/` unless noted.

```bash
go run cmd/main.go          # run the server directly (listens on :8080)
make run                    # build + run (same as above)
make build                  # go build -o recipe-server ./cmd
go test ./...                # all tests
make test                   # same
go test ./internal/handlers/... -run TestName   # single package / single test
make lint                   # golangci-lint run
make db-reset               # delete recipeapp.db (auto-recreated + reseeded on next run)
```

No manual DB setup: on startup the server opens `recipeapp.db` (SQLite, pure-Go `modernc.org/sqlite`, `CGO_ENABLED=0`), applies `internal/database/schema.sql` (embedded via `go:embed`), and seeds sample data.

CI (`.github/workflows/go.yml`) runs `go test ./...` then `go build ./...` on push/PR to `main`.

## Architecture

Request flow: `cmd/main.go` (wiring — opens the DB, applies schema, seeds data, constructs repositories → handlers, wires chi routes) → `internal/handlers` (thin: parse/validate → call repository → respond) → `internal/repositories` (raw parameterized `database/sql`, one struct per aggregate) → SQLite. `internal/services` and `internal/storage` exist but are intentionally empty — add a service only when logic outgrows a handler.

Key points that span multiple files:

- **No ORM.** Every query is raw, parameterized SQL (`?` placeholders) — never string-concatenated user input.
- **Dual response mode (HTMX):** handlers check `HX-Request: true` and branch between rendering an `html/template` partial (from `backend/web/templates/`) and writing JSON. This pattern runs throughout `internal/handlers`.
- **Auth:** JWT (`golang-jwt/jwt/v5`) + bcrypt, wired via `appmiddleware.NewAuthService` in `cmd/main.go`. Routes opt into `authService.AuthMiddleware` (required) or `OptionalAuthMiddleware` (optional) per-route. Resource ownership always comes from the auth context (`GetUserID(ctx)`), never the request body.
- **Middleware chain order** (set in `cmd/main.go`, keep intact): RequestID → Recoverer → CanonicalLoopbackHost → RequestLogger → ErrorHandler → CORS → RateLimit → SecurityHeaders.
- **Schema:** single source of truth is `internal/database/schema.sql`; no migration tool — it's embedded and applied at startup.
- Module path is `recipe-app`; internal imports are `recipe-app/internal/...`.
- **Mobile (`mobile/`):** native Android (Kotlin), multi-module Gradle (`app`, `data`, `domain`, `network`) — scaffold stage, little implemented yet.

## Authoritative rule docs

This repo carries detailed, up-to-date agent instructions — read the relevant one before making non-trivial backend changes rather than relying on this summary:

- [`.github/copilot-instructions.md`](.github/copilot-instructions.md) — global engineering principles, dependency policy, and a "never violate" safety/security list.
- [`.github/instructions/go.instructions.md`](.github/instructions/go.instructions.md) (`applyTo: backend/**/*.go`) — canonical Go rules: handler/repository/model conventions, status codes, testing pattern, prohibited list.
- [`.github/references/go-shared-rules.md`](.github/references/go-shared-rules.md) — condensed checklist version of the above.
- [`.github/AGENTS.md`](.github/AGENTS.md) — indexes repo-specific agents (`implement-feature`, `review-codebase`, `update-deps`) and skills (`go-feature`, `go-testing`, `project-structure`, `git-commit`, `plan-to-issues`) under `.github/agents/` and `.github/skills/`.

`backend/agent.md` is a higher-level summary of the same rules; where it conflicts with `.github/`, **follow `.github/`** — it describes the code as it exists today.

## Known debt (respect, don't copy)

- Some handlers take a `ctx interface{}` param and call `r.Context()` internally instead of threading `context.Context` through. Don't propagate this pattern in new code; don't silently "fix" existing signatures unless asked.
- `CreateRecipe` writes its children (ingredients, instructions, tags) without a transaction. Prefer wrapping related writes in a `db.BeginTx`/`tx.Commit` transaction when touching that code.
