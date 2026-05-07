package monitor

import (
	"fmt"
	"sync"
	"time"
)

// DependencyLink represents a dependency relationship between two jobs.
type DependencyLink struct {
	Upstream   string    `json:"upstream"`
	Downstream string    `json:"downstream"`
	CreatedAt  time.Time `json:"created_at"`
}

// DependencyStore tracks upstream/downstream relationships between jobs.
type DependencyStore struct {
	mu    sync.RWMutex
	links []DependencyLink
}

// NewDependencyStore creates an empty DependencyStore.
func NewDependencyStore() *DependencyStore {
	return &DependencyStore{}
}

// Add records that downstream depends on upstream.
func (d *DependencyStore) Add(upstream, downstream string) error {
	if upstream == downstream {
		return fmt.Errorf("job cannot depend on itself: %q", upstream)
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, l := range d.links {
		if l.Upstream == upstream && l.Downstream == downstream {
			return fmt.Errorf("dependency already exists: %q -> %q", upstream, downstream)
		}
	}
	d.links = append(d.links, DependencyLink{
		Upstream:   upstream,
		Downstream: downstream,
		CreatedAt:  time.Now().UTC(),
	})
	return nil
}

// Remove deletes the dependency between upstream and downstream.
func (d *DependencyStore) Remove(upstream, downstream string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	for i, l := range d.links {
		if l.Upstream == upstream && l.Downstream == downstream {
			d.links = append(d.links[:i], d.links[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("dependency not found: %q -> %q", upstream, downstream)
}

// UpstreamsOf returns all jobs that the given job depends on.
func (d *DependencyStore) UpstreamsOf(job string) []string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var result []string
	for _, l := range d.links {
		if l.Downstream == job {
			result = append(result, l.Upstream)
		}
	}
	return result
}

// DownstreamsOf returns all jobs that depend on the given job.
func (d *DependencyStore) DownstreamsOf(job string) []string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var result []string
	for _, l := range d.links {
		if l.Upstream == job {
			result = append(result, l.Downstream)
		}
	}
	return result
}

// All returns every recorded dependency link.
func (d *DependencyStore) All() []DependencyLink {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]DependencyLink, len(d.links))
	copy(out, d.links)
	return out
}
