package monitor

import (
	"time"
)

// EvaluationResult holds the outcome of evaluating all jobs.
type EvaluationResult struct {
	Missed  []*Job
	Failed  []*Job
	Healthy []*Job
}

// Evaluator periodically checks all registered jobs and reports results.
type Evaluator struct {
	registry *Registry
	interval time.Duration
	stopCh   chan struct{}
	OnAlert  func(result EvaluationResult)
}

// NewEvaluator creates an Evaluator that runs checks at the given interval.
func NewEvaluator(r *Registry, interval time.Duration, onAlert func(EvaluationResult)) *Evaluator {
	return &Evaluator{
		registry: r,
		interval: interval,
		stopCh:   make(chan struct{}),
		OnAlert:  onAlert,
	}
}

// Start begins the evaluation loop in a background goroutine.
func (e *Evaluator) Start() {
	go func() {
		ticker := time.NewTicker(e.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				e.runOnce(time.Now())
			case <-e.stopCh:
				return
			}
		}
	}()
}

// Stop halts the evaluation loop.
func (e *Evaluator) Stop() {
	close(e.stopCh)
}

// runOnce evaluates all jobs at the given time and fires OnAlert if any are non-healthy.
func (e *Evaluator) runOnce(now time.Time) {
	jobs := e.registry.All()
	result := EvaluationResult{}
	for _, j := range jobs {
		switch j.Evaluate(now) {
		case StatusMissed:
			result.Missed = append(result.Missed, j)
		case StatusFailed:
			result.Failed = append(result.Failed, j)
		default:
			result.Healthy = append(result.Healthy, j)
		}
	}
	if len(result.Missed) > 0 || len(result.Failed) > 0 {
		if e.OnAlert != nil {
			e.OnAlert(result)
		}
	}
}
