package monitor

import (
	"testing"
	"time"
)

func TestMetricStore_RecordAndSnapshot(t *testing.T) {
	s := NewMetricStore(100)
	s.Record("job1", false, 200*time.Millisecond)
	s.Record("job1", false, 400*time.Millisecond)

	snap, ok := s.Snapshot("job1")
	if !ok {
		t.Fatal("expected snapshot to exist")
	}
	if snap.RunCount != 2 {
		t.Errorf("expected RunCount=2, got %d", snap.RunCount)
	}
	if snap.FailCount != 0 {
		t.Errorf("expected FailCount=0, got %d", snap.FailCount)
	}
	if snap.AvgDuration != 300*time.Millisecond {
		t.Errorf("expected AvgDuration=300ms, got %v", snap.AvgDuration)
	}
}

func TestMetricStore_FailCount(t *testing.T) {
	s := NewMetricStore(100)
	s.Record("job2", false, 100*time.Millisecond)
	s.Record("job2", true, 100*time.Millisecond)
	s.Record("job2", true, 100*time.Millisecond)

	snap, ok := s.Snapshot("job2")
	if !ok {
		t.Fatal("expected snapshot")
	}
	if snap.FailCount != 2 {
		t.Errorf("expected FailCount=2, got %d", snap.FailCount)
	}
	if snap.RunCount != 3 {
		t.Errorf("expected RunCount=3, got %d", snap.RunCount)
	}
}

func TestMetricStore_GetUnknownJob(t *testing.T) {
	s := NewMetricStore(0)
	_, ok := s.Snapshot("ghost")
	if ok {
		t.Error("expected no snapshot for unknown job")
	}
}

func TestMetricStore_BoundedSize(t *testing.T) {
	s := NewMetricStore(3)
	for i := 0; i < 10; i++ {
		s.Record("job3", false, time.Millisecond)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.records["job3"]) != 3 {
		t.Errorf("expected bounded to 3, got %d", len(s.records["job3"]))
	}
}

func TestMetricStore_DefaultMaxSize(t *testing.T) {
	s := NewMetricStore(0)
	if s.maxSize != defaultMetricMaxSize {
		t.Errorf("expected default max size %d, got %d", defaultMetricMaxSize, s.maxSize)
	}
}

func TestMetricStore_All(t *testing.T) {
	s := NewMetricStore(100)
	s.Record("alpha", false, 50*time.Millisecond)
	s.Record("beta", true, 75*time.Millisecond)

	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(all))
	}
}
