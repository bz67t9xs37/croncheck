package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/example/croncheck/internal/monitor"
)

// AnnotationHandler serves GET/POST/DELETE for job annotations.
type AnnotationHandler struct {
	store *monitor.AnnotationStore
}

// NewAnnotationHandler creates a new AnnotationHandler.
func NewAnnotationHandler(store *monitor.AnnotationStore) *AnnotationHandler {
	return &AnnotationHandler{store: store}
}

func (h *AnnotationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (h *AnnotationHandler) list(w http.ResponseWriter, r *http.Request) {
	jobID := strings.TrimSpace(r.URL.Query().Get("job_id"))
	var result interface{}
	if jobID != "" {
		result = h.store.Get(jobID)
	} else {
		result = h.store.All()
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func (h *AnnotationHandler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		JobID   string `json:"job_id"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if err := h.store.Add(req.JobID, req.Message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *AnnotationHandler) delete(w http.ResponseWriter, r *http.Request) {
	jobID := strings.TrimSpace(r.URL.Query().Get("job_id"))
	if jobID == "" {
		http.Error(w, "job_id required", http.StatusBadRequest)
		return
	}
	h.store.Delete(jobID)
	w.WriteHeader(http.StatusNoContent)
}
