package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/croncheck/internal/monitor"
)

// SnapshotHandler serves job snapshot data over HTTP.
type SnapshotHandler struct {
	store *monitor.SnapshotStore
}

// NewSnapshotHandler creates a SnapshotHandler backed by the given store.
func NewSnapshotHandler(store *monitor.SnapshotStore) *SnapshotHandler {
	return &SnapshotHandler{store: store}
}

// ServeHTTP routes GET /snapshots and GET /snapshots?job=<id>.
func (h *SnapshotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	jobID := strings.TrimSpace(r.URL.Query().Get("job"))
	if jobID != "" {
		snap, ok := h.store.Get(jobID)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(snap)
		return
	}

	all := h.store.All()
	if all == nil {
		all = []monitor.JobSnapshot{}
	}
	json.NewEncoder(w).Encode(all)
}
