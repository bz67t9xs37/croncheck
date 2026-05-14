package monitor

import (
	"testing"
	"time"
)

func TestAuditStore_RecordAndAll(t *testing.T) {
	s := NewAuditStore(10)
	s.Record(AuditEntry{JobID: "job1", Action: AuditActionCheckIn})
	s.Record(AuditEntry{JobID: "job2", Action: AuditActionCreated})

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if all[0].JobID != "job1" {
		t.Errorf("expected job1, got %s", all[0].JobID)
	}
}

func TestAuditStore_BoundedSize(t *testing.T) {
	s := NewAuditStore(3)
	for i := 0; i < 5; i++ {
		s.Record(AuditEntry{JobID: "job", Action: AuditActionUpdated})
	}
	if len(s.All()) != 3 {
		t.Fatalf("expected 3 entries after overflow, got %d", len(s.All()))
	}
}

func TestAuditStore_DefaultMaxSize(t *testing.T) {
	s := NewAuditStore(0)
	if s.maxSize != defaultMaxAuditEntries {
		t.Errorf("expected default max size %d, got %d", defaultMaxAuditEntries, s.maxSize)
	}
}

func TestAuditStore_ForJob(t *testing.T) {
	s := NewAuditStore(20)
	s.Record(AuditEntry{JobID: "alpha", Action: AuditActionCheckIn})
	s.Record(AuditEntry{JobID: "beta", Action: AuditActionCreated})
	s.Record(AuditEntry{JobID: "alpha", Action: AuditActionAlertion})

	entries := s.ForJob("alpha")
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries for alpha, got %d", len(entries))
	}
	for _, e := range entries {
		if e.JobID != "alpha" {
			t.Errorf("unexpected job_id %s", e.JobID)
		}
	}
}

func TestAuditStore_ForJob_Unknown(t *testing.T) {
	s := NewAuditStore(10)
	s.Record(AuditEntry{JobID: "job1", Action: AuditActionCheckIn})
	if entries := s.ForJob("unknown"); len(entries) != 0 {
		t.Errorf("expected empty slice for unknown job, got %d", len(entries))
	}
}

func TestAuditStore_TimestampAutoSet(t *testing.T) {
	s := NewAuditStore(10)
	before := time.Now().UTC()
	s.Record(AuditEntry{JobID: "job1", Action: AuditActionCreated})
	after := time.Now().UTC()

	all := s.All()
	ts := all[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", ts, before, after)
	}
}
