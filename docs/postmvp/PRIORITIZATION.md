# Post-MVP Initiative Prioritization

Scope: everything left once epic #48's tasks (#50–#58) land — the MVP itself (backend API
+ web UI + basic Android auth/recipes) is covered by those tasks and `docs/ROADMAP.md`,
not this doc. This ranks what comes *after*, based on gaps found while executing that
epic — not a speculative wishlist. Every item below traces to either an unused/partial
piece of the existing codebase or an explicitly-noted limitation in another doc from this
epic.

## Method

Each initiative is scored on three axes, 1 (low) – 3 (high):

- **Impact** — how much it improves the product for real users, not internal tooling.
- **Effort** — inverted (3 = cheap, 1 = expensive) so higher is always "more attractive."
- **Confidence** — how sure we are the scope/estimate is right, given what's already in
  the codebase (an unused model column is higher confidence than a feature with no
  existing scaffolding at all).

Priority tier = rough judgment from the three scores, not a mechanical formula — a 3-3-3
item is obviously P0, but ties are broken by dependency position (unblocking other work
ranks higher) and by risk (user-facing correctness/security gaps outrank convenience
features at equal impact).

## Ranked initiatives

### P0 — do next

| # | Initiative | Impact | Effort | Confidence | Why P0 |
|---|---|---|---|---|---|
| 1 | **Implement the FTS5 migration** (`docs/search/FTS_DESIGN.md`) | 3 | 2 | 3 | Design is already done and reviewed (#55); `docs/search/BENCHMARKS.md` proved the current LIKE-based search is not just slow but **quadratic** (35.9s at 5,000 recipes) — this stops being a nice-to-have the moment the recipe count grows past low thousands. Highest confidence item on this list because the design doc already specifies exact SQL. |
| 2 | **Recipe edit/delete on Android** | 3 | 3 | 3 | The backend already has `PUT`/`DELETE /recipes/{id}` (see `openapi.yaml`) and the mobile app already has the repository/DI/navigation scaffolding from #54 — this is "wire up two more screens using the pattern that already exists," not new architecture. Currently a user can create a recipe on Android but never fix a typo or remove it, which is a glaring gap for a "recipe management app." |
| 3 | **Image upload from Android** | 2 | 2 | 2 | Backend already has working multipart upload endpoints (`/api/upload`, `/api/upload/multiple`) that the web UI presumably exercises, but Android has zero camera/gallery integration. A recipe app where mobile users can't attach a photo is missing a core expected feature. Effort/confidence are middling because it touches platform permissions (camera/storage) not yet dealt with anywhere in `mobile/`. |

### P1 — after P0, real value but lower urgency or higher risk

| # | Initiative | Impact | Effort | Confidence | Why P1 |
|---|---|---|---|---|---|
| 4 | **Pagination/infinite scroll on Android recipe list** | 2 | 3 | 3 | `RecipeApi.getRecipes` already accepts `limit`/`offset`; the mobile list screen just doesn't use them (fixed `limit=20, offset=0`). Cheap, but low-impact until item #1 makes it possible to *have* thousands of searchable recipes worth paginating through — sequenced after FTS. |
| 5 | **Collections on Android** | 2 | 1 | 2 | Backend has a full collections API (`/api/collections/*`) and the web UI already integrated it (`feat(web): integrate search and collections web UI`), but there is no mobile domain model, repository, or screen for collections at all — this is a bigger lift than #2/#3 (new vertical, not "add a button to an existing one"), hence P1 not P0 despite similar impact. |
| 6 | **Ratings & reviews** | 2 | 1 | 1 | `models.Rating` already exists in `backend/internal/models/user.go` but is completely unused — no repository, no handler, no schema table for it beyond the struct. Real user-facing value (helps users pick a good recipe) but lowest confidence here: the existing struct is a stub, not a design, so this needs real scoping (aggregate rating storage, one-rating-per-user-per-recipe constraint, API surface) before effort can be estimated honestly. |
| 7 | **Metrics + alerting** (`docs/ops/RELEASE_AND_OBSERVABILITY.md` gaps) | 2 | 2 | 2 | Structured logs and `/healthz` exist (#58); Prometheus-style metrics, a dashboard, and live alerting on the search-latency SLO proposed in `BENCHMARKS.md` do not. Impact is operational, not user-facing, so it's P1 even though it unblocks safely operating everything above at real scale. |

### P2 — valuable eventually, genuinely speculative or expensive right now

| # | Initiative | Impact | Effort | Confidence | Why P2 |
|---|---|---|---|---|---|
| 8 | **Nutrition info** | 1 | 1 | 1 | `models.NutritionInfo` exists, unused, same as `Rating` — but sourcing accurate nutrition data (manual entry is error-prone and low-trust; an external API is a new integration, cost, and rate-limit surface) makes this the least-scoped item on the list. Don't start until there's a concrete answer to "where does the data come from." |
| 9 | **CD / automatic deploy target** | 2 | 1 | 1 | `release.yml` (#58) publishes an image to GHCR but nothing deploys it — because there is no chosen hosting target yet. This is blocked on a decision (where does this actually run?) that's outside engineering scope, not an effort estimate; revisit once that's answered. |
| 10 | **Systematic i18n** | 1 | 1 | 2 | `appmiddleware/error.go` already has ad hoc Spanish user-facing error strings alongside English ones (`"No se pudo completar la solicitud."`) — real bilingual signal exists, but formalizing it (a proper i18n library, translated templates, locale negotiation) is a cross-cutting rewrite touching every handler and template, for a benefit that's nice but not blocking anything else on this list. |
| 11 | **Social features** (public sharing, following users, activity feed) | 3 | 1 | 1 | Highest *potential* impact on the whole list but zero existing scaffolding anywhere in the codebase (no visibility/privacy model beyond `collections.is_public`, no follow graph, no feed). Genuinely a v2-of-the-product-scale initiative — sequenced last on purpose; starting here before the P0/P1 foundation (real search, complete CRUD, observability) would be building on sand. |

## Key dependencies across initiatives

- **#1 (FTS5) should land before #4 (pagination)** — paginating through a LIKE-scanned
  list that degrades quadratically just makes the problem more visible sooner, not
  better.
- **#2 and #3 (Android edit/delete, image upload) are independent of each other and of
  #1** — both build directly on the #54 Android vertical, not on search.
- **#5 (Collections) and #6 (Ratings) are independent of each other**, but #6 needs its
  own design pass before it can be scoped the way #5 already can (the collections API
  contract already exists in `openapi.yaml`; ratings' doesn't).
- **#7 (Metrics/alerting) has no hard blocker** but is most valuable once #1 and #2/#3
  ship — it exists to catch regressions in exactly the features being added.
- **#9 (CD) blocks nothing else** on this list and is blocked by a decision, not code.

See [`ROADMAP.md`](ROADMAP.md) for how this ranking turns into concrete, sequenced
delivery waves.
