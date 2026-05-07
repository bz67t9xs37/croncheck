package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/example/croncheck/internal/monitor"
)

// SilenceHandler handles HTTP requests for managing job silences.
type SilenceHandler struct {
	store *monitor.SilenceStore
}

// NewSilenceHandler creates a new SilenceHandler backed by the given store.
func NewSilenceHandler(store *monitor.SilenceStore) *SilenceHandler {
	return &SilenceHandler{store: store}
}

func (h *SilenceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listSilences(w, r)
	case http.MethodPost:
		h.createSilence(w, r)
	case http.MethodDelete:
		h.deleteSilence(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *SilenceHandler) listSilences(w http.ResponseWriter, _ *http.Request) {
	silences := h.store.All()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(silences)
}

func (h *SilenceHandler) createSilence(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Job      string `json:"job"`
		Duration string `json:"duration"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Job == "" || body.Duration == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	d, err := time.ParseDuration(body.Duration)
	if err != nil {
		http.Error(w, "invalid duration format", http.StatusBadRequest)
		return
	}

	h.store.Silence(body.Job, time.Now().Add(d))
	w.WriteHeader(http.StatusCreated)
}

func (h *SilenceHandler) deleteSilence(w http.ResponseWriter, r *http.Request) {
	job := r.URL.Query().Get("job")
	if job == "" {
		http.Error(w, "missing job query parameter", http.StatusBadRequest)
		return
	}
	h.store.Unsilence(job)
	w.WriteHeader(http.StatusNoContent)
}
