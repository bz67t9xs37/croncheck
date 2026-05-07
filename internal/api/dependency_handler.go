package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/croncheck/internal/monitor"
)

// DependencyHandler exposes CRUD endpoints for job dependency links.
type DependencyHandler struct {
	store *monitor.DependencyStore
}

// NewDependencyHandler creates a DependencyHandler backed by store.
func NewDependencyHandler(store *monitor.DependencyStore) *DependencyHandler {
	return &DependencyHandler{store: store}
}

// ServeHTTP routes requests to the appropriate sub-handler.
//
//	POST   /dependencies          – create a link
//	GET    /dependencies          – list all links
//	DELETE /dependencies          – remove a link
func (h *DependencyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCreate(w, r)
	case http.MethodGet:
		h.handleList(w, r)
	case http.MethodDelete:
		h.handleDelete(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

type dependencyRequest struct {
	Upstream   string `json:"upstream"`
	Downstream string `json:"downstream"`
}

func (h *DependencyHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req dependencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Upstream == "" || req.Downstream == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.store.Add(req.Upstream, req.Downstream); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *DependencyHandler) handleList(w http.ResponseWriter, r *http.Request) {
	job := strings.TrimSpace(r.URL.Query().Get("job"))
	var result interface{}
	if job != "" {
		result = map[string]interface{}{
			"upstreams":   h.store.UpstreamsOf(job),
			"downstreams": h.store.DownstreamsOf(job),
		}
	} else {
		result = h.store.All()
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *DependencyHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	var req dependencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Upstream == "" || req.Downstream == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.store.Remove(req.Upstream, req.Downstream); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
