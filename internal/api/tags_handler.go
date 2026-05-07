package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/example/croncheck/internal/monitor"
)

// TagsHandler handles CRUD operations for job tags.
type TagsHandler struct {
	store *monitor.TagStore
}

// NewTagsHandler creates a TagsHandler backed by the given TagStore.
func NewTagsHandler(store *monitor.TagStore) *TagsHandler {
	return &TagsHandler{store: store}
}

// ServeHTTP routes /api/tags and /api/tags/{jobID}.
func (h *TagsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Trim prefix and split path segments.
	path := strings.TrimPrefix(r.URL.Path, "/api/tags")
	path = strings.Trim(path, "/")

	switch r.Method {
	case http.MethodGet:
		if path == "" {
			h.listAll(w)
		} else {
			h.getJob(w, path)
		}
	case http.MethodPut:
		if path == "" {
			http.Error(w, "job ID required", http.StatusBadRequest)
			return
		}
		h.setTag(w, r, path)
	case http.MethodDelete:
		if path == "" {
			http.Error(w, "job ID and key required", http.StatusBadRequest)
			return
		}
		h.deleteTag(w, r, path)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TagsHandler) listAll(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.store.All())
}

func (h *TagsHandler) getJob(w http.ResponseWriter, jobID string) {
	w.Header().Set("Content-Type", "application/json")
	tags := h.store.Get(jobID)
	if tags == nil {
		tags = map[string]string{}
	}
	json.NewEncoder(w).Encode(tags)
}

func (h *TagsHandler) setTag(w http.ResponseWriter, r *http.Request, jobID string) {
	var body struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Key == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.store.Set(jobID, body.Key, body.Value); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TagsHandler) deleteTag(w http.ResponseWriter, r *http.Request, jobID string) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key query param required", http.StatusBadRequest)
		return
	}
	h.store.Delete(jobID, key)
	w.WriteHeader(http.StatusNoContent)
}
