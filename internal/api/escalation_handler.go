package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/example/croncheck/internal/monitor"
)

// EscalationHandler handles HTTP requests for escalation policy management.
type EscalationHandler struct {
	store *monitor.EscalationStore
}

// NewEscalationHandler creates a new EscalationHandler.
func NewEscalationHandler(store *monitor.EscalationStore) *EscalationHandler {
	return &EscalationHandler{store: store}
}

type escalationRequest struct {
	JobID     string `json:"job_id"`
	Threshold int    `json:"threshold"`
	IntervalS int    `json:"interval_seconds"`
	Webhook   string `json:"webhook"`
}

// ServeHTTP dispatches to create/list or delete based on method and path.
func (h *EscalationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		h.list(w, r)
	case r.Method == http.MethodPost:
		h.create(w, r)
	case r.Method == http.MethodDelete:
		h.delete(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *EscalationHandler) list(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job_id")
	policies := h.store.All()
	if jobID != "" {
		filtered := policies[:0]
		for _, p := range policies {
			if p.JobID == jobID {
				filtered = append(filtered, p)
			}
		}
		policies = filtered
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policies)
}

func (h *EscalationHandler) create(w http.ResponseWriter, r *http.Request) {
	var req escalationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.JobID == "" || req.Threshold <= 0 {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	policy := monitor.EscalationPolicy{
		JobID:     req.JobID,
		Threshold: req.Threshold,
		Interval:  time.Duration(req.IntervalS) * time.Second,
		Webhook:   req.Webhook,
	}
	h.store.Set(policy)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(policy)
}

func (h *EscalationHandler) delete(w http.ResponseWriter, r *http.Request) {
	jobID := strings.TrimPrefix(r.URL.Path, "/escalation/")
	if jobID == "" {
		http.Error(w, "job_id required", http.StatusBadRequest)
		return
	}
	if !h.store.Delete(jobID) {
		http.Error(w, "policy not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
