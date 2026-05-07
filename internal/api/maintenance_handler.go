package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/croncheck/internal/monitor"
)

// MaintenanceHandler handles maintenance window API requests.
type MaintenanceHandler struct {
	store *monitor.MaintenanceStore
}

// NewMaintenanceHandler creates a new MaintenanceHandler.
func NewMaintenanceHandler(store *monitor.MaintenanceStore) *MaintenanceHandler {
	return &MaintenanceHandler{store: store}
}

type createMaintenanceRequest struct {
	JobName string    `json:"job_name"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
}

// ServeHTTP routes maintenance requests.
func (h *MaintenanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listWindows(w, r)
	case http.MethodPost:
		h.createWindow(w, r)
	case http.MethodDelete:
		h.deleteWindows(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *MaintenanceHandler) listWindows(w http.ResponseWriter, r *http.Request) {
	job := r.URL.Query().Get("job")
	all := h.store.All()
	if job != "" {
		filtered := all[:0]
		for _, win := range all {
			if win.JobName == job {
				filtered = append(filtered, win)
			}
		}
		all = filtered
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(all)
}

func (h *MaintenanceHandler) createWindow(w http.ResponseWriter, r *http.Request) {
	var req createMaintenanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.JobName == "" || req.End.IsZero() || !req.End.After(req.Start) {
		http.Error(w, "job_name, start, and a valid end time are required", http.StatusBadRequest)
		return
	}
	win := h.store.Add(req.JobName, req.Start, req.End)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(win)
}

func (h *MaintenanceHandler) deleteWindows(w http.ResponseWriter, r *http.Request) {
	job := r.URL.Query().Get("job")
	if job == "" {
		http.Error(w, "job query parameter required", http.StatusBadRequest)
		return
	}
	count := h.store.Remove(job)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"removed": count})
}
