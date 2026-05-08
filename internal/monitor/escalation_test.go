package monitor

import (
	"testing"
	"time"
)

func TestEscalationStore_SetAndGet(t *testing.T) {
	s := NewEscalationStore()
	p := EscalationPolicy{JobID: "job1", Threshold: 3, Interval: time.Minute, Webhook: "http://example.com"}
	s.Set(p)
	got, ok := s.Get("job1")
	if !ok {
		t.Fatal("expected policy to exist")
	}
	if got.Threshold != 3 {
		t.Errorf("expected threshold 3, got %d", got.Threshold)
	}
}

func TestEscalationStore_GetUnknown(t *testing.T) {
	s := NewEscalationStore()
	_, ok := s.Get("missing")
	if ok {
		t.Error("expected no policy for unknown job")
	}
}

func TestEscalationStore_Delete(t *testing.T) {
	s := NewEscalationStore()
	s.Set(EscalationPolicy{JobID: "job1", Threshold: 2, Interval: time.Minute})
	if !s.Delete("job1") {
		t.Fatal("expected delete to return true")
	}
	_, ok := s.Get("job1")
	if ok {
		t.Error("expected policy to be removed")
	}
}

func TestEscalationStore_Delete_NotFound(t *testing.T) {
	s := NewEscalationStore()
	if s.Delete("ghost") {
		t.Error("expected delete to return false for unknown job")
	}
}

func TestEscalationStore_ShouldEscalate_BelowThreshold(t *testing.T) {
	s := NewEscalationStore()
	s.Set(EscalationPolicy{JobID: "job1", Threshold: 3, Interval: time.Minute})
	now := time.Now()
	if s.ShouldEscalate("job1", now) {
		t.Error("should not escalate below threshold")
	}
	if s.ShouldEscalate("job1", now) {
		t.Error("should not escalate below threshold on second miss")
	}
}

func TestEscalationStore_ShouldEscalate_AtThreshold(t *testing.T) {
	s := NewEscalationStore()
	s.Set(EscalationPolicy{JobID: "job1", Threshold: 2, Interval: time.Minute})
	now := time.Now()
	s.ShouldEscalate("job1", now) // miss 1
	if !s.ShouldEscalate("job1", now) { // miss 2 — threshold reached
		t.Error("expected escalation at threshold")
	}
}

func TestEscalationStore_ShouldEscalate_IntervalNotElapsed(t *testing.T) {
	s := NewEscalationStore()
	s.Set(EscalationPolicy{JobID: "job1", Threshold: 1, Interval: time.Hour})
	now := time.Now()
	s.ShouldEscalate("job1", now) // triggers first escalation
	if s.ShouldEscalate("job1", now.Add(time.Minute)) {
		t.Error("should not escalate again within interval")
	}
}

func TestEscalationStore_ShouldEscalate_AfterInterval(t *testing.T) {
	s := NewEscalationStore()
	s.Set(EscalationPolicy{JobID: "job1", Threshold: 1, Interval: time.Minute})
	now := time.Now()
	s.ShouldEscalate("job1", now)
	if !s.ShouldEscalate("job1", now.Add(2*time.Minute)) {
		t.Error("expected escalation after interval elapsed")
	}
}

func TestEscalationStore_ResetMisses(t *testing.T) {
	s := NewEscalationStore()
	s.Set(EscalationPolicy{JobID: "job1", Threshold: 2, Interval: time.Minute})
	now := time.Now()
	s.ShouldEscalate("job1", now) // 1 miss
	s.ResetMisses("job1")
	if s.ShouldEscalate("job1", now) { // should be back to 1 after reset
		t.Error("should not escalate immediately after reset")
	}
}

func TestEscalationStore_All(t *testing.T) {
	s := NewEscalationStore()
	s.Set(EscalationPolicy{JobID: "a", Threshold: 1, Interval: time.Minute})
	s.Set(EscalationPolicy{JobID: "b", Threshold: 2, Interval: time.Minute})
	if len(s.All()) != 2 {
		t.Errorf("expected 2 policies, got %d", len(s.All()))
	}
}
