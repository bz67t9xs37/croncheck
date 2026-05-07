package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/croncheck/internal/monitor"
)

func newTestSilenceStore() *monitor.SilenceStore {
	return monitor.NewSilenceStore()
}

func TestSilenceHandler_CreateSilence(t *testing.T) {
	store := newTestSilenceStore()
	h := NewSilenceHandler(store)

	body := map[string]interface{}{
		"job":      "backup",
		"duration": "2h",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/silences", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	if !store.IsSilenced("backup", time.Now()) {
		t.Error("expected job to be silenced")
	}
}

func TestSilenceHandler_ListSilences(t *testing.T) {
	store := newTestSilenceStore()
	store.Silence("backup", time.Now().Add(1*time.Hour))
	store.Silence("cleanup", time.Now().Add(2*time.Hour))

	h := NewSilenceHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/silences", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 silences, got %d", len(result))
	}
}

func TestSilenceHandler_DeleteSilence(t *testing.T) {
	store := newTestSilenceStore()
	store.Silence("backup", time.Now().Add(1*time.Hour))

	h := NewSilenceHandler(store)
	req := httptest.NewRequest(http.MethodDelete, "/silences?job=backup", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}

	if store.IsSilenced("backup", time.Now()) {
		t.Error("expected job to be unsilenced")
	}
}

func TestSilenceHandler_MethodNotAllowed(t *testing.T) {
	store := newTestSilenceStore()
	h := NewSilenceHandler(store)

	req := httptest.NewRequest(http.MethodPut, "/silences", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestSilenceHandler_CreateSilence_InvalidBody(t *testing.T) {
	store := newTestSilenceStore()
	h := NewSilenceHandler(store)

	req := httptest.NewRequest(http.MethodPost, "/silences", bytes.NewReader([]byte("not-json")))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
