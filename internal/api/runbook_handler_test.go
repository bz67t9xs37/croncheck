package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourorg/croncheck/internal/monitor"
)

func newTestRunbookStore() *monitor.RunbookStore {
	return monitor.NewRunbookStore()
}

func TestRunbookHandler_CreateAndGet(t *testing.T) {
	store := newTestRunbookStore()
	h := NewRunbookHandler(store)

	body, _ := json.Marshal(monitor.RunbookEntry{
		JobID: "deploy",
		URL:   "https://wiki.example.com/deploy",
		Notes: "Check pod restarts",
	})
	req := httptest.NewRequest(http.MethodPost, "/runbooks", bytes.NewReader(body))
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rw.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/runbooks?job_id=deploy", nil)
	rw2 := httptest.NewRecorder()
	h.ServeHTTP(rw2, req2)
	if rw2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw2.Code)
	}
	var got monitor.RunbookEntry
	json.NewDecoder(rw2.Body).Decode(&got)
	if got.URL != "https://wiki.example.com/deploy" {
		t.Errorf("unexpected URL: %s", got.URL)
	}
}

func TestRunbookHandler_ListAll(t *testing.T) {
	store := newTestRunbookStore()
	_ = store.Set(monitor.RunbookEntry{JobID: "j1", URL: "https://example.com/j1"})
	_ = store.Set(monitor.RunbookEntry{JobID: "j2", URL: "https://example.com/j2"})
	h := NewRunbookHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/runbooks", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var entries []monitor.RunbookEntry
	json.NewDecoder(rw.Body).Decode(&entries)
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestRunbookHandler_Delete(t *testing.T) {
	store := newTestRunbookStore()
	_ = store.Set(monitor.RunbookEntry{JobID: "old", URL: "https://example.com/old"})
	h := NewRunbookHandler(store)

	req := httptest.NewRequest(http.MethodDelete, "/runbooks?job_id=old", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rw.Code)
	}
}

func TestRunbookHandler_Delete_NotFound(t *testing.T) {
	h := NewRunbookHandler(newTestRunbookStore())
	req := httptest.NewRequest(http.MethodDelete, "/runbooks?job_id=ghost", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rw.Code)
	}
}

func TestRunbookHandler_MethodNotAllowed(t *testing.T) {
	h := NewRunbookHandler(newTestRunbookStore())
	req := httptest.NewRequest(http.MethodPut, "/runbooks", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}
}
