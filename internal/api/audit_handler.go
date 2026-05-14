package api

import (
	"encoding/json"
	"net/http"

	"github.com/example/croncheck/internal/monitor"
)

// AuditHandler exposes the audit log over HTTP.
type AuditHandler struct {
	store *monitor.AuditStore
}

// NewAuditHandler creates an AuditHandler backed by the given store.
func NewAuditHandler(store *monitor.AuditStore) *AuditHandler {
	return &AuditHandler{store: store}
}

// ServeHTTP handles GET /audit.
// An optional ?job= query parameter filters entries by job ID.
func (h *AuditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var entries []monitor.AuditEntry
	if jobID := r.URL.Query().Get("job"); jobID != "" {
		entries = h.store.ForJob(jobID)
	} else {
		entries = h.store.All()
	}

	if entries == nil {
		entries = []monitor.AuditEntry{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries) //nolint:errcheck
}
