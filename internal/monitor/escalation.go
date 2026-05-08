package monitor

import (
	"sync"
	"time"
)

// EscalationPolicy defines how alerts should be escalated over time.
type EscalationPolicy struct {
	JobID     string        `json:"job_id"`
	Threshold int           `json:"threshold"` // number of consecutive missed checks before escalation
	Interval  time.Duration `json:"interval"`  // minimum time between escalated alerts
	Webhook   string        `json:"webhook"`   // escalation webhook URL (overrides default)
}

// escalationState tracks per-job escalation state.
type escalationState struct {
	ConsecutiveMisses int
	LastEscalated     time.Time
}

// EscalationStore manages escalation policies and their runtime state.
type EscalationStore struct {
	mu       sync.RWMutex
	policies map[string]EscalationPolicy
	states   map[string]*escalationState
}

// NewEscalationStore creates a new EscalationStore.
func NewEscalationStore() *EscalationStore {
	return &EscalationStore{
		policies: make(map[string]EscalationPolicy),
		states:   make(map[string]*escalationState),
	}
}

// Set registers or updates an escalation policy for a job.
func (s *EscalationStore) Set(policy EscalationPolicy) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.policies[policy.JobID] = policy
}

// Get returns the escalation policy for a job, if any.
func (s *EscalationStore) Get(jobID string) (EscalationPolicy, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.policies[jobID]
	return p, ok
}

// Delete removes the escalation policy and state for a job.
func (s *EscalationStore) Delete(jobID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.policies[jobID]
	if !ok {
		return false
	}
	delete(s.policies, jobID)
	delete(s.states, jobID)
	return true
}

// All returns all registered escalation policies.
func (s *EscalationStore) All() []EscalationPolicy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]EscalationPolicy, 0, len(s.policies))
	for _, p := range s.policies {
		out = append(out, p)
	}
	return out
}

// ShouldEscalate records a miss for the given job and returns true if the
// policy threshold is reached and the escalation interval has elapsed.
func (s *EscalationStore) ShouldEscalate(jobID string, now time.Time) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	policy, ok := s.policies[jobID]
	if !ok {
		return false
	}
	st, exists := s.states[jobID]
	if !exists {
		st = &escalationState{}
		s.states[jobID] = st
	}
	st.ConsecutiveMisses++
	if st.ConsecutiveMisses < policy.Threshold {
		return false
	}
	if !st.LastEscalated.IsZero() && now.Sub(st.LastEscalated) < policy.Interval {
		return false
	}
	st.LastEscalated = now
	return true
}

// ResetMisses resets the consecutive miss counter (e.g. on successful check-in).
func (s *EscalationStore) ResetMisses(jobID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if st, ok := s.states[jobID]; ok {
		st.ConsecutiveMisses = 0
	}
}
