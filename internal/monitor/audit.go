package monitor

import (
	"sync"
	"time"
)

const defaultMaxAuditEntries = 200

// AuditAction represents the type of action recorded in the audit log.
type AuditAction string

const (
	AuditActionCreated  AuditAction = "created"
	AuditActionUpdated  AuditAction = "updated"
	AuditActionDeleted  AuditAction = "deleted"
	AuditActionCheckIn  AuditAction = "check_in"
	AuditActionSilenced AuditAction = "silenced"
	AuditActionAlertion AuditAction = "alerted"
)

// AuditEntry records a single auditable event for a job.
type AuditEntry struct {
	JobID     string      `json:"job_id"`
	Action    AuditAction `json:"action"`
	Actor     string      `json:"actor,omitempty"`
	Detail    string      `json:"detail,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// AuditStore maintains a bounded, in-memory audit log of job-related events.
type AuditStore struct {
	mu      sync.Mutex
	entries []AuditEntry
	maxSize int
}

// NewAuditStore creates an AuditStore with the given maximum capacity.
// If maxSize <= 0 the default is used.
func NewAuditStore(maxSize int) *AuditStore {
	if maxSize <= 0 {
		maxSize = defaultMaxAuditEntries
	}
	return &AuditStore{maxSize: maxSize}
}

// Record appends an audit entry, evicting the oldest when the store is full.
func (s *AuditStore) Record(entry AuditEntry) {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.entries) >= s.maxSize {
		s.entries = s.entries[1:]
	}
	s.entries = append(s.entries, entry)
}

// All returns a copy of all audit entries in chronological order.
func (s *AuditStore) All() []AuditEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]AuditEntry, len(s.entries))
	copy(out, s.entries)
	return out
}

// ForJob returns audit entries for a specific job ID.
func (s *AuditStore) ForJob(jobID string) []AuditEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []AuditEntry
	for _, e := range s.entries {
		if e.JobID == jobID {
			out = append(out, e)
		}
	}
	return out
}
