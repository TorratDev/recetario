# API Changelog

Tracks changes to the documented contract in `backend/api/openapi.yaml`. See
`docs/api/VERSIONING.md` for what counts as breaking vs. additive.

## 1.0.0 — Initial documented contract

- First formal OpenAPI 3.0 specification of the existing `/api` surface (issue #51). No behavior
  changed — this release documents the API as it already existed in code, including two known
  inconsistencies carried forward intentionally (see `docs/api/VERSIONING.md`):
  - `GET /recipes/` (list) returns a narrower field set than `GET /recipes/{id}/` (detail).
  - Error `code` values vary in granularity depending on which handler served the request.
- Endpoints documented: `auth/{register,login,logout,refresh}`, `recipes/` CRUD, recipe-scoped
  `ingredients` (+ reorder), recipe-scoped + global `tags`, `search` (+ suggestions + popular
  tags), `upload` (+ multiple + delete), `users/profile`, `collections/` CRUD (+ recipe
  add/remove).
