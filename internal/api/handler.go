package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/user/croncheck/internal/monitor"
)

// Handler provides HTTP endpoints for job check-ins and status queries.
type Handler struct {
	registry *monitor.Registry
}

// NewHandler creates a new Handler backed by the given registry.
func NewHandler(r *monitor.Registry) *Handler {
	return &Handler{registry: r}
}

// RegisterRoutes attaches the handler routes to the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/checkin/", h.handleCheckIn)
	mux.HandleFunc("/status", h.handleStatus)
}

// handleCheckIn processes a POST /checkin/{jobName} request.
func (h *Handler) handleCheckIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	jobName := r.URL.Path[len("/checkin/"):]
	if jobName == "" {
		http.Error(w, "job name required", http.StatusBadRequest)
		return
	}
	if err := h.registry.CheckIn(jobName, time.Now()); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ok":true}`))
}

// statusResponse is the JSON shape returned by /status.
type statusResponse struct {
	Jobs []jobStatusEntry `json:"jobs"`
}

type jobStatusEntry struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	LastSeen  time.Time `json:"last_seen,omitempty"`
}

// handleStatus returns a JSON summary of all monitored jobs.
func (h *Handler) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	snapshot := h.registry.Snapshot()
	resp := statusResponse{Jobs: make([]jobStatusEntry, 0, len(snapshot))}
	for _, job := range snapshot {
		resp.Jobs = append(resp.Jobs, jobStatusEntry{
			Name:     job.Name,
			Status:   job.Status.String(),
			LastSeen: job.LastCheckIn,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
