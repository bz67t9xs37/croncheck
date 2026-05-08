package monitor

import (
	"sync"
	"time"
)

// JobSnapshot captures a point-in-time view of a job's state.
type JobSnapshot struct {
	JobID       string     `json:"job_id"`
	Status      JobStatus  `json:"status"`
	LastCheckIn *time.Time `json:"last_check_in,omitempty"`
	CapturedAt  time.Time  `json:"captured_at"`
}

// SnapshotStore holds the most recent snapshot for each job.
type SnapshotStore struct {
	mu        sync.RWMutex
	snapshots map[string]JobSnapshot
}

// NewSnapshotStore creates an empty SnapshotStore.
func NewSnapshotStore() *SnapshotStore {
	return &SnapshotStore{
		snapshots: make(map[string]JobSnapshot),
	}
}

// Record saves or updates the snapshot for a job.
func (s *SnapshotStore) Record(snap JobSnapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	snap.CapturedAt = time.Now()
	s.snapshots[snap.JobID] = snap
}

// Get returns the latest snapshot for a job and whether it exists.
func (s *SnapshotStore) Get(jobID string) (JobSnapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap, ok := s.snapshots[jobID]
	return snap, ok
}

// All returns a copy of all current snapshots.
func (s *SnapshotStore) All() []JobSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]JobSnapshot, 0, len(s.snapshots))
	for _, snap := range s.snapshots {
		out = append(out, snap)
	}
	return out
}

// Delete removes the snapshot for a job.
func (s *SnapshotStore) Delete(jobID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.snapshots, jobID)
}
