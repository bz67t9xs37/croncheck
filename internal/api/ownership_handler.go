package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/croncheck/internal/monitor"
)

// OwnershipHandler handles HTTP requests for job ownership records.
type OwnershipHandler struct {
	store *monitor.OwnershipStore
}

// NewOwnershipHandler creates a new OwnershipHandler.
func NewOwnershipHandler(store *monitor.OwnershipStore) *OwnershipHandler {
	return &OwnershipHandler{store: store}
}

func (h *OwnershipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPut:
		h.handlePut(w, r)
	case http.MethodDelete:
		h.handleDelete(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *OwnershipHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	job := strings.TrimSpace(r.URL.Query().Get("job"))
	if job == "" {
		all := h.store.All()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(all)
		return
	}
	owner, ok := h.store.Get(job)
	if !ok {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(owner)
}

func (h *OwnershipHandler) handlePut(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Job   string `json:"job"`
		Email string `json:"email"`
		Team  string `json:"team"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.store.Set(req.Job, req.Email, req.Team); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *OwnershipHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	job := strings.TrimSpace(r.URL.Query().Get("job"))
	if job == "" {
		http.Error(w, "job parameter required", http.StatusBadRequest)
		return
	}
	h.store.Delete(job)
	w.WriteHeader(http.StatusNoContent)
}
