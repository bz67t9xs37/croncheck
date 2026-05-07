package monitor

import (
	"sync"
	"time"
)

// MaintenanceWindow represents a scheduled window during which alerts are suppressed.
type MaintenanceWindow struct {
	JobName   string    `json:"job_name"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	CreatedAt time.Time `json:"created_at"`
}

// IsActive returns true if the window covers the given time.
func (w MaintenanceWindow) IsActive(t time.Time) bool {
	return !t.Before(w.Start) && t.Before(w.End)
}

// MaintenanceStore holds scheduled maintenance windows.
type MaintenanceStore struct {
	mu      sync.RWMutex
	windows []MaintenanceWindow
}

// NewMaintenanceStore creates an empty MaintenanceStore.
func NewMaintenanceStore() *MaintenanceStore {
	return &MaintenanceStore{}
}

// Add registers a new maintenance window.
func (s *MaintenanceStore) Add(jobName string, start, end time.Time) MaintenanceWindow {
	w := MaintenanceWindow{
		JobName:   jobName,
		Start:     start,
		End:       end,
		CreatedAt: time.Now(),
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.windows = append(s.windows, w)
	return w
}

// IsInMaintenance returns true if the given job has an active window at time t.
func (s *MaintenanceStore) IsInMaintenance(jobName string, t time.Time) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, w := range s.windows {
		if w.JobName == jobName && w.IsActive(t) {
			return true
		}
	}
	return false
}

// All returns a copy of all windows (including expired).
func (s *MaintenanceStore) All() []MaintenanceWindow {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]MaintenanceWindow, len(s.windows))
	copy(out, s.windows)
	return out
}

// Remove deletes all windows for the given job.
func (s *MaintenanceStore) Remove(jobName string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	filtered := s.windows[:0]
	count := 0
	for _, w := range s.windows {
		if w.JobName == jobName {
			count++
			continue
		}
		filtered = append(filtered, w)
	}
	s.windows = filtered
	return count
}
