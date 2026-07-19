# Search Benchmarks & SLO Baseline

Real numbers from the current (LIKE-based, unindexed) search implementation, captured to
give [`FTS_DESIGN.md`](FTS_DESIGN.md) something concrete to compare against once
implemented. Source: `backend/internal/repositories/search_bench_test.go`, run via
`make bench-search`.

## Methodology

- Each benchmark seeds a throwaway on-disk SQLite database (`b.TempDir()`, real schema
  via `database.ApplySchema`) with N synthetic recipes through `RecipeRepository.CreateRecipe`
  — the same code path production writes use, each with 5 randomly-drawn ingredients — so
  seeded data exercises the same `ingredients` join real recipes would.
- Setup runs before `b.ResetTimer()`; only the repository call itself is timed.
- Three query shapes, each at N = 100 / 1,000 / 5,000 recipes:
  - **TextQuery** — `RecipeFilter{Query: "Chicken"}`, i.e. the `title/description LIKE`
    OR `EXISTS (... ingredients ...)` path.
  - **FilterOnly** — `RecipeFilter{Difficulty: "medium", Cuisine: "Italian"}`, i.e. plain
    equality `WHERE` clauses, no text search.
  - **SuggestTitles** — the autocomplete endpoint's left-anchored `title LIKE 'term%'`.
- Run with `-benchtime=1x` (single real invocation per size, not an averaged loop) — the
  goal here is an order-of-magnitude baseline, not noise-free micro-benchmarking.

## Environment this run was captured on

AMD Ryzen 5 5600X 6-Core, Windows/amd64, Go 1.25.6, `modernc.org/sqlite` (pure-Go, no
CGO). Absolute numbers will differ on other hardware; the *shape* of the scaling is the
part that matters.

## Results

| Benchmark | N=100 | N=1,000 | N=5,000 |
|---|---|---|---|
| `TextQuery` (LIKE + ingredient EXISTS) | 7.2 ms | 687 ms | **35.9 s** |
| `FilterOnly` (equality WHERE, no text) | 0.38 ms | 1.0 ms | 2.7 ms |
| `SuggestTitles` (prefix LIKE) | 0.19 ms | 0.42 ms | 1.9 ms |

## Reading this

`FilterOnly` and `SuggestTitles` scale roughly linearly with table size — as expected for
a full-table scan of just `recipes`, which stays cheap even unindexed at these sizes.

`TextQuery` does not scale linearly — it's **~95x slower at 10x the rows** (100→1,000)
and **~52x slower at 5x the rows** (1,000→5,000). That's quadratic-shaped, not
coincidental noise. The cause is `buildRecipeSearchWhere`'s ingredient clause:

```sql
EXISTS (SELECT 1 FROM ingredients i WHERE i.recipe_id = recipes.id AND i.name LIKE ?)
```

`ingredients` has no index on `recipe_id` (see `schema.sql`), so this correlated
subquery does a full scan of the *entire* `ingredients` table for every candidate row in
`recipes` — at N=5,000 recipes × 5 ingredients each, that's ~5,000 × 25,000 = 125M row
comparisons for a single search. This is exactly the query plan `FTS_DESIGN.md`'s
`idx_ingredients_recipe_id` index (plus the FTS5 migration for the text-match itself) is
meant to fix — expect that change to bring `TextQuery` down to roughly the same shape as
`FilterOnly`.

## Proposed SLO baseline

There's no production traffic yet to derive a real SLO from, so this is a starting
target, not a measured commitment — revisit once `FTS_DESIGN.md` ships and real usage
exists:

- **p95 search latency: < 200 ms** at the current seed-data scale (dozens of recipes) —
  trivially met by `FilterOnly`/`SuggestTitles` today, and by `TextQuery` up to roughly
  N≈1,000 based on the table above.
- **`TextQuery` must not regress super-linearly with N.** The table above is the
  documented "before" state specifically because it violates this — any future change to
  the search path should be benchmarked against these numbers before merging, using
  `make bench-search`.
- Alert threshold, once `/api/search` has real observability (see
  [`../ops/RELEASE_AND_OBSERVABILITY.md`](../ops/RELEASE_AND_OBSERVABILITY.md)): p95 >
  500 ms sustained over 5 minutes.

## Reproducing

```bash
make bench-search
```

Runs `go test ./internal/repositories/... -bench . -benchtime=1x -benchmem` from
`backend/`. Re-run and update the table above whenever the search implementation
changes materially (new indexes, FTS5 migration, etc.) — this file is a point-in-time
snapshot, not a live dashboard.
