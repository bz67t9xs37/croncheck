package monitor

import (
	"fmt"
	"sync"
)

// Registry holds all monitored jobs.
type Registry struct {
	mu   sync.RWMutex
	jobs map[string]*Job
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		jobs: make(map[string]*Job),
	}
}

// Register adds a job to the registry.
func (r *Registry) Register(j *Job) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[j.Name] = j
}

// Get retrieves a job by name.
func (r *Registry) Get(name string) (*Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	j, ok := r.jobs[name]
	if !ok {
		return nil, fmt.Errorf("job %q not found", name)
	}
	return j, nil
}

// All returns a snapshot of all registered jobs.
func (r *Registry) All() []*Job {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Job, 0, len(r.jobs))
	for _, j := range r.jobs {
		out = append(out, j)
	}
	return out
}

// CheckIn records a check-in for the named job.
func (r *Registry) CheckIn(name string) error {
	j, err := r.Get(name)
	if err != nil {
		return err
	}
	j.CheckIn()
	return nil
}

// Len returns the number of registered jobs.
func (r *Registry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.jobs)
}
