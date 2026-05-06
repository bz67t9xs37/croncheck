package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/user/croncheck/internal/monitor"
)

// SchedulerStatusResponse is returned by the scheduler status endpoint.
type SchedulerStatusResponse struct {
	Running   bool      `json:"running"`
	CheckedAt time.Time `json:"checked_at"`
	JobCount  int       `json:"job_count"`
}

// SchedulerHandler exposes a read-only status endpoint for the scheduler.
type SchedulerHandler struct {
	registry *monitor.Registry
	running  func() bool
}

// NewSchedulerHandler creates a handler that reports scheduler status.
func NewSchedulerHandler(registry *monitor.Registry, running func() bool) *SchedulerHandler {
	return &SchedulerHandler{registry: registry, running: running}
}

// ServeHTTP handles GET /scheduler/status.
func (h *SchedulerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := SchedulerStatusResponse{
		Running:   h.running(),
		CheckedAt: time.Now().UTC(),
		JobCount:  len(h.registry.All()),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
