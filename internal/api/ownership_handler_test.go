package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/croncheck/internal/api"
	"github.com/croncheck/internal/monitor"
)

func newTestOwnershipStore() *monitor.OwnershipStore {
	return monitor.NewOwnershipStore()
}

func TestOwnershipHandler_SetAndGet(t *testing.T) {
	store := newTestOwnershipStore()
	h := api.NewOwnershipHandler(store)

	body, _ := json.Marshal(map[string]string{"job": "backup", "email": "ops@example.com", "team": "infra"})
	req := httptest.NewRequest(http.MethodPut, "/ownership", bytes.NewReader(body))
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rw.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/ownership?job=backup", nil)
	rw2 := httptest.NewRecorder()
	h.ServeHTTP(rw2, req2)
	if rw2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw2.Code)
	}
	var result map[string]string
	json.NewDecoder(rw2.Body).Decode(&result)
	if result["email"] != "ops@example.com" {
		t.Errorf("expected ops@example.com, got %s", result["email"])
	}
	if result["team"] != "infra" {
		t.Errorf("expected infra, got %s", result["team"])
	}
}

func TestOwnershipHandler_ListAll(t *testing.T) {
	store := newTestOwnershipStore()
	h := api.NewOwnershipHandler(store)

	for _, entry := range []struct{ job, email, team string }{
		{"job-a", "a@example.com", "team-a"},
		{"job-b", "b@example.com", "team-b"},
	} {
		body, _ := json.Marshal(map[string]string{"job": entry.job, "email": entry.email, "team": entry.team})
		req := httptest.NewRequest(http.MethodPut, "/ownership", bytes.NewReader(body))
		h.ServeHTTP(httptest.NewRecorder(), req)
	}

	req := httptest.NewRequest(http.MethodGet, "/ownership", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var all map[string]interface{}
	json.NewDecoder(rw.Body).Decode(&all)
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}

func TestOwnershipHandler_Delete(t *testing.T) {
	store := newTestOwnershipStore()
	h := api.NewOwnershipHandler(store)

	body, _ := json.Marshal(map[string]string{"job": "cleanup", "email": "dev@example.com", "team": "dev"})
	req := httptest.NewRequest(http.MethodPut, "/ownership", bytes.NewReader(body))
	h.ServeHTTP(httptest.NewRecorder(), req)

	req2 := httptest.NewRequest(http.MethodDelete, "/ownership?job=cleanup", nil)
	rw2 := httptest.NewRecorder()
	h.ServeHTTP(rw2, req2)
	if rw2.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rw2.Code)
	}

	req3 := httptest.NewRequest(http.MethodGet, "/ownership?job=cleanup", nil)
	rw3 := httptest.NewRecorder()
	h.ServeHTTP(rw3, req3)
	if rw3.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", rw3.Code)
	}
}

func TestOwnershipHandler_MethodNotAllowed(t *testing.T) {
	store := newTestOwnershipStore()
	h := api.NewOwnershipHandler(store)

	req := httptest.NewRequest(http.MethodPost, "/ownership", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rw.Code)
	}
}

func TestOwnershipHandler_GetUnknownJob(t *testing.T) {
	store := newTestOwnershipStore()
	h := api.NewOwnershipHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/ownership?job=ghost", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rw.Code)
	}
}
