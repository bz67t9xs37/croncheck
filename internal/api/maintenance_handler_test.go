package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/croncheck/internal/monitor"
)

func newTestMaintenanceStore() *monitor.MaintenanceStore {
	return monitor.NewMaintenanceStore()
}

func TestMaintenanceHandler_CreateAndList(t *testing.T) {
	store := newTestMaintenanceStore()
	h := NewMaintenanceHandler(store)

	now := time.Now().UTC().Truncate(time.Second)
	body, _ := json.Marshal(map[string]interface{}{
		"job_name": "daily-report",
		"start":    now,
		"end":      now.Add(2 * time.Hour),
	})
	req := httptest.NewRequest(http.MethodPost, "/maintenance", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/maintenance", nil)
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)

	var windows []monitor.MaintenanceWindow
	json.NewDecoder(rec2.Body).Decode(&windows)
	if len(windows) != 1 || windows[0].JobName != "daily-report" {
		t.Errorf("unexpected windows: %+v", windows)
	}
}

func TestMaintenanceHandler_FilterByJob(t *testing.T) {
	store := newTestMaintenanceStore()
	now := time.Now()
	store.Add("job-a", now, now.Add(time.Hour))
	store.Add("job-b", now, now.Add(time.Hour))
	h := NewMaintenanceHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/maintenance?job=job-a", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var windows []monitor.MaintenanceWindow
	json.NewDecoder(rec.Body).Decode(&windows)
	if len(windows) != 1 || windows[0].JobName != "job-a" {
		t.Errorf("expected only job-a, got %+v", windows)
	}
}

func TestMaintenanceHandler_Delete(t *testing.T) {
	store := newTestMaintenanceStore()
	now := time.Now()
	store.Add("job-a", now, now.Add(time.Hour))
	store.Add("job-a", now.Add(time.Hour), now.Add(2*time.Hour))
	h := NewMaintenanceHandler(store)

	req := httptest.NewRequest(http.MethodDelete, "/maintenance?job=job-a", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var result map[string]int
	json.NewDecoder(rec.Body).Decode(&result)
	if result["removed"] != 2 {
		t.Errorf("expected 2 removed, got %d", result["removed"])
	}
}

func TestMaintenanceHandler_MethodNotAllowed(t *testing.T) {
	h := NewMaintenanceHandler(newTestMaintenanceStore())
	req := httptest.NewRequest(http.MethodPatch, "/maintenance", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestMaintenanceHandler_CreateInvalidBody(t *testing.T) {
	h := NewMaintenanceHandler(newTestMaintenanceStore())
	req := httptest.NewRequest(http.MethodPost, "/maintenance", bytes.NewBufferString("not-json"))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}
