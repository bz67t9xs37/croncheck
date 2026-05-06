package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/croncheck/internal/monitor"
)

func newTestHistory() *monitor.History {
	return monitor.NewHistory(10)
}

func TestHistoryHandler_AllRecords(t *testing.T) {
	h := newTestHistory()
	h.Record("job-a", monitor.StatusHealthy)
	h.Record("job-b", monitor.StatusMissed)

	handler := NewHistoryHandler(h)
	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var records []monitor.CheckInRecord
	if err := json.NewDecoder(rr.Body).Decode(&records); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(records) != 2 {
		t.Errorf("expected 2 records, got %d", len(records))
	}
}

func TestHistoryHandler_FilterByJob(t *testing.T) {
	h := newTestHistory()
	h.Record("job-a", monitor.StatusHealthy)
	h.Record("job-b", monitor.StatusMissed)
	h.Record("job-a", monitor.StatusHealthy)

	handler := NewHistoryHandler(h)
	req := httptest.NewRequest(http.MethodGet, "/history/job-a", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var records []monitor.CheckInRecord
	if err := json.NewDecoder(rr.Body).Decode(&records); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(records) != 2 {
		t.Errorf("expected 2 records for job-a, got %d", len(records))
	}
	for _, rec := range records {
		if rec.JobName != "job-a" {
			t.Errorf("unexpected job name %q in filtered results", rec.JobName)
		}
	}
}

func TestHistoryHandler_MethodNotAllowed(t *testing.T) {
	h := newTestHistory()
	handler := NewHistoryHandler(h)
	req := httptest.NewRequest(http.MethodPost, "/history", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestHistoryHandler_UnknownJob_ReturnsEmpty(t *testing.T) {
	h := newTestHistory()
	handler := NewHistoryHandler(h)
	req := httptest.NewRequest(http.MethodGet, "/history/ghost", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var records []monitor.CheckInRecord
	if err := json.NewDecoder(rr.Body).Decode(&records); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected 0 records for unknown job, got %d", len(records))
	}
}
