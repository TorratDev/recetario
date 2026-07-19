package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"recipe-app/internal/database"
)

func TestHealthHandler_HandleHealth(t *testing.T) {
	dir := t.TempDir()
	db, err := sql.Open("sqlite", filepath.Join(dir, "health.db"))
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	if err := database.ApplySchema(db); err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	handler := NewHealthHandler(db)

	t.Run("database reachable", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		w := httptest.NewRecorder()

		handler.HandleHealth(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 (body %q)", w.Code, w.Body.String())
		}
		if ct := w.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
	})

	t.Run("database unreachable", func(t *testing.T) {
		db.Close() // last sub-test: force PingContext to fail by closing the pool

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		w := httptest.NewRecorder()

		handler.HandleHealth(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Fatalf("status = %d, want 503 (body %q)", w.Code, w.Body.String())
		}
	})
}
