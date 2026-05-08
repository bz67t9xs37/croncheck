package monitor

import (
	"testing"
)

func TestRunbookStore_SetAndGet(t *testing.T) {
	s := NewRunbookStore()
	entry := RunbookEntry{JobID: "backup", URL: "https://wiki.example.com/backup", Notes: "Check S3 bucket"}
	if err := s.Set(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := s.Get("backup")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got.URL != entry.URL {
		t.Errorf("expected URL %q, got %q", entry.URL, got.URL)
	}
	if got.Notes != entry.Notes {
		t.Errorf("expected Notes %q, got %q", entry.Notes, got.Notes)
	}
}

func TestRunbookStore_GetUnknownJob(t *testing.T) {
	s := NewRunbookStore()
	_, ok := s.Get("nonexistent")
	if ok {
		t.Fatal("expected no entry for unknown job")
	}
}

func TestRunbookStore_Delete(t *testing.T) {
	s := NewRunbookStore()
	_ = s.Set(RunbookEntry{JobID: "cleanup", URL: "https://wiki.example.com/cleanup"})
	if !s.Delete("cleanup") {
		t.Fatal("expected delete to return true")
	}
	_, ok := s.Get("cleanup")
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestRunbookStore_Delete_NotFound(t *testing.T) {
	s := NewRunbookStore()
	if s.Delete("ghost") {
		t.Fatal("expected delete to return false for unknown job")
	}
}

func TestRunbookStore_All(t *testing.T) {
	s := NewRunbookStore()
	_ = s.Set(RunbookEntry{JobID: "job1", URL: "https://wiki.example.com/job1"})
	_ = s.Set(RunbookEntry{JobID: "job2", URL: "https://wiki.example.com/job2"})
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}

func TestRunbookStore_Set_Validation(t *testing.T) {
	s := NewRunbookStore()
	if err := s.Set(RunbookEntry{JobID: "", URL: "https://example.com"}); err == nil {
		t.Error("expected error for empty job_id")
	}
	if err := s.Set(RunbookEntry{JobID: "job1", URL: ""}); err == nil {
		t.Error("expected error for empty url")
	}
}
