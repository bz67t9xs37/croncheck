package monitor

import (
	"sync"
	"time"
)

// SilenceEntry represents a silenced job for a given time window.
type SilenceEntry struct {
	JobName  string    `json:"job_name"`
	SilentUntil time.Time `json:"silent_until"`
	Reason   string    `json:"reason"`
}

// SilenceStore tracks jobs that should not fire alerts during a window.
type SilenceStore struct {
	mu      sync.RWMutex
	silences map[string]SilenceEntry
}

// NewSilenceStore creates an empty SilenceStore.
func NewSilenceStore() *SilenceStore {
	return &SilenceStore{
		silences: make(map[string]SilenceEntry),
	}
}

// Silence adds or updates a silence for a job until the given time.
func (s *SilenceStore) Silence(jobName string, until time.Time, reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.silences[jobName] = SilenceEntry{
		JobName:     jobName,
		SilentUntil: until,
		Reason:      reason,
	}
}

// Unsilence removes a silence entry for a job.
func (s *SilenceStore) Unsilence(jobName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.silences, jobName)
}

// IsSilenced returns true if the job has an active silence at the given time.
func (s *SilenceStore) IsSilenced(jobName string, at time.Time) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.silences[jobName]
	if !ok {
		return false
	}
	if at.After(entry.SilentUntil) {
		return false
	}
	return true
}

// All returns all current silence entries, including expired ones.
func (s *SilenceStore) All() []SilenceEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]SilenceEntry, 0, len(s.silences))
	for _, e := range s.silences {
		out = append(out, e)
	}
	return out
}
