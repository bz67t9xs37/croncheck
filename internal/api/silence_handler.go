package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/user/croncheck/internal/monitor"
)

// SilenceHandler handles HTTP requests for managing job silences.
type SilenceHandler struct {
	store *monitor.SilenceStore
}

// NewSilenceHandler creates a new SilenceHandler.
func NewSilenceHandler(store *monitor.SilenceStore) *SilenceHandler {
	return &SilenceHandler{store: store}
}

type silenceRequest struct {
	JobName  string `json:"job_name"`
	Duration string `json:"duration"`
	Reason   string `json:"reason"`
}

// ServeHTTP routes GET (list all), POST (add silence), DELETE (remove silence).
func (h *SilenceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listSilences(w, r)
	case http.MethodPost:
		h.addSilence(w, r)
	case http.MethodDelete:
		h.removeSilence(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *SilenceHandler) listSilences(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.store.All())
}

func (h *SilenceHandler) addSilence(w http.ResponseWriter, r *http.Request) {
	var req silenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.JobName == "" {
		http.Error(w, "job_name is required", http.StatusBadRequest)
		return
	}
	dur, err := time.ParseDuration(req.Duration)
	if err != nil || dur <= 0 {
		http.Error(w, "invalid duration", http.StatusBadRequest)
		return
	}
	h.store.Silence(req.JobName, time.Now().Add(dur), req.Reason)
	w.WriteHeader(http.StatusNoContent)
}

func (h *SilenceHandler) removeSilence(w http.ResponseWriter, r *http.Request) {
	jobName := r.URL.Query().Get("job")
	if jobName == "" {
		http.Error(w, "job query param required", http.StatusBadRequest)
		return
	}
	h.store.Unsilence(jobName)
	w.WriteHeader(http.StatusNoContent)
}
