package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/user/croncheck/internal/monitor"
)

// HistoryHandler serves check-in history via HTTP.
type HistoryHandler struct {
	history *monitor.History
}

// NewHistoryHandler creates a HistoryHandler backed by the given History store.
func NewHistoryHandler(h *monitor.History) *HistoryHandler {
	return &HistoryHandler{history: h}
}

// ServeHTTP routes GET /history and GET /history/{job} requests.
func (hh *HistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Strip leading slash and split: "/history" or "/history/jobname"
	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/history"), "/", 2)
	jobName := ""
	if len(parts) == 2 {
		jobName = strings.TrimPrefix(parts[1], "/")
	}

	var records []monitor.CheckInRecord
	if jobName != "" {
		records = hh.history.Get(jobName)
	} else {
		records = hh.history.All()
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(records); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
