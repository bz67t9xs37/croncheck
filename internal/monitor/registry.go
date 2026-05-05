package monitor

import (
	"fmt"
	"sync"
)

// Registry holds all monitored jobs and provides thread-safe access.
type Registry struct {
	mu   sync.RWMutex
	jobs map[string]*Job
}

// NewRegistry creates an empty job registry.
func NewRegistry() *Registry {
	return &Registry{
		jobs: make(map[string]*Job),
	}
}

// Register adds a new job to the registry.
// Returns an error if a job with the same ID already exists.
func (r *Registry) Register(job *Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.jobs[job.ID]; exists {
		return fmt.Errorf("job %q is already registered", job.ID)
	}
	r.jobs[job.ID] = job
	return nil
}

// Get retrieves a job by ID.
func (r *Registry) Get(id string) (*Job, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	j, ok := r.jobs[id]
	return j, ok
}

// CheckIn records a successful run for the given job ID.
func (r *Registry) CheckIn(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	j, ok := r.jobs[id]
	if !ok {
		return fmt.Errorf("job %q not found", id)
	}
	j.CheckIn()
	return nil
}

// All returns a snapshot of all registered jobs.
func (r *Registry) All() []*Job {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*Job, 0, len(r.jobs))
	for _, j := range r.jobs {
		list = append(list, j)
	}
	return list
}
