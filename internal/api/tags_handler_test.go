package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/croncheck/internal/monitor"
)

func newTestTagStore() *monitor.TagStore {
	return monitor.NewTagStore()
}

func TestTagsHandler_SetAndGet(t *testing.T) {
	h := NewTagsHandler(newTestTagStore())

	body := `{"key":"env","value":"prod"}`
	req := httptest.NewRequest(http.MethodPut, "/api/tags/job1", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	req.URL.Path = "/api/tags/job1"
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/tags/job1", nil)
	req2.URL.Path = "/api/tags/job1"
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)
	var tags map[string]string
	if err := json.NewDecoder(rec2.Body).Decode(&tags); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if tags["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", tags["env"])
	}
}

func TestTagsHandler_ListAll(t *testing.T) {
	store := newTestTagStore()
	_ = store.Set("jobA", "tier", "critical")
	_ = store.Set("jobB", "env", "dev")
	h := NewTagsHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/api/tags", nil)
	req.URL.Path = "/api/tags"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var all map[string]map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&all); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(all))
	}
}

func TestTagsHandler_Delete(t *testing.T) {
	store := newTestTagStore()
	_ = store.Set("job1", "env", "staging")
	h := NewTagsHandler(store)

	req := httptest.NewRequest(http.MethodDelete, "/api/tags/job1?key=env", nil)
	req.URL.Path = "/api/tags/job1"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	if tags := store.Get("job1"); tags != nil {
		t.Error("expected tags to be empty after delete")
	}
}

func TestTagsHandler_MethodNotAllowed(t *testing.T) {
	h := NewTagsHandler(newTestTagStore())
	req := httptest.NewRequest(http.MethodPost, "/api/tags", nil)
	req.URL.Path = "/api/tags"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestTagsHandler_GetUnknownJob_ReturnsEmpty(t *testing.T) {
	h := NewTagsHandler(newTestTagStore())
	req := httptest.NewRequest(http.MethodGet, "/api/tags/ghost", nil)
	req.URL.Path = "/api/tags/ghost"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	var tags map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&tags); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected empty map, got %v", tags)
	}
}
