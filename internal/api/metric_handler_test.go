package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/croncheck/internal/monitor"
)

func newTestMetricStore() *monitor.MetricStore {
	return monitor.NewMetricStore(0)
}

func TestMetricHandler_GetMetrics_AllJobs(t *testing.T) {
	store := newTestMetricStore()

	now := time.Now()
	store.Record("job-a", monitor.MetricEntry{Timestamp: now, DurationMs: 120, Success: true})
	store.Record("job-b", monitor.MetricEntry{Timestamp: now, DurationMs: 340, Success: false})

	h := NewMetricHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(rw.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if _, ok := result["job-a"]; !ok {
		t.Error("expected job-a in response")
	}
	if _, ok := result["job-b"]; !ok {
		t.Error("expected job-b in response")
	}
}

func TestMetricHandler_GetMetrics_FilterByJob(t *testing.T) {
	store := newTestMetricStore()

	now := time.Now()
	store.Record("job-a", monitor.MetricEntry{Timestamp: now, DurationMs: 200, Success: true})
	store.Record("job-b", monitor.MetricEntry{Timestamp: now, DurationMs: 500, Success: false})

	h := NewMetricHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/metrics?job=job-a", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(rw.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if _, ok := result["job-a"]; !ok {
		t.Error("expected job-a in filtered response")
	}
	if _, ok := result["job-b"]; ok {
		t.Error("did not expect job-b in filtered response")
	}
}

func TestMetricHandler_GetMetrics_UnknownJob_ReturnsEmpty(t *testing.T) {
	store := newTestMetricStore()

	h := NewMetricHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/metrics?job=ghost", nil)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(rw.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected empty result for unknown job, got %v", result)
	}
}

func TestMetricHandler_MethodNotAllowed(t *testing.T) {
	store := newTestMetricStore()
	h := NewMetricHandler(store)

	for _, method := range []string{http.MethodPost, http.MethodDelete, http.MethodPut} {
		req := httptest.NewRequest(method, "/metrics", nil)
		rw := httptest.NewRecorder()

		h.ServeHTTP(rw, req)

		if rw.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s: expected 405, got %d", method, rw.Code)
		}
	}
}
