# syntax=docker/dockerfile:1

# --- build stage -------------------------------------------------------
FROM golang:1.25-alpine AS build
WORKDIR /src

# Cache module downloads separately from source changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0: the app uses modernc.org/sqlite (pure Go, no CGO) precisely
# so it can produce a static binary that runs on a scratch/distroless base
# with no libc — keep that property here.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/recipe-server ./cmd

# --- runtime stage -------------------------------------------------------
# distroless: no shell, no package manager, minimal attack surface. "nonroot"
# variant already runs as an unprivileged UID (65532) by default.
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=build /out/recipe-server ./recipe-server
# Templates and static assets are read from disk at request time
# (html/template.ParseFiles / http.Dir), not go:embed — they must ship
# alongside the binary, not just at build time.
COPY --from=build /src/web ./web

# recipeapp.db (DB_PATH) and uploads/ (UPLOAD_DIR) are written at runtime —
# mount volumes at these paths from docker-compose.yml to persist them.
VOLUME ["/app/data", "/app/uploads"]

ENV DB_PATH=/app/data/recipeapp.db
ENV UPLOAD_DIR=/app/uploads
ENV PORT=8080

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["./recipe-server"]
