package monitor

import (
	"fmt"
	"sync"
)

// TagStore manages arbitrary key-value tags associated with monitored jobs.
// Tags can be used to group, filter, or annotate jobs (e.g. team, env, tier).
type TagStore struct {
	mu   sync.RWMutex
	tags map[string]map[string]string // jobID -> key -> value
}

// NewTagStore creates an empty TagStore.
func NewTagStore() *TagStore {
	return &TagStore{
		tags: make(map[string]map[string]string),
	}
}

// Set adds or updates a tag key/value for a job.
func (s *TagStore) Set(jobID, key, value string) error {
	if jobID == "" || key == "" {
		return fmt.Errorf("jobID and key must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.tags[jobID] == nil {
		s.tags[jobID] = make(map[string]string)
	}
	s.tags[jobID][key] = value
	return nil
}

// Get returns the tags for a job. Returns nil if the job has no tags.
func (s *TagStore) Get(jobID string) map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if m, ok := s.tags[jobID]; ok {
		copy := make(map[string]string, len(m))
		for k, v := range m {
			copy[k] = v
		}
		return copy
	}
	return nil
}

// Delete removes a single tag key from a job.
func (s *TagStore) Delete(jobID, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if m, ok := s.tags[jobID]; ok {
		delete(m, key)
		if len(m) == 0 {
			delete(s.tags, jobID)
		}
	}
}

// All returns a snapshot of all job tags.
func (s *TagStore) All() map[string]map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]map[string]string, len(s.tags))
	for jobID, m := range s.tags {
		copy := make(map[string]string, len(m))
		for k, v := range m {
			copy[k] = v
		}
		out[jobID] = copy
	}
	return out
}
