# Running RecipeApp with Docker

## Quick start

```bash
cp .env.example .env
# edit .env and set JWT_SECRET (openssl rand -base64 32)
docker compose up --build
```

The API is then reachable at `http://localhost:8080` (both the JSON API under `/api/*`
and the server-rendered HTMX web UI under `/`).

## What the image contains

`backend/Dockerfile` is a two-stage build:

1. **Build stage** (`golang:1.25-alpine`) — `CGO_ENABLED=0 go build` produces a static
   binary. This only works because the backend already uses `modernc.org/sqlite`, a
   pure-Go SQLite driver (see `CLAUDE.md`) — there's no C driver to link against, so
   there's nothing stopping a fully static build.
2. **Runtime stage** (`gcr.io/distroless/static-debian12:nonroot`) — just the compiled
   binary, the `web/` directory (templates + static assets — the handlers read these
   from disk via `html/template.ParseFiles`/`http.Dir` at request time, not `go:embed`,
   so they have to ship alongside the binary), and nothing else. No shell, no package
   manager, no libc. Runs as the image's built-in unprivileged `nonroot` user
   (UID 65532), not root.

## Persistence

Two things the app writes at runtime need to survive container restarts/recreates:

- The SQLite database (`DB_PATH`, default `/app/data/recipeapp.db` in the container)
- Uploaded recipe images (`UPLOAD_DIR`, default `/app/uploads`)

`docker-compose.yml` mounts both as named volumes (`recipeapp-db`, `recipeapp-uploads`)
so `docker compose down` (without `-v`) preserves them. `docker compose down -v` deletes
both — same caution as `make db-reset` in the bare-metal workflow.

## Configuration

All three of these are read directly by `backend/cmd/main.go`; unset in Docker vs. bare
`go run` only changes the fallback:

| Env var | Compose default | Bare `go run` fallback | Purpose |
|---|---|---|---|
| `JWT_SECRET` | **required**, no default (compose fails fast if unset) | `dev-secret-change-in-production` | Signs/verifies JWTs |
| `DB_PATH` | `/app/data/recipeapp.db` | `recipeapp.db` | SQLite file location |
| `PORT` | `8080` | `8080` | HTTP listen port |
| `UPLOAD_DIR` | `/app/uploads` | `./uploads` | Recipe image upload storage |

`docker-compose.yml` deliberately does **not** default `JWT_SECRET` — running with the
bare-metal dev fallback secret in a container that looks production-shaped is a worse
failure mode than refusing to start.

## Known limitations

- **No `HEALTHCHECK`.** Distroless has no shell or `curl`/`wget` to run one from inside
  the container. A real health endpoint + how to probe it from outside the container
  (or via a Go-based/scratch healthcheck binary) is tracked separately in
  [`RELEASE_AND_OBSERVABILITY.md`](RELEASE_AND_OBSERVABILITY.md).
- **No CI image build/push yet.** This Dockerfile is for local/manual use; wiring it into
  CI (build on merge, push to a registry, tag by release) is also tracked in
  `RELEASE_AND_OBSERVABILITY.md`.
- **Single container, no reverse proxy/TLS.** Fine for local use or behind an existing
  ingress/load balancer; this compose file doesn't attempt to add either.

## Building the image standalone (without compose)

```bash
docker build -t recipeapp:local ./backend
docker run -p 8080:8080 \
  -e JWT_SECRET=$(openssl rand -base64 32) \
  -v recipeapp-db:/app/data \
  -v recipeapp-uploads:/app/uploads \
  recipeapp:local
```
