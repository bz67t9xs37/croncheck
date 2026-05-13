package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/example/croncheck/internal/monitor"
)

// SLAHandler exposes CRUD for SLA policies and a read endpoint for violations.
type SLAHandler struct {
	store *monitor.SLAStore
}

// NewSLAHandler returns a new SLAHandler.
func NewSLAHandler(store *monitor.SLAStore) *SLAHandler {
	return &SLAHandler{store: store}
}

func (h *SLAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasSuffix(r.URL.Path, "/violations"):
		h.handleViolations(w, r)
	default:
		h.handlePolicy(w, r)
	}
}

func (h *SLAHandler) handlePolicy(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		jobID := r.URL.Query().Get("job_id")
		if jobID != "" {
			p, ok := h.store.Get(jobID)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			jsonResponse(w, p)
			return
		}
		jsonResponse(w, h.store.All())
	case http.MethodPost:
		var req struct {
			JobID          string  `json:"job_id"`
			MinSuccessRate float64 `json:"min_success_rate"`
			MaxDowntimeSec int     `json:"max_downtime_sec"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if req.JobID == "" {
			http.Error(w, "job_id required", http.StatusBadRequest)
			return
		}
		h.store.Set(monitor.SLAPolicy{
			JobID:          req.JobID,
			MinSuccessRate: req.MinSuccessRate,
			MaxDowntime:    time.Duration(req.MaxDowntimeSec) * time.Second,
		})
		w.WriteHeader(http.StatusCreated)
	case http.MethodDelete:
		jobID := r.URL.Query().Get("job_id")
		if !h.store.Delete(jobID) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *SLAHandler) handleViolations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	jobID := r.URL.Query().Get("job_id")
	jsonResponse(w, h.store.Violations(jobID))
}

func jsonResponse(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
