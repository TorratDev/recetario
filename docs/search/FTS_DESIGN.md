# Search & Indexing Design

Status: proposed, not yet implemented. Companion to [`BENCHMARKS.md`](BENCHMARKS.md), which
measures the current implementation this document proposes replacing.

## Current implementation (as of this doc)

`backend/internal/repositories/search_repository.go` — `SearchRecipes` builds a fully
parameterized `WHERE` clause and matches `title`/`description`/`ingredients.name` with
`LIKE '%term%'` (leading wildcard, so no index can serve it — every call is a full table
scan of `recipes`, plus a correlated `EXISTS` subquery scan of `ingredients` per row).
`schema.sql` defines **no indexes at all** beyond the implicit ones from `PRIMARY KEY`.
Exact-match filters (`difficulty`, `category`, `cuisine`, tag `EXISTS`) also scan.

This is fine at seed-data scale (dozens of rows) and is not a bug — it just hasn't needed
to be fast yet. `BENCHMARKS.md` quantifies exactly where it stops being fine.

## Why not `EXTERNAL CONTENT` FTS5

SQLite FTS5's `content=` external-content mode ties each FTS row to the source table's
`rowid`. Every table in this schema uses a `TEXT PRIMARY KEY` (UUID strings) —
`recipes.id`, `ingredients.id`, etc. — so there is no integer rowid to key against, and
`content=recipes` mode isn't usable without either (a) adding a shadow integer rowid
column to `recipes`, which ripples into every repository query and the API's `id` type
contract, or (b) using `WITHOUT ROWID`, which FTS5 doesn't support as a content table.
Neither is worth the disruption.

## Proposed approach: standalone FTS5 table, `id` as a regular column

Use a **standalone** (not external-content) FTS5 table that stores the UUID as an
ordinary `UNINDEXED` column instead of relying on rowid mapping:

```sql
CREATE VIRTUAL TABLE IF NOT EXISTS recipe_search USING fts5(
    id UNINDEXED,
    title,
    description,
    ingredients_text,
    tokenize = 'porter unicode61'
);
```

FTS5 still allocates its own internal integer rowid for this table; the app never
touches it and always joins/filters on the `id` column instead:

```sql
SELECT r.id, r.title, ...
FROM recipes r
JOIN recipe_search s ON s.id = r.id
WHERE recipe_search MATCH ?
ORDER BY bm25(recipe_search)
LIMIT ? OFFSET ?;
```

### Keeping it in sync

Standalone FTS5 tables don't auto-sync — the app is responsible for writes. Two options:

1. **Triggers** on `recipes` and `ingredients` (`AFTER INSERT/UPDATE/DELETE`) that
   `DELETE FROM recipe_search WHERE id = ...` then re-`INSERT` a recomputed row
   (`ingredients_text` is an aggregate of all of a recipe's ingredient names, so any
   ingredient change must recompute the whole row — FTS5 has no partial-row update).
   Triggers keep every write path (including any future direct-SQL scripts) correct for
   free, at the cost of a few extra statements per write.
2. **Application-level sync** inside `RecipeRepository.CreateRecipe`/`UpdateRecipe`/
   `DeleteRecipe`, mirroring the existing pattern where ingredients/instructions/tags are
   already written alongside the recipe row.

Recommend **triggers**: this schema has no service layer enforcing "always go through
the repository," and a trigger can't be bypassed by a forgotten call site. Both the
virtual table and its triggers are `CREATE ... IF NOT EXISTS` / idempotent, so they slot
into `schema.sql` the same way every other table does today — no migration tool needed,
consistent with [`CLAUDE.md`](../../CLAUDE.md)'s "single source of truth, applied at
startup" model.

### Query-time concerns

- **Sanitize user input before `MATCH`.** FTS5 query syntax treats `-`, `"`, `*`, `AND`/
  `OR`/`NOT`, etc. as operators. A raw user query like `chicken - rice` is a valid FTS5
  expression (NOT), not a literal phrase, and unbalanced quotes are a syntax error, not a
  zero-result match. Wrap each whitespace-split term in double quotes and join with a
  bareword `AND`/`OR` chosen by the caller — never interpolate the raw string.
- **Ranking**: order by `bm25(recipe_search)` (ascending — lower is more relevant),
  optionally weighted per-column (`bm25(recipe_search, 3.0, 1.0, 1.0)` to weight title
  above description/ingredients).
- `SuggestTitles` (prefix autocomplete) stays on `LIKE 'term%'` — that's a *left-anchored*
  LIKE, which a plain `CREATE INDEX idx_recipes_title ON recipes(title)` already serves
  via SQLite's LIKE-optimization, so it doesn't need FTS5 at all.

## Non-FTS indexes to add regardless

These help the exact-match and sort paths independent of the FTS5 decision:

```sql
CREATE INDEX IF NOT EXISTS idx_recipes_difficulty ON recipes(difficulty);
CREATE INDEX IF NOT EXISTS idx_recipes_category    ON recipes(category);
CREATE INDEX IF NOT EXISTS idx_recipes_cuisine     ON recipes(cuisine);
CREATE INDEX IF NOT EXISTS idx_recipes_created_at  ON recipes(created_at);
CREATE INDEX IF NOT EXISTS idx_recipes_title       ON recipes(title);
CREATE INDEX IF NOT EXISTS idx_ingredients_recipe_id  ON ingredients(recipe_id);
CREATE INDEX IF NOT EXISTS idx_instructions_recipe_id ON instructions(recipe_id);
```

`recipe_tags` and `collection_recipes` already have covering composite primary keys for
their `EXISTS`/join patterns; no additional index needed there.

## Migration/rollout plan

1. Add the indexes above and the `recipe_search` table + triggers to `schema.sql`
   (additive, `IF NOT EXISTS` — safe on every existing `recipeapp.db`, since the schema
   is re-applied on every startup per `CLAUDE.md`).
2. Backfill: on first run after deploy, existing rows in `recipes` predate the triggers,
   so a one-time `INSERT INTO recipe_search SELECT id, title, description, (SELECT
   group_concat(name, ' ') FROM ingredients WHERE recipe_id = recipes.id) FROM recipes`
   is needed. Guard it (e.g. `INSERT ... WHERE NOT EXISTS (SELECT 1 FROM recipe_search
   LIMIT 1)`) so it only runs once.
3. Swap `buildRecipeSearchWhere`'s `LIKE` clause for the FTS5 `MATCH` join when
   `f.Query != ""`; leave every other filter (`difficulty`, `category`, tags, etc.)
   exactly as-is — they're unrelated to text search and already correctly parameterized.
4. Compare against the `BENCHMARKS.md` baseline at the same data volumes to confirm the
   change actually helped before removing the old code path.

## Out of scope for this doc

Ranking tuning, synonym/stemming language selection beyond the default `porter
unicode61` tokenizer, and multi-language search are left for a follow-up once real usage
data exists — no point over-designing ranking weights against seed data.
