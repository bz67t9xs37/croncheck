package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/example/croncheck/internal/monitor"
)

// RateLimitHandler exposes rate limit state for inspection and manual reset.
type RateLimitHandler struct {
	store *monitor.RateLimitStore
}

// NewRateLimitHandler creates a new RateLimitHandler.
func NewRateLimitHandler(store *monitor.RateLimitStore) *RateLimitHandler {
	return &RateLimitHandler{store: store}
}

type rateLimitRecord struct {
	JobID     string    `json:"job_id"`
	LastAlert time.Time `json:"last_alert"`
}

// ServeHTTP routes GET (list all) and DELETE (reset by job) requests.
func (h *RateLimitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listAll(w, r)
	case http.MethodDelete:
		h.reset(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RateLimitHandler) listAll(w http.ResponseWriter, _ *http.Request) {
	all := h.store.All()
	records := make([]rateLimitRecord, 0, len(all))
	for jobID, last := range all {
		records = append(records, rateLimitRecord{JobID: jobID, LastAlert: last})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func (h *RateLimitHandler) reset(w http.ResponseWriter, r *http.Request) {
	// Expect path: /ratelimit/{jobID}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 || parts[len(parts)-1] == "" {
		http.Error(w, "job ID required", http.StatusBadRequest)
		return
	}
	jobID := parts[len(parts)-1]
	h.store.Reset(jobID)
	w.WriteHeader(http.StatusNoContent)
}
