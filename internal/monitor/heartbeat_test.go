package monitor

import (
	"testing"
	"time"
)

func TestHeartbeatRecord_IsStale_NotStale(t *testing.T) {
	r := HeartbeatRecord{
		JobID:         "job1",
		LastSeen:      time.Now(),
		ExpectedEvery: 5 * time.Minute,
	}
	if r.IsStale(30*time.Second, time.Now()) {
		t.Error("expected record to not be stale immediately after check-in")
	}
}

func TestHeartbeatRecord_IsStale_Stale(t *testing.T) {
	r := HeartbeatRecord{
		JobID:         "job1",
		LastSeen:      time.Now().Add(-10 * time.Minute),
		ExpectedEvery: 5 * time.Minute,
	}
	if !r.IsStale(30*time.Second, time.Now()) {
		t.Error("expected record to be stale after missing interval + grace")
	}
}

func TestHeartbeatRecord_IsStale_ZeroLastSeen(t *testing.T) {
	r := HeartbeatRecord{
		JobID:         "job1",
		ExpectedEvery: 1 * time.Minute,
	}
	if r.IsStale(0, time.Now()) {
		t.Error("expected zero LastSeen to never be considered stale")
	}
}

func TestHeartbeatStore_RecordAndGet(t *testing.T) {
	s := NewHeartbeatStore()
	s.Record("jobA", 10*time.Minute)

	r, ok := s.Get("jobA")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if r.JobID != "jobA" {
		t.Errorf("expected jobID jobA, got %s", r.JobID)
	}
	if r.ExpectedEvery != 10*time.Minute {
		t.Errorf("expected interval 10m, got %v", r.ExpectedEvery)
	}
}

func TestHeartbeatStore_GetUnknownJob(t *testing.T) {
	s := NewHeartbeatStore()
	_, ok := s.Get("ghost")
	if ok {
		t.Error("expected no record for unknown job")
	}
}

func TestHeartbeatStore_All(t *testing.T) {
	s := NewHeartbeatStore()
	s.Record("j1", time.Minute)
	s.Record("j2", 2*time.Minute)

	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 records, got %d", len(all))
	}
}

func TestHeartbeatStore_StaleJobs(t *testing.T) {
	s := NewHeartbeatStore()
	s.Record("fresh", 10*time.Minute)

	// Manually inject a stale record
	s.mu.Lock()
	s.records["stale"] = HeartbeatRecord{
		JobID:         "stale",
		LastSeen:      time.Now().Add(-20 * time.Minute),
		ExpectedEvery: 5 * time.Minute,
	}
	s.mu.Unlock()

	stale := s.StaleJobs(30 * time.Second)
	if len(stale) != 1 {
		t.Fatalf("expected 1 stale job, got %d", len(stale))
	}
	if stale[0].JobID != "stale" {
		t.Errorf("expected stale job 'stale', got %s", stale[0].JobID)
	}
}
