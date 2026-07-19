# Release Process & Observability

## Release process

### Trigger

`.github/workflows/release.yml` builds and pushes `backend/Dockerfile` to
`ghcr.io/<owner>/<repo>` whenever a tag matching `v*.*.*` (e.g. `v1.2.0`) is pushed. This
is separate from `.github/workflows/go.yml`, which runs `go test`/`go build` on every
push/PR to `main` and gates merges — `release.yml` only runs on an explicit release tag
and additionally publishes an artifact.

### Cutting a release

```bash
git checkout main && git pull
git tag v1.2.0
git push origin v1.2.0
```

The workflow: runs the full test suite (never pushes an image that doesn't pass its own
tests), logs into GHCR with the built-in `GITHUB_TOKEN`, then builds and pushes with three
tags: the exact version (`v1.2.0`), the `major.minor` floating tag (`v1.2`), and `latest`.

Versioning follows semver against the **API contract** described in
[`../api/VERSIONING.md`](../api/VERSIONING.md) — a breaking API change is a major bump,
regardless of how small the code diff is; an additive API change or backend-only fix is
minor/patch. `../api/CHANGELOG.md` should get an entry before tagging, not after.

### Rollback

Every previous tag's image stays in GHCR. Rolling back means re-deploying the prior tag's
image — this repo doesn't (yet) prescribe a specific deploy target/orchestrator, so the
exact rollback command depends on wherever the image is actually run. There is no
database migration step to reverse (schema.sql is additive, see `CLAUDE.md`), so a
rollback is just "run the old image" unless a release happened to include a genuinely
destructive schema change (which the additive-only `schema.sql` model makes rare by
construction).

### What's deliberately not here yet

- No CD (automatic deploy after a successful release build) — publishing the image is as
  far as this pipeline goes. Wiring it to an actual deploy target is future work once
  there is a real hosting target to deploy to.
- No image vulnerability scanning step. Worth adding (`docker/scout-action` or
  `anchore/scan-action`) before this is used for anything beyond local/demo deploys.

## Observability

### What exists today

- **Structured logs**: `internal/logger` wraps `log/slog` with JSON output.
  `appmiddleware.RequestLogger` (part of the middleware chain in `cmd/main.go`) attaches a
  correlation ID to every request's context and logs method, URI, duration, and status
  code for every request — this is already wired up, not proposed.
- **`GET /healthz`**: added alongside this doc (`internal/handlers/health.go`). Returns
  `200 {"status":"ok"}` when the process is up and `db.PingContext` succeeds, `503
  {"status":"error"}` otherwise. Intended for container/orchestrator liveness+readiness
  probes and uptime monitors. Not wired into `backend/Dockerfile` as a `HEALTHCHECK`
  directive — the distroless runtime image has no shell/`curl`/`wget` to run one from
  inside the container (see [`DOCKER.md`](DOCKER.md)); probe it from outside the
  container instead (e.g. a Kubernetes `httpGet` readiness probe, or an external uptime
  check hitting the published port).
- **Request-ID propagation**: `chiMiddleware.RequestID` + the correlation ID above means
  every log line for a request can be tied together, and (per `CLAUDE.md`'s middleware
  order) this happens before `Recoverer`, so a panic's recovery log still carries it.

### What's missing (future work, not implemented here)

- **Metrics** (request rate/latency/error-rate, ideally Prometheus-format on a `/metrics`
  endpoint). None exist yet. The search latency numbers in
  [`../search/BENCHMARKS.md`](../search/BENCHMARKS.md) are offline benchmark numbers, not
  live production metrics — there's no dashboard today that would catch a live regression
  automatically.
- **Distributed tracing**. Not applicable yet at single-service scale, but would matter
  once/if the mobile client, web UI, and backend need correlated traces across a slower
  network hop.
- **Alerting**. `BENCHMARKS.md` proposes a p95 > 500ms / 5min alert threshold for search —
  there's no alerting pipeline (Prometheus Alertmanager, Grafana alerts, etc.) to actually
  fire it yet. That threshold is a target to build toward, not a live alert.
- **Log aggregation/shipping**. Logs currently go to stdout as JSON (correct for a
  container — 12-factor "logs as event streams") but nothing yet ships them to a
  queryable backend (Loki, CloudWatch, etc.). `docker-compose.yml` doesn't configure a
  logging driver beyond Docker's default `json-file`.

None of the above blocks the current release pipeline from working — they're gaps to
close as the app gets real traffic, not prerequisites for cutting a release today.
