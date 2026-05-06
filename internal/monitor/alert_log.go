package monitor

import (
	"sync"
	"time"
)

// AlertEntry records a single alert that was fired for a job.
type AlertEntry struct {
	JobName   string
	Status    JobStatus
	FiredAt   time.Time
	Message   string
}

// AlertLog maintains an in-memory log of recent alerts fired by the evaluator.
type AlertLog struct {
	mu      sync.RWMutex
	entries []AlertEntry
	maxSize int
}

const defaultAlertLogSize = 100

// NewAlertLog creates an AlertLog with an optional maximum size.
// If maxSize <= 0, defaultAlertLogSize is used.
func NewAlertLog(maxSize int) *AlertLog {
	if maxSize <= 0 {
		maxSize = defaultAlertLogSize
	}
	return &AlertLog{
		entries: make([]AlertEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Record appends a new alert entry, evicting the oldest if at capacity.
func (al *AlertLog) Record(entry AlertEntry) {
	al.mu.Lock()
	defer al.mu.Unlock()
	if len(al.entries) >= al.maxSize {
		al.entries = al.entries[1:]
	}
	al.entries = append(al.entries, entry)
}

// All returns a copy of all recorded alert entries.
func (al *AlertLog) All() []AlertEntry {
	al.mu.RLock()
	defer al.mu.RUnlock()
	result := make([]AlertEntry, len(al.entries))
	copy(result, al.entries)
	return result
}

// ForJob returns all alert entries for a specific job name.
func (al *AlertLog) ForJob(name string) []AlertEntry {
	al.mu.RLock()
	defer al.mu.RUnlock()
	var result []AlertEntry
	for _, e := range al.entries {
		if e.JobName == name {
			result = append(result, e)
		}
	}
	return result
}
