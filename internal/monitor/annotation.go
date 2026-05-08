package monitor

import (
	"fmt"
	"sync"
	"time"
)

// Annotation holds a timestamped note attached to a job.
type Annotation struct {
	JobID     string    `json:"job_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// AnnotationStore stores annotations keyed by job ID.
type AnnotationStore struct {
	mu      sync.RWMutex
	entries map[string][]Annotation
	maxPer  int
}

// NewAnnotationStore returns an AnnotationStore with a per-job cap.
func NewAnnotationStore(maxPerJob int) *AnnotationStore {
	if maxPerJob <= 0 {
		maxPerJob = 50
	}
	return &AnnotationStore{
		entries: make(map[string][]Annotation),
		maxPer:  maxPerJob,
	}
}

// Add appends an annotation for the given job, trimming oldest if over cap.
func (s *AnnotationStore) Add(jobID, message string) error {
	if jobID == "" {
		return fmt.Errorf("job ID must not be empty")
	}
	if message == "" {
		return fmt.Errorf("message must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	a := Annotation{JobID: jobID, Message: message, CreatedAt: time.Now().UTC()}
	s.entries[jobID] = append(s.entries[jobID], a)
	if len(s.entries[jobID]) > s.maxPer {
		s.entries[jobID] = s.entries[jobID][len(s.entries[jobID])-s.maxPer:]
	}
	return nil
}

// Get returns all annotations for a job.
func (s *AnnotationStore) Get(jobID string) []Annotation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := s.entries[jobID]
	out := make([]Annotation, len(list))
	copy(out, list)
	return out
}

// All returns every annotation across all jobs.
func (s *AnnotationStore) All() []Annotation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []Annotation
	for _, list := range s.entries {
		out = append(out, list...)
	}
	return out
}

// Delete removes all annotations for a job.
func (s *AnnotationStore) Delete(jobID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, jobID)
}
