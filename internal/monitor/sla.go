package monitor

import (
	"fmt"
	"sync"
	"time"
)

// SLAPolicy defines the expected success rate and max allowed downtime for a job.
type SLAPolicy struct {
	JobID          string        `json:"job_id"`
	MinSuccessRate float64       `json:"min_success_rate"` // 0.0 - 1.0
	MaxDowntime    time.Duration `json:"max_downtime"`
	CreatedAt      time.Time     `json:"created_at"`
}

// SLAViolation records a detected SLA breach.
type SLAViolation struct {
	JobID       string    `json:"job_id"`
	Reason      string    `json:"reason"`
	DetectedAt  time.Time `json:"detected_at"`
}

// SLAStore manages SLA policies and tracks violations.
type SLAStore struct {
	mu         sync.RWMutex
	policies   map[string]SLAPolicy
	violations []SLAViolation
	maxViolations int
}

const defaultMaxViolations = 200

// NewSLAStore returns an initialised SLAStore.
func NewSLAStore(maxViolations int) *SLAStore {
	if maxViolations <= 0 {
		maxViolations = defaultMaxViolations
	}
	return &SLAStore{
		policies:      make(map[string]SLAPolicy),
		maxViolations: maxViolations,
	}
}

// Set adds or replaces an SLA policy for a job.
func (s *SLAStore) Set(p SLAPolicy) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.CreatedAt = time.Now()
	s.policies[p.JobID] = p
}

// Get returns the SLA policy for a job, if one exists.
func (s *SLAStore) Get(jobID string) (SLAPolicy, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.policies[jobID]
	return p, ok
}

// Delete removes the SLA policy for a job.
func (s *SLAStore) Delete(jobID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.policies[jobID]
	delete(s.policies, jobID)
	return ok
}

// All returns all current SLA policies.
func (s *SLAStore) All() []SLAPolicy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]SLAPolicy, 0, len(s.policies))
	for _, p := range s.policies {
		out = append(out, p)
	}
	return out
}

// RecordViolation appends a new SLA violation, evicting the oldest if at capacity.
func (s *SLAStore) RecordViolation(jobID, reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v := SLAViolation{JobID: jobID, Reason: reason, DetectedAt: time.Now()}
	if len(s.violations) >= s.maxViolations {
		s.violations = s.violations[1:]
	}
	s.violations = append(s.violations, v)
}

// Violations returns all recorded violations, optionally filtered by jobID.
func (s *SLAStore) Violations(jobID string) []SLAViolation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if jobID == "" {
		out := make([]SLAViolation, len(s.violations))
		copy(out, s.violations)
		return out
	}
	var out []SLAViolation
	for _, v := range s.violations {
		if v.JobID == jobID {
			out = append(out, v)
		}
	}
	return out
}

// Evaluate checks a job's metrics against its SLA policy and records any violations.
func (s *SLAStore) Evaluate(jobID string, metrics JobMetrics) {
	s.mu.RLock()
	p, ok := s.policies[jobID]
	s.mu.RUnlock()
	if !ok {
		return
	}
	total := metrics.SuccessCount + metrics.FailCount
	if total > 0 {
		rate := float64(metrics.SuccessCount) / float64(total)
		if rate < p.MinSuccessRate {
			s.RecordViolation(jobID, fmt.Sprintf("success rate %.2f below threshold %.2f", rate, p.MinSuccessRate))
		}
	}
	if p.MaxDowntime > 0 && metrics.LastDowntime > p.MaxDowntime {
		s.RecordViolation(jobID, fmt.Sprintf("downtime %s exceeds max %s", metrics.LastDowntime, p.MaxDowntime))
	}
}

// JobMetrics is a minimal snapshot used for SLA evaluation.
type JobMetrics struct {
	SuccessCount int
	FailCount    int
	LastDowntime time.Duration
}
