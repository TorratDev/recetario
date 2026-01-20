# 🍲 Recipe App

A cross-platform recipe management application with Go backend, HTMX web frontend, and Android mobile client.

## Project Structure

```
recipe-app/
├── backend/          # Go API server
│   ├── cmd/         # Application entry points
│   ├── internal/    # Private application code
│   ├── pkg/         # Public library code
│   ├── migrations/  # Database migrations
│   └── configs/     # Configuration files
├── web/             # HTMX frontend
└── mobile/          # Android application
```

## Quick Start

### Backend

```bash
cd backend
go run cmd/main.go
```

The server will start on `:8080` with:

- API endpoints at `/api/*`
- Web interface at `/`

### Database Setup

1. Install PostgreSQL
2. Create database: `createdb recipe_app`
3. Run migrations: `psql recipe_app < migrations/001_initial_schema.sql`

## API Endpoints

- `GET /api/recipes` - List all recipes
- `POST /api/recipes` - Create new recipe
- `GET /api/recipes/{id}` - Get specific recipe
- `PUT /api/recipes/{id}` - Update recipe
- `DELETE /api/recipes/{id}` - Delete recipe

## Tech Stack

- **Backend**: Go, PostgreSQL, Gin
- **Web**: HTMX, Tailwind CSS
- **Mobile**: Android (Kotlin)
- **Database**: PostgreSQL with full-text search
