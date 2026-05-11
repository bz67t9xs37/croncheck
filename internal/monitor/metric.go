package monitor

import (
	"sync"
	"time"
)

// MetricSnapshot holds aggregated runtime metrics for a job.
type MetricSnapshot struct {
	JobID       string        `json:"job_id"`
	RunCount    int           `json:"run_count"`
	FailCount   int           `json:"fail_count"`
	AvgDuration time.Duration `json:"avg_duration_ms"`
	LastRun     time.Time     `json:"last_run"`
}

// metricEntry is the internal record for a single run observation.
type metricEntry struct {
	at       time.Time
	failed   bool
	duration time.Duration
}

// MetricStore accumulates per-job run metrics.
type MetricStore struct {
	mu      sync.RWMutex
	records map[string][]metricEntry
	maxSize int
}

const defaultMetricMaxSize = 500

// NewMetricStore returns a MetricStore with a bounded history per job.
func NewMetricStore(maxSize int) *MetricStore {
	if maxSize <= 0 {
		maxSize = defaultMetricMaxSize
	}
	return &MetricStore{
		records: make(map[string][]metricEntry),
		maxSize: maxSize,
	}
}

// Record appends a run observation for the given job.
func (s *MetricStore) Record(jobID string, failed bool, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := metricEntry{at: time.Now(), failed: failed, duration: duration}
	s.records[jobID] = append(s.records[jobID], entry)
	if len(s.records[jobID]) > s.maxSize {
		s.records[jobID] = s.records[jobID][len(s.records[jobID])-s.maxSize:]
	}
}

// Snapshot returns aggregated metrics for a job.
func (s *MetricStore) Snapshot(jobID string) (MetricSnapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries, ok := s.records[jobID]
	if !ok || len(entries) == 0 {
		return MetricSnapshot{}, false
	}
	var totalDur time.Duration
	var failCount int
	for _, e := range entries {
		totalDur += e.duration
		if e.failed {
			failCount++
		}
	}
	return MetricSnapshot{
		JobID:       jobID,
		RunCount:    len(entries),
		FailCount:   failCount,
		AvgDuration: totalDur / time.Duration(len(entries)),
		LastRun:     entries[len(entries)-1].at,
	}, true
}

// All returns snapshots for every tracked job.
func (s *MetricStore) All() []MetricSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]MetricSnapshot, 0, len(s.records))
	for jobID, entries := range s.records {
		if len(entries) == 0 {
			continue
		}
		var totalDur time.Duration
		var failCount int
		for _, e := range entries {
			totalDur += e.duration
			if e.failed {
				failCount++
			}
		}
		out = append(out, MetricSnapshot{
			JobID:       jobID,
			RunCount:    len(entries),
			FailCount:   failCount,
			AvgDuration: totalDur / time.Duration(len(entries)),
			LastRun:     entries[len(entries)-1].at,
		})
	}
	return out
}
