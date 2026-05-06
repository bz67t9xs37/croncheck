package monitor

import (
	"testing"
)

func TestHistory_RecordAndGet(t *testing.T) {
	h := NewHistory(10)
	h.Record("backup", StatusHealthy)
	h.Record("backup", StatusMissed)

	records := h.Get("backup")
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].Status != StatusHealthy {
		t.Errorf("expected first record StatusHealthy, got %s", records[0].Status)
	}
	if records[1].Status != StatusMissed {
		t.Errorf("expected second record StatusMissed, got %s", records[1].Status)
	}
}

func TestHistory_BoundedSize(t *testing.T) {
	const max = 5
	h := NewHistory(max)

	for i := 0; i < 12; i++ {
		h.Record("cleanup", StatusHealthy)
	}

	records := h.Get("cleanup")
	if len(records) != max {
		t.Errorf("expected %d records after overflow, got %d", max, len(records))
	}
}

func TestHistory_GetUnknownJob(t *testing.T) {
	h := NewHistory(10)
	records := h.Get("nonexistent")
	if records == nil {
		t.Error("expected non-nil slice for unknown job")
	}
	if len(records) != 0 {
		t.Errorf("expected 0 records for unknown job, got %d", len(records))
	}
}

func TestHistory_All(t *testing.T) {
	h := NewHistory(10)
	h.Record("job-a", StatusHealthy)
	h.Record("job-b", StatusMissed)
	h.Record("job-a", StatusHealthy)

	all := h.All()
	if len(all) != 3 {
		t.Errorf("expected 3 total records, got %d", len(all))
	}
}

func TestHistory_DefaultMaxSize(t *testing.T) {
	h := NewHistory(0) // should default to 50
	for i := 0; i < 60; i++ {
		h.Record("job", StatusHealthy)
	}
	records := h.Get("job")
	if len(records) != 50 {
		t.Errorf("expected default max 50, got %d", len(records))
	}
}
