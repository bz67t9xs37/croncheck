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

func newTestEscalationStore() *monitor.EscalationStore {
	return monitor.NewEscalationStore()
}

func TestEscalationHandler_CreateAndList(t *testing.T) {
	store := newTestEscalationStore()
	h := NewEscalationHandler(store)

	body := `{"job_id":"backup","threshold":3,"interval_seconds":300,"webhook":"http://hook"}`
	req := httptest.NewRequest(http.MethodPost, "/escalation", bytes.NewBufferString(body))
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rw.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/escalation", nil)
	rw2 := httptest.NewRecorder()
	h.ServeHTTP(rw2, req2)

	var policies []monitor.EscalationPolicy
	if err := json.NewDecoder(rw2.Body).Decode(&policies); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(policies) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(policies))
	}
	if policies[0].Threshold != 3 {
		t.Errorf("expected threshold 3, got %d", policies[0].Threshold)
	}
	if policies[0].Interval != 300*time.Second {
		t.Errorf("unexpected interval: %v", policies[0].Interval)
	}
}

func TestEscalationHandler_FilterByJob(t *testing.T) {
	store := newTestEscalationStore()
	store.Set(monitor.EscalationPolicy{JobID: "alpha", Threshold: 1, Interval: time.Minute})
	store.Set(monitor.EscalationPolicy{JobID: "beta", Threshold: 2, Interval: time.Minute})
	h := NewEscalationHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/escalation?job_id=alpha", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)

	var policies []monitor.EscalationPolicy
	json.NewDecoder(rw.Body).Decode(&policies)
	if len(policies) != 1 || policies[0].JobID != "alpha" {
		t.Errorf("expected only alpha, got %+v", policies)
	}
}

func TestEscalationHandler_Delete(t *testing.T) {
	store := newTestEscalationStore()
	store.Set(monitor.EscalationPolicy{JobID: "job1", Threshold: 2, Interval: time.Minute})
	h := NewEscalationHandler(store)

	req := httptest.NewRequest(http.MethodDelete, "/escalation/job1", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rw.Code)
	}
	_, ok := store.Get("job1")
	if ok {
		t.Error("expected policy to be deleted")
	}
}

func TestEscalationHandler_Delete_NotFound(t *testing.T) {
	store := newTestEscalationStore()
	h := NewEscalationHandler(store)

	req := httptest.NewRequest(http.MethodDelete, "/escalation/ghost", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rw.Code)
	}
}

func TestEscalationHandler_MethodNotAllowed(t *testing.T) {
	store := newTestEscalationStore()
	h := NewEscalationHandler(store)

	req := httptest.NewRequest(http.MethodPatch, "/escalation", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rw.Code)
	}
}

func TestEscalationHandler_InvalidBody(t *testing.T) {
	store := newTestEscalationStore()
	h := NewEscalationHandler(store)

	req := httptest.NewRequest(http.MethodPost, "/escalation", bytes.NewBufferString(`{"threshold":0}`))
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rw.Code)
	}
}
