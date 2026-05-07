package monitor

import (
	"sync"
	"time"
)

// RetryPolicy defines how many times an alert should be retried and the interval.
type RetryPolicy struct {
	MaxAttempts int
	Interval    time.Duration
}

// RetryState tracks retry attempts for a specific job.
type RetryState struct {
	JobName     string
	Attempts    int
	LastAttempt time.Time
	Exhausted   bool
}

// RetryTracker manages retry state for missed job alerts.
type RetryTracker struct {
	mu     sync.Mutex
	policy RetryPolicy
	states map[string]*RetryState
}

// NewRetryTracker creates a RetryTracker with the given policy.
func NewRetryTracker(policy RetryPolicy) *RetryTracker {
	if policy.MaxAttempts <= 0 {
		policy.MaxAttempts = 3
	}
	if policy.Interval <= 0 {
		policy.Interval = 5 * time.Minute
	}
	return &RetryTracker{
		policy: policy,
		states: make(map[string]*RetryState),
	}
}

// ShouldRetry returns true if a retry attempt should be made for the given job.
func (r *RetryTracker) ShouldRetry(jobName string, now time.Time) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, ok := r.states[jobName]
	if !ok {
		r.states[jobName] = &RetryState{
			JobName:     jobName,
			Attempts:    1,
			LastAttempt: now,
		}
		return true
	}
	if s.Exhausted {
		return false
	}
	if now.Sub(s.LastAttempt) < r.policy.Interval {
		return false
	}
	s.Attempts++
	s.LastAttempt = now
	if s.Attempts >= r.policy.MaxAttempts {
		s.Exhausted = true
	}
	return true
}

// Reset clears the retry state for a job (e.g. after a successful check-in).
func (r *RetryTracker) Reset(jobName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.states, jobName)
}

// All returns a snapshot of all current retry states.
func (r *RetryTracker) All() []RetryState {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]RetryState, 0, len(r.states))
	for _, s := range r.states {
		out = append(out, *s)
	}
	return out
}
