package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/croncheck/internal/monitor"
)

func newTestAnnotationStore() *monitor.AnnotationStore {
	return monitor.NewAnnotationStore(10)
}

func TestAnnotationHandler_CreateAndList(t *testing.T) {
	h := NewAnnotationHandler(newTestAnnotationStore())
	body := `{"job_id":"job1","message":"deployed v2"}`
	req := httptest.NewRequest(http.MethodPost, "/annotations", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/annotations?job_id=job1", nil)
	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr2.Code)
	}
	var list []monitor.Annotation
	_ = json.NewDecoder(rr2.Body).Decode(&list)
	if len(list) != 1 || list[0].Message != "deployed v2" {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestAnnotationHandler_ListAll(t *testing.T) {
	h := NewAnnotationHandler(newTestAnnotationStore())
	for _, id := range []string{"job1", "job2"} {
		body := bytes.NewBufferString(`{"job_id":"` + id + `","message":"note"}`)
		req := httptest.NewRequest(http.MethodPost, "/annotations", body)
		h.ServeHTTP(httptest.NewRecorder(), req)
	}
	req := httptest.NewRequest(http.MethodGet, "/annotations", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	var list []monitor.Annotation
	_ = json.NewDecoder(rr.Body).Decode(&list)
	if len(list) != 2 {
		t.Errorf("expected 2 annotations, got %d", len(list))
	}
}

func TestAnnotationHandler_Delete(t *testing.T) {
	h := NewAnnotationHandler(newTestAnnotationStore())
	body := bytes.NewBufferString(`{"job_id":"job1","message":"note"}`)
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/annotations", body))

	req := httptest.NewRequest(http.MethodDelete, "/annotations?job_id=job1", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}

func TestAnnotationHandler_MethodNotAllowed(t *testing.T) {
	h := NewAnnotationHandler(newTestAnnotationStore())
	req := httptest.NewRequest(http.MethodPatch, "/annotations", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestAnnotationHandler_Create_InvalidJSON(t *testing.T) {
	h := NewAnnotationHandler(newTestAnnotationStore())
	req := httptest.NewRequest(http.MethodPost, "/annotations", bytes.NewBufferString("not-json"))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}
