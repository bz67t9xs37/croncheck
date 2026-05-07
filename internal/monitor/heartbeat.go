package monitor

import (
	"sync"
	"time"
)

// HeartbeatRecord holds the last check-in time and interval expectation for a job.
type HeartbeatRecord struct {
	JobID        string
	LastSeen     time.Time
	ExpectedEvery time.Duration
}

// IsStale returns true if the job has not checked in within the expected interval plus grace.
func (h HeartbeatRecord) IsStale(grace time.Duration, now time.Time) bool {
	if h.LastSeen.IsZero() {
		return false
	}
	return now.After(h.LastSeen.Add(h.ExpectedEvery + grace))
}

// HeartbeatStore tracks the last heartbeat for each monitored job.
type HeartbeatStore struct {
	mu      sync.RWMutex
	records map[string]HeartbeatRecord
}

// NewHeartbeatStore creates an empty HeartbeatStore.
func NewHeartbeatStore() *HeartbeatStore {
	return &HeartbeatStore{
		records: make(map[string]HeartbeatRecord),
	}
}

// Record updates the last-seen timestamp for a job.
func (s *HeartbeatStore) Record(jobID string, expectedEvery time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[jobID] = HeartbeatRecord{
		JobID:         jobID,
		LastSeen:      time.Now(),
		ExpectedEvery: expectedEvery,
	}
}

// Get returns the heartbeat record for a job, and whether it exists.
func (s *HeartbeatStore) Get(jobID string) (HeartbeatRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.records[jobID]
	return r, ok
}

// All returns a copy of all heartbeat records.
func (s *HeartbeatStore) All() []HeartbeatRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]HeartbeatRecord, 0, len(s.records))
	for _, r := range s.records {
		out = append(out, r)
	}
	return out
}

// StaleJobs returns all records that are stale given the provided grace duration.
func (s *HeartbeatStore) StaleJobs(grace time.Duration) []HeartbeatRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	now := time.Now()
	var stale []HeartbeatRecord
	for _, r := range s.records {
		if r.IsStale(grace, now) {
			stale = append(stale, r)
		}
	}
	return stale
}
