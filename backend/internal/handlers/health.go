package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// HealthHandler reports basic liveness/readiness for load balancers, container
// orchestrators, and uptime checks. It has no auth requirement — health
// checks run before a caller could plausibly have a session.
type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HandleHealth returns 200 when the process is up and can reach its
// database, 503 otherwise. Docker/Kubernetes and uptime monitors alike treat
// a non-2xx here as "not ready" — the SQLite ping is what actually catches a
// broken DB_PATH mount or a corrupted file, versus just "process exists."
func (h *HealthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := h.db.PingContext(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "error",
			"error":  "database unreachable",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
