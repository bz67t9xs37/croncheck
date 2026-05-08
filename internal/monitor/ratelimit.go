package monitor

import (
	"sync"
	"time"
)

// RateLimitStore tracks alert rate limiting per job to prevent alert storms.
type RateLimitStore struct {
	mu       sync.Mutex
	records  map[string]time.Time
	cooldown time.Duration
}

// NewRateLimitStore creates a RateLimitStore with the given cooldown duration.
// Alerts for a given job will be suppressed if one was sent within the cooldown window.
func NewRateLimitStore(cooldown time.Duration) *RateLimitStore {
	if cooldown <= 0 {
		cooldown = 15 * time.Minute
	}
	return &RateLimitStore{
		records:  make(map[string]time.Time),
		cooldown: cooldown,
	}
}

// Allow returns true if an alert for the given job is permitted (i.e. not rate-limited).
// If allowed, it records the current time as the last alert time for the job.
func (r *RateLimitStore) Allow(jobID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	last, exists := r.records[jobID]
	if exists && time.Since(last) < r.cooldown {
		return false
	}
	r.records[jobID] = time.Now()
	return true
}

// Reset clears the rate limit record for the given job, allowing the next alert immediately.
func (r *RateLimitStore) Reset(jobID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, jobID)
}

// LastAlert returns the time of the last alert for the given job and whether one exists.
func (r *RateLimitStore) LastAlert(jobID string) (time.Time, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.records[jobID]
	return t, ok
}

// All returns a copy of all current rate limit records keyed by job ID.
func (r *RateLimitStore) All() map[string]time.Time {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make(map[string]time.Time, len(r.records))
	for k, v := range r.records {
		out[k] = v
	}
	return out
}
