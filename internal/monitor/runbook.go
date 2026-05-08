package monitor

import (
	"fmt"
	"sync"
)

// RunbookStore maps job IDs to runbook URLs or notes for on-call reference.
type RunbookStore struct {
	mu      sync.RWMutex
	entries map[string]RunbookEntry
}

// RunbookEntry holds the runbook URL and optional notes for a job.
type RunbookEntry struct {
	JobID string `json:"job_id"`
	URL   string `json:"url"`
	Notes string `json:"notes,omitempty"`
}

// NewRunbookStore creates an empty RunbookStore.
func NewRunbookStore() *RunbookStore {
	return &RunbookStore{
		entries: make(map[string]RunbookEntry),
	}
}

// Set adds or replaces the runbook entry for a job.
func (s *RunbookStore) Set(entry RunbookEntry) error {
	if entry.JobID == "" {
		return fmt.Errorf("job_id is required")
	}
	if entry.URL == "" {
		return fmt.Errorf("url is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[entry.JobID] = entry
	return nil
}

// Get returns the runbook entry for a job, if present.
func (s *RunbookStore) Get(jobID string) (RunbookEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[jobID]
	return e, ok
}

// Delete removes the runbook entry for a job.
func (s *RunbookStore) Delete(jobID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.entries[jobID]
	if ok {
		delete(s.entries, jobID)
	}
	return ok
}

// All returns all runbook entries.
func (s *RunbookStore) All() []RunbookEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]RunbookEntry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}
