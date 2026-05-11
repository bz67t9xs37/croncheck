package api

import (
	"encoding/json"
	"net/http"

	"github.com/example/croncheck/internal/monitor"
)

// MetricHandler serves job runtime metrics.
type MetricHandler struct {
	store *monitor.MetricStore
}

// NewMetricHandler returns a MetricHandler backed by the given store.
func NewMetricHandler(store *monitor.MetricStore) *MetricHandler {
	return &MetricHandler{store: store}
}

// ServeHTTP handles GET /metrics and GET /metrics?job=<id>.
func (h *MetricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	jobID := r.URL.Query().Get("job")
	if jobID != "" {
		snap, ok := h.store.Snapshot(jobID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "job not found"})
			return
		}
		_ = json.NewEncoder(w).Encode(snap)
		return
	}

	all := h.store.All()
	if all == nil {
		all = []monitor.MetricSnapshot{}
	}
	_ = json.NewEncoder(w).Encode(all)
}
