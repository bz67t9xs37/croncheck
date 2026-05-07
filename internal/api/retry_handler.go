package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/user/croncheck/internal/monitor"
)

// RetryHandler exposes retry state over HTTP.
type RetryHandler struct {
	tracker *monitor.RetryTracker
}

// NewRetryHandler creates a RetryHandler.
func NewRetryHandler(tracker *monitor.RetryTracker) *RetryHandler {
	return &RetryHandler{tracker: tracker}
}

// ServeHTTP handles GET /retries and DELETE /retries/{job}.
func (h *RetryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listRetries(w, r)
	case http.MethodDelete:
		h.resetRetry(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RetryHandler) listRetries(w http.ResponseWriter, r *http.Request) {
	jobName := r.URL.Query().Get("job")
	all := h.tracker.All()

	if jobName != "" {
		filtered := all[:0]
		for _, s := range all {
			if s.JobName == jobName {
				filtered = append(filtered, s)
			}
		}
		all = filtered
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(all)
}

func (h *RetryHandler) resetRetry(w http.ResponseWriter, r *http.Request) {
	// Expect path: /retries/{jobName}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 || parts[len(parts)-1] == "" {
		http.Error(w, "job name required", http.StatusBadRequest)
		return
	}
	jobName := parts[len(parts)-1]
	h.tracker.Reset(jobName)
	w.WriteHeader(http.StatusNoContent)
}
