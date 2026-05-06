package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/croncheck/internal/monitor"
)

func newTestAlertLog() *monitor.AlertLog {
	log := monitor.NewAlertLog(10)
	log.Record(monitor.AlertEntry{JobName: "backup", Status: monitor.StatusMissed, FiredAt: time.Now(), Message: "missed run"})
	log.Record(monitor.AlertEntry{JobName: "sync", Status: monitor.StatusMissed, FiredAt: time.Now(), Message: "missed run"})
	log.Record(monitor.AlertEntry{JobName: "backup", Status: monitor.StatusMissed, FiredAt: time.Now(), Message: "missed again"})
	return log
}

func TestAlertLogHandler_AllAlerts(t *testing.T) {
	h := NewAlertLogHandler(newTestAlertLog())
	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var entries []monitor.AlertEntry
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestAlertLogHandler_FilterByJob(t *testing.T) {
	h := NewAlertLogHandler(newTestAlertLog())
	req := httptest.NewRequest(http.MethodGet, "/alerts?job=backup", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var entries []monitor.AlertEntry
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries for backup, got %d", len(entries))
	}
}

func TestAlertLogHandler_UnknownJob_ReturnsEmpty(t *testing.T) {
	h := NewAlertLogHandler(newTestAlertLog())
	req := httptest.NewRequest(http.MethodGet, "/alerts?job=ghost", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var entries []monitor.AlertEntry
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestAlertLogHandler_MethodNotAllowed(t *testing.T) {
	h := NewAlertLogHandler(newTestAlertLog())
	req := httptest.NewRequest(http.MethodPost, "/alerts", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}
