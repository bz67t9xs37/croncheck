package monitor

import (
	"fmt"
	"sync"
)

// Owner holds contact information for a job owner.
type Owner struct {
	Job   string `json:"job"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Team  string `json:"team,omitempty"`
}

// OwnershipStore maps jobs to their owners.
type OwnershipStore struct {
	mu     sync.RWMutex
	owners map[string]Owner
}

// NewOwnershipStore creates an empty OwnershipStore.
func NewOwnershipStore() *OwnershipStore {
	return &OwnershipStore{
		owners: make(map[string]Owner),
	}
}

// Set assigns an owner to a job, overwriting any existing entry.
func (s *OwnershipStore) Set(o Owner) error {
	if o.Job == "" {
		return fmt.Errorf("job name must not be empty")
	}
	if o.Email == "" {
		return fmt.Errorf("owner email must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.owners[o.Job] = o
	return nil
}

// Get returns the owner for a job, or an error if none is registered.
func (s *OwnershipStore) Get(job string) (Owner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.owners[job]
	if !ok {
		return Owner{}, fmt.Errorf("no owner registered for job %q", job)
	}
	return o, nil
}

// Delete removes the owner entry for a job.
func (s *OwnershipStore) Delete(job string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.owners[job]; !ok {
		return fmt.Errorf("no owner registered for job %q", job)
	}
	delete(s.owners, job)
	return nil
}

// All returns a snapshot of all registered owners.
func (s *OwnershipStore) All() []Owner {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Owner, 0, len(s.owners))
	for _, o := range s.owners {
		out = append(out, o)
	}
	return out
}
