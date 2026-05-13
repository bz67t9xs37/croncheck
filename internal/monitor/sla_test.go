package monitor

import (
	"testing"
	"time"
)

func TestSLAStore_SetAndGet(t *testing.T) {
	s := NewSLAStore(0)
	p := SLAPolicy{JobID: "job1", MinSuccessRate: 0.95, MaxDowntime: 5 * time.Minute}
	s.Set(p)

	got, ok := s.Get("job1")
	if !ok {
		t.Fatal("expected policy to exist")
	}
	if got.MinSuccessRate != 0.95 {
		t.Errorf("expected 0.95, got %f", got.MinSuccessRate)
	}
}

func TestSLAStore_GetUnknown(t *testing.T) {
	s := NewSLAStore(0)
	_, ok := s.Get("ghost")
	if ok {
		t.Error("expected no policy for unknown job")
	}
}

func TestSLAStore_Delete(t *testing.T) {
	s := NewSLAStore(0)
	s.Set(SLAPolicy{JobID: "job1"})
	if !s.Delete("job1") {
		t.Error("expected delete to return true")
	}
	_, ok := s.Get("job1")
	if ok {
		t.Error("expected policy to be removed")
	}
	if s.Delete("job1") {
		t.Error("expected false when deleting non-existent policy")
	}
}

func TestSLAStore_RecordViolation_BoundedSize(t *testing.T) {
	s := NewSLAStore(3)
	for i := 0; i < 5; i++ {
		s.RecordViolation("job1", "test violation")
	}
	if len(s.Violations("")) != 3 {
		t.Errorf("expected 3 violations (bounded), got %d", len(s.Violations("")))
	}
}

func TestSLAStore_Violations_FilterByJob(t *testing.T) {
	s := NewSLAStore(0)
	s.RecordViolation("job1", "reason a")
	s.RecordViolation("job2", "reason b")
	s.RecordViolation("job1", "reason c")

	v := s.Violations("job1")
	if len(v) != 2 {
		t.Errorf("expected 2 violations for job1, got %d", len(v))
	}
}

func TestSLAStore_Evaluate_SuccessRateViolation(t *testing.T) {
	s := NewSLAStore(0)
	s.Set(SLAPolicy{JobID: "job1", MinSuccessRate: 0.9})
	s.Evaluate("job1", JobMetrics{SuccessCount: 7, FailCount: 3})

	v := s.Violations("job1")
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestSLAStore_Evaluate_NoViolationWhenHealthy(t *testing.T) {
	s := NewSLAStore(0)
	s.Set(SLAPolicy{JobID: "job1", MinSuccessRate: 0.8, MaxDowntime: 10 * time.Minute})
	s.Evaluate("job1", JobMetrics{SuccessCount: 9, FailCount: 1, LastDowntime: 2 * time.Minute})

	if len(s.Violations("job1")) != 0 {
		t.Error("expected no violations for healthy job")
	}
}

func TestSLAStore_Evaluate_DowntimeViolation(t *testing.T) {
	s := NewSLAStore(0)
	s.Set(SLAPolicy{JobID: "job1", MinSuccessRate: 0.0, MaxDowntime: 5 * time.Minute})
	s.Evaluate("job1", JobMetrics{LastDowntime: 10 * time.Minute})

	v := s.Violations("job1")
	if len(v) != 1 {
		t.Fatalf("expected 1 downtime violation, got %d", len(v))
	}
}

func TestSLAStore_All(t *testing.T) {
	s := NewSLAStore(0)
	s.Set(SLAPolicy{JobID: "job1"})
	s.Set(SLAPolicy{JobID: "job2"})
	if len(s.All()) != 2 {
		t.Errorf("expected 2 policies, got %d", len(s.All()))
	}
}
