package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/yourorg/croncheck/internal/monitor"
)

// RunbookHandler serves CRUD operations for job runbook entries.
type RunbookHandler struct {
	store *monitor.RunbookStore
}

// NewRunbookHandler creates a RunbookHandler backed by the given store.
func NewRunbookHandler(store *monitor.RunbookStore) *RunbookHandler {
	return &RunbookHandler{store: store}
}

func (h *RunbookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	case http.MethodDelete:
		h.delete(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RunbookHandler) list(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job_id")
	if jobID != "" {
		entry, ok := h.store.Get(jobID)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entry)
		return
	}
	all := h.store.All()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(all)
}

func (h *RunbookHandler) create(w http.ResponseWriter, r *http.Request) {
	var entry monitor.RunbookEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if err := h.store.Set(entry); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *RunbookHandler) delete(w http.ResponseWriter, r *http.Request) {
	jobID := strings.TrimSpace(r.URL.Query().Get("job_id"))
	if jobID == "" {
		http.Error(w, "job_id is required", http.StatusBadRequest)
		return
	}
	if !h.store.Delete(jobID) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
