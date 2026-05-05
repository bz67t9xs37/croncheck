package monitor

import (
	"sync"
	"testing"
	"time"
)

func newTestRegistry(jobs ...*Job) *Registry {
	r := NewRegistry()
	for _, j := range jobs {
		r.Register(j)
	}
	return r
}

func TestEvaluator_NoAlert_WhenAllHealthy(t *testing.T) {
	now := time.Now()
	j := &Job{
		Name:        "healthy-job",
		Schedule:    "* * * * *",
		GracePeriod: 5 * time.Minute,
		LastCheckIn: now.Add(-30 * time.Second),
	}
	r := newTestRegistry(j)
	alerted := false
	e := NewEvaluator(r, time.Hour, func(_ EvaluationResult) {
		alerted = true
	})
	e.runOnce(now)
	if alerted {
		t.Error("expected no alert for healthy job")
	}
}

func TestEvaluator_AlertFired_WhenJobMissed(t *testing.T) {
	now := time.Now()
	j := &Job{
		Name:        "missed-job",
		Schedule:    "* * * * *",
		GracePeriod: 1 * time.Minute,
		LastCheckIn: now.Add(-10 * time.Minute),
	}
	r := newTestRegistry(j)
	var mu sync.Mutex
	var captured EvaluationResult
	e := NewEvaluator(r, time.Hour, func(res EvaluationResult) {
		mu.Lock()
		captured = res
		mu.Unlock()
	})
	e.runOnce(now)
	mu.Lock()
	defer mu.Unlock()
	if len(captured.Missed) != 1 {
		t.Errorf("expected 1 missed job, got %d", len(captured.Missed))
	}
	if captured.Missed[0].Name != "missed-job" {
		t.Errorf("unexpected job name: %s", captured.Missed[0].Name)
	}
}

func TestEvaluator_StartStop(t *testing.T) {
	r := NewRegistry()
	called := make(chan struct{}, 1)
	e := NewEvaluator(r, 20*time.Millisecond, func(_ EvaluationResult) {
		select {
		case called <- struct{}{}:
		default:
		}
	})
	e.Start()
	time.Sleep(60 * time.Millisecond)
	e.Stop()
	// No panic and no deadlock means success.
}
