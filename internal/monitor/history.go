package monitor

import (
	"sync"
	"time"
)

// CheckInRecord represents a single check-in event for a job.
type CheckInRecord struct {
	JobName   string
	Timestamp time.Time
	Status    JobStatus
}

// History stores a bounded ring buffer of check-in records per job.
type History struct {
	mu      sync.RWMutex
	records map[string][]CheckInRecord
	maxSize int
}

// NewHistory creates a History that retains up to maxSize records per job.
func NewHistory(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 50
	}
	return &History{
		records: make(map[string][]CheckInRecord),
		maxSize: maxSize,
	}
}

// Record appends a check-in event for the given job.
func (h *History) Record(jobName string, status JobStatus) {
	h.mu.Lock()
	defer h.mu.Unlock()

	entry := CheckInRecord{
		JobName:   jobName,
		Timestamp: time.Now().UTC(),
		Status:    status,
	}

	buf := h.records[jobName]
	buf = append(buf, entry)
	if len(buf) > h.maxSize {
		buf = buf[len(buf)-h.maxSize:]
	}
	h.records[jobName] = buf
}

// Get returns a copy of the recorded events for the given job.
func (h *History) Get(jobName string) []CheckInRecord {
	h.mu.RLock()
	defer h.mu.RUnlock()

	src := h.records[jobName]
	out := make([]CheckInRecord, len(src))
	copy(out, src)
	return out
}

// All returns a copy of every recorded event across all jobs.
func (h *History) All() []CheckInRecord {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var out []CheckInRecord
	for _, recs := range h.records {
		out = append(out, recs...)
	}
	return out
}
