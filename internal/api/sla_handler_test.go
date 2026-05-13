package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/croncheck/internal/monitor"
)

func newTestSLAStore() *monitor.SLAStore {
	return monitor.NewSLAStore(0)
}

func TestSLAHandler_CreateAndGet(t *testing.T) {
	h := NewSLAHandler(newTestSLAStore())

	body := `{"job_id":"job1","min_success_rate":0.95,"max_downtime_sec":300}`
	req := httptest.NewRequest(http.MethodPost, "/sla", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/sla?job_id=job1", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
	var p monitor.SLAPolicy
	if err := json.NewDecoder(w2.Body).Decode(&p); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if p.MinSuccessRate != 0.95 {
		t.Errorf("expected 0.95, got %f", p.MinSuccessRate)
	}
}

func TestSLAHandler_ListAll(t *testing.T) {
	store := newTestSLAStore()
	store.Set(monitor.SLAPolicy{JobID: "job1"})
	store.Set(monitor.SLAPolicy{JobID: "job2"})
	h := NewSLAHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/sla", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	var policies []monitor.SLAPolicy
	if err := json.NewDecoder(w.Body).Decode(&policies); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(policies))
	}
}

func TestSLAHandler_Delete(t *testing.T) {
	store := newTestSLAStore()
	store.Set(monitor.SLAPolicy{JobID: "job1"})
	h := NewSLAHandler(store)

	req := httptest.NewRequest(http.MethodDelete, "/sla?job_id=job1", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}

	req2 := httptest.NewRequest(http.MethodDelete, "/sla?job_id=job1", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, req2)
	if w2.Code != http.StatusNotFound {
		t.Errorf("expected 404 on second delete, got %d", w2.Code)
	}
}

func TestSLAHandler_Violations(t *testing.T) {
	store := newTestSLAStore()
	store.RecordViolation("job1", "rate too low")
	store.RecordViolation("job2", "downtime exceeded")
	h := NewSLAHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/sla/violations?job_id=job1", nil)
	req.URL.Path = "/sla/violations"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	var violations []monitor.SLAViolation
	if err := json.NewDecoder(w.Body).Decode(&violations); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(violations) != 1 {
		t.Errorf("expected 1 violation for job1, got %d", len(violations))
	}
}

func TestSLAHandler_MethodNotAllowed(t *testing.T) {
	h := NewSLAHandler(newTestSLAStore())
	req := httptest.NewRequest(http.MethodPatch, "/sla", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestSLAHandler_Create_MissingJobID(t *testing.T) {
	h := NewSLAHandler(newTestSLAStore())
	body := `{"min_success_rate":0.9}`
	req := httptest.NewRequest(http.MethodPost, "/sla", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
