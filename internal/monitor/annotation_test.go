package monitor

import (
	"testing"
)

func TestAnnotationStore_AddAndGet(t *testing.T) {
	s := NewAnnotationStore(10)
	if err := s.Add("job1", "first note"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := s.Add("job1", "second note"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list := s.Get("job1")
	if len(list) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(list))
	}
	if list[0].Message != "first note" {
		t.Errorf("expected 'first note', got %q", list[0].Message)
	}
}

func TestAnnotationStore_BoundedSize(t *testing.T) {
	s := NewAnnotationStore(3)
	for i := 0; i < 5; i++ {
		_ = s.Add("job1", "note")
	}
	if len(s.Get("job1")) != 3 {
		t.Errorf("expected cap of 3")
	}
}

func TestAnnotationStore_DefaultMaxSize(t *testing.T) {
	s := NewAnnotationStore(0)
	if s.maxPer != 50 {
		t.Errorf("expected default maxPer=50, got %d", s.maxPer)
	}
}

func TestAnnotationStore_GetUnknownJob(t *testing.T) {
	s := NewAnnotationStore(10)
	if list := s.Get("unknown"); len(list) != 0 {
		t.Errorf("expected empty slice for unknown job")
	}
}

func TestAnnotationStore_All(t *testing.T) {
	s := NewAnnotationStore(10)
	_ = s.Add("job1", "note a")
	_ = s.Add("job2", "note b")
	if len(s.All()) != 2 {
		t.Errorf("expected 2 total annotations")
	}
}

func TestAnnotationStore_Delete(t *testing.T) {
	s := NewAnnotationStore(10)
	_ = s.Add("job1", "note")
	s.Delete("job1")
	if len(s.Get("job1")) != 0 {
		t.Errorf("expected annotations to be deleted")
	}
}

func TestAnnotationStore_Add_Validation(t *testing.T) {
	s := NewAnnotationStore(10)
	if err := s.Add("", "msg"); err == nil {
		t.Error("expected error for empty job ID")
	}
	if err := s.Add("job1", ""); err == nil {
		t.Error("expected error for empty message")
	}
}
