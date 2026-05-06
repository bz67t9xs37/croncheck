package api

import (
	"encoding/json"
	"net/http"

	"github.com/user/croncheck/internal/monitor"
)

// AlertLogHandler serves recorded alert entries over HTTP.
type AlertLogHandler struct {
	log *monitor.AlertLog
}

// NewAlertLogHandler creates a new AlertLogHandler.
func NewAlertLogHandler(log *monitor.AlertLog) *AlertLogHandler {
	return &AlertLogHandler{log: log}
}

// ServeHTTP handles GET /alerts[?job=<name>].
// Returns a JSON array of alert entries, optionally filtered by job name.
func (h *AlertLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var entries []monitor.AlertEntry
	if job := r.URL.Query().Get("job"); job != "" {
		entries = h.log.ForJob(job)
	} else {
		entries = h.log.All()
	}

	if entries == nil {
		entries = []monitor.AlertEntry{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
