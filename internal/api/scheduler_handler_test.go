package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/croncheck/internal/monitor"
)

func newTestSchedulerRegistry() *monitor.Registry {
	reg := monitor.NewRegistry()
	reg.Register(monitor.Job{
		Name:          "job-a",
		Schedule:      "@hourly",
		GracePeriod:   5 * time.Minute,
		ExpectedEvery: time.Hour,
	})
	return reg
}

func TestSchedulerHandler_Status_Running(t *testing.T) {
	reg := newTestSchedulerRegistry()
	h := NewSchedulerHandler(reg, func() bool { return true })

	req := httptest.NewRequest(http.MethodGet, "/scheduler/status", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp SchedulerStatusResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Running {
		t.Error("expected running=true")
	}
	if resp.JobCount != 1 {
		t.Errorf("expected job_count=1, got %d", resp.JobCount)
	}
}

func TestSchedulerHandler_Status_Stopped(t *testing.T) {
	reg := newTestSchedulerRegistry()
	h := NewSchedulerHandler(reg, func() bool { return false })

	req := httptest.NewRequest(http.MethodGet, "/scheduler/status", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp SchedulerStatusResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Running {
		t.Error("expected running=false")
	}
}

func TestSchedulerHandler_MethodNotAllowed(t *testing.T) {
	reg := newTestSchedulerRegistry()
	h := NewSchedulerHandler(reg, func() bool { return true })

	for _, method := range []string{http.MethodPost, http.MethodDelete, http.MethodPut} {
		req := httptest.NewRequest(method, "/scheduler/status", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected 405, got %d", method, rec.Code)
		}
	}
}
