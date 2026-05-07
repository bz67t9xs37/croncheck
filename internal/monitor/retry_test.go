package monitor

import (
	"testing"
	"time"
)

func TestRetryTracker_FirstAttempt(t *testing.T) {
	rt := NewRetryTracker(RetryPolicy{MaxAttempts: 3, Interval: time.Minute})
	now := time.Now()
	if !rt.ShouldRetry("job1", now) {
		t.Fatal("expected first attempt to be allowed")
	}
	s := rt.states["job1"]
	if s.Attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", s.Attempts)
	}
}

func TestRetryTracker_TooSoon(t *testing.T) {
	rt := NewRetryTracker(RetryPolicy{MaxAttempts: 3, Interval: time.Minute})
	now := time.Now()
	rt.ShouldRetry("job1", now)
	if rt.ShouldRetry("job1", now.Add(10*time.Second)) {
		t.Fatal("expected retry to be denied when interval not elapsed")
	}
}

func TestRetryTracker_AllowsAfterInterval(t *testing.T) {
	rt := NewRetryTracker(RetryPolicy{MaxAttempts: 3, Interval: time.Minute})
	now := time.Now()
	rt.ShouldRetry("job1", now)
	if !rt.ShouldRetry("job1", now.Add(2*time.Minute)) {
		t.Fatal("expected retry after interval")
	}
}

func TestRetryTracker_ExhaustsAfterMaxAttempts(t *testing.T) {
	rt := NewRetryTracker(RetryPolicy{MaxAttempts: 3, Interval: time.Second})
	now := time.Now()
	for i := 0; i < 3; i++ {
		rt.ShouldRetry("job1", now.Add(time.Duration(i)*2*time.Second))
	}
	if rt.ShouldRetry("job1", now.Add(10*time.Second)) {
		t.Fatal("expected retries to be exhausted")
	}
	if !rt.states["job1"].Exhausted {
		t.Fatal("expected state to be marked exhausted")
	}
}

func TestRetryTracker_Reset(t *testing.T) {
	rt := NewRetryTracker(RetryPolicy{MaxAttempts: 2, Interval: time.Second})
	now := time.Now()
	rt.ShouldRetry("job1", now)
	rt.ShouldRetry("job1", now.Add(2*time.Second))
	rt.Reset("job1")
	if _, ok := rt.states["job1"]; ok {
		t.Fatal("expected state to be cleared after reset")
	}
	if !rt.ShouldRetry("job1", now.Add(3*time.Second)) {
		t.Fatal("expected fresh retry after reset")
	}
}

func TestRetryTracker_All(t *testing.T) {
	rt := NewRetryTracker(RetryPolicy{MaxAttempts: 3, Interval: time.Minute})
	now := time.Now()
	rt.ShouldRetry("jobA", now)
	rt.ShouldRetry("jobB", now)
	all := rt.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 states, got %d", len(all))
	}
}

func TestRetryTracker_Defaults(t *testing.T) {
	rt := NewRetryTracker(RetryPolicy{})
	if rt.policy.MaxAttempts != 3 {
		t.Errorf("expected default MaxAttempts=3, got %d", rt.policy.MaxAttempts)
	}
	if rt.policy.Interval != 5*time.Minute {
		t.Errorf("expected default Interval=5m, got %v", rt.policy.Interval)
	}
}
