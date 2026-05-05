package monitor

import (
	"time"
)

// JobStatus represents the current state of a monitored cron job.
type JobStatus int

const (
	StatusUnknown JobStatus = iota
	StatusOK
	StatusMissed
	StatusFailed
)

func (s JobStatus) String() string {
	switch s {
	case StatusOK:
		return "ok"
	case StatusMissed:
		return "missed"
	case StatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// Job represents a cron job being monitored.
type Job struct {
	ID          string
	Name        string
	Schedule    string        // cron expression, e.g. "0 * * * *"
	GracePeriod time.Duration // allowed delay before marking as missed
	LastCheckin time.Time
	LastStatus  JobStatus
}

// CheckIn records a successful execution of the job.
func (j *Job) CheckIn() {
	j.LastCheckin = time.Now()
	j.LastStatus = StatusOK
}

// Evaluate determines whether the job has missed its expected run window.
// expectedAt is the most recent scheduled time the job should have run.
func (j *Job) Evaluate(expectedAt time.Time) JobStatus {
	deadline := expectedAt.Add(j.GracePeriod)
	if time.Now().After(deadline) && j.LastCheckin.Before(expectedAt) {
		j.LastStatus = StatusMissed
		return StatusMissed
	}
	return j.LastStatus
}
