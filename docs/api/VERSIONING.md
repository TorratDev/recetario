# API Versioning & Compatibility Policy

## Current state

The API is served entirely under a flat `/api` prefix (`backend/cmd/main.go`) with no version
segment, header, or field anywhere. `backend/api/openapi.yaml` (added in #51) documents this as
`1.0.0` — the first formally versioned snapshot of the contract, not a claim that versioning
machinery already exists in the code.

## Versioning scheme

**Path-based versioning**: the next breaking change introduces `/api/v2/...`, mounted alongside
the existing `/api/...` routes via an additional top-level `r.Route("/api/v2", ...)` in
`cmd/main.go` (chi already nests routes this way — no new routing mechanism needed). The
unversioned `/api` prefix is treated as an implicit `v1` for as long as it's the only version
that exists; it is not renamed to `/api/v1` retroactively, since that path change would itself be
breaking for existing clients (web UI, Android) with no benefit.

Rejected alternatives:
- **Header/media-type versioning** (`Accept: application/vnd.recipeapp.v2+json`) — harder to
  test with a browser/HTMX client and curl, no material benefit over path versioning at this
  scale.
- **Query-param versioning** (`?version=2`) — easy to omit by accident, weaker as a contract.

## Compatibility rule

Within a major version (`/api`, then `/api/v2`, ...), only **additive, backward-compatible**
changes are allowed:
- New optional request fields, new response fields, new endpoints, new query parameters.
- Loosening a validation rule (e.g. allowing a previously-rejected value).

Anything else — removing/renaming a field, changing a field's type or semantics, tightening
validation, changing a status code for an existing case, changing the shape of an error `code` —
requires a new major version.

## Two known real inconsistencies (documented, not silently fixed)

These exist in the current `1.0.0` contract and are called out here explicitly so they are fixed
under a version bump, not patched quietly under the same version:

1. **List vs. detail field mismatch** — `GET /recipes/` returns a narrower field set
   (`RecipeSummary`) than `GET /recipes/{id}/` (`Recipe`). Unifying these would add/remove fields
   clients may already depend on the absence/presence of — treat as a `v2` change.
2. **Inconsistent error codes** — handlers that call `WriteJSONError` directly produce specific
   `code` values (e.g. `RECIPE_NOT_FOUND`), while handlers using plain `http.Error` get rewritten
   by the `ErrorHandler` middleware into generic codes (`NOT_FOUND`, `BAD_REQUEST`, ...). A client
   that pattern-matches on `code` today sees different granularity depending on which endpoint it
   called. Standardizing all handlers onto `WriteJSONError` is a `v2`-worthy change to the `code`
   enum, even though the envelope shape itself (`ErrorResponse`) stays the same.

## Deprecation process

1. Mark the old version's affected endpoint(s) as deprecated in `backend/api/openapi.yaml`
   (`deprecated: true`) and in `docs/api/CHANGELOG.md`.
2. Emit `Deprecation: true` and `Sunset: <date>` response headers on the deprecated endpoint for
   at least one full release cycle before removal.
3. Record the sunset date and replacement endpoint in `docs/api/CHANGELOG.md`.
4. Only remove the deprecated version/endpoint after its `Sunset` date has passed.

## Breaking-change checklist (use before merging any API change)

- [ ] Does this remove or rename a response field, or change its type?
- [ ] Does this change an existing status code for a case that previously returned a different one?
- [ ] Does this tighten request validation (a previously-accepted request would now be rejected)?
- [ ] Does this change the meaning of an existing `code` value in `ErrorResponse`?
- [ ] Does this change default behavior for an existing endpoint (e.g. default sort order, default
      pagination size)?

If any box is checked: this is a breaking change — it goes in the next major version (`/api/v2`),
not the current one, and must be recorded in `docs/api/CHANGELOG.md`.

## Governance

Any PR touching `backend/api/openapi.yaml` must also update `docs/api/CHANGELOG.md` in the same
PR, and requires review against this checklist before merge.
