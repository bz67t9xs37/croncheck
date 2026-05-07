package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/croncheck/internal/monitor"
)

func newTestDependencyStore() *monitor.DependencyStore {
	return monitor.NewDependencyStore()
}

func TestDependencyHandler_CreateAndList(t *testing.T) {
	h := NewDependencyHandler(newTestDependencyStore())

	body := `{"upstream":"job-a","downstream":"job-b"}`
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/dependencies", bytes.NewBufferString(body)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, httptest.NewRequest(http.MethodGet, "/dependencies", nil))
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec2.Code)
	}
	var links []map[string]interface{}
	json.NewDecoder(rec2.Body).Decode(&links)
	if len(links) != 1 {
		t.Errorf("expected 1 link, got %d", len(links))
	}
}

func TestDependencyHandler_FilterByJob(t *testing.T) {
	store := newTestDependencyStore()
	_ = store.Add("job-a", "job-b")
	h := NewDependencyHandler(store)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/dependencies?job=job-b", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var result map[string][]string
	json.NewDecoder(rec.Body).Decode(&result)
	if len(result["upstreams"]) != 1 || result["upstreams"][0] != "job-a" {
		t.Errorf("unexpected upstreams: %v", result["upstreams"])
	}
}

func TestDependencyHandler_Delete(t *testing.T) {
	store := newTestDependencyStore()
	_ = store.Add("job-a", "job-b")
	h := NewDependencyHandler(store)

	body := `{"upstream":"job-a","downstream":"job-b"}`
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/dependencies", bytes.NewBufferString(body)))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	if len(store.All()) != 0 {
		t.Error("expected store to be empty after delete")
	}
}

func TestDependencyHandler_Delete_NotFound(t *testing.T) {
	h := NewDependencyHandler(newTestDependencyStore())
	body := `{"upstream":"ghost","downstream":"ghost2"}`
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/dependencies", bytes.NewBufferString(body)))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestDependencyHandler_MethodNotAllowed(t *testing.T) {
	h := NewDependencyHandler(newTestDependencyStore())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/dependencies", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestDependencyHandler_DuplicateConflict(t *testing.T) {
	h := NewDependencyHandler(newTestDependencyStore())
	body := `{"upstream":"job-a","downstream":"job-b"}`
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/dependencies", bytes.NewBufferString(body)))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/dependencies", bytes.NewBufferString(body)))
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", rec.Code)
	}
}
