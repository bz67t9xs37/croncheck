package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/user/croncheck/internal/monitor"
)

func newTestRegistry(t *testing.T) *monitor.Registry {
	t.Helper()
	r := monitor.NewRegistry()
	r.Register(&monitor.Job{
		Name:        "backup",
		GracePeriod: 5 * time.Minute,
	})
	return r
}

func TestHandleCheckIn_Success(t *testing.T) {
	r := newTestRegistry(t)
	h := NewHandler(r)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/checkin/backup", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rw.Code)
	}
}

func TestHandleCheckIn_UnknownJob(t *testing.T) {
	r := newTestRegistry(t)
	h := NewHandler(r)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/checkin/unknown", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rw.Code)
	}
}

func TestHandleCheckIn_MethodNotAllowed(t *testing.T) {
	r := newTestRegistry(t)
	h := NewHandler(r)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/checkin/backup", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rw.Code)
	}
}

func TestHandleStatus(t *testing.T) {
	r := newTestRegistry(t)
	h := NewHandler(r)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var resp statusResponse
	if err := json.NewDecoder(strings.NewReader(rw.Body.String())).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if len(resp.Jobs) != 1 {
		t.Errorf("expected 1 job in status, got %d", len(resp.Jobs))
	}
}
