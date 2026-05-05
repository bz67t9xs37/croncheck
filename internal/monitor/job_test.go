package monitor

import (
	"testing"
	"time"
)

func TestJobStatus_String(t *testing.T) {
	cases := []struct {
		status   JobStatus
		expected string
	}{
		{StatusOK, "ok"},
		{StatusMissed, "missed"},
		{StatusFailed, "failed"},
		{StatusUnknown, "unknown"},
	}
	for _, tc := range cases {
		if got := tc.status.String(); got != tc.expected {
			t.Errorf("String() = %q, want %q", got, tc.expected)
		}
	}
}

func TestJob_CheckIn(t *testing.T) {
	j := &Job{ID: "job-1", Name: "backup"}
	before := time.Now()
	j.CheckIn()
	after := time.Now()

	if j.LastStatus != StatusOK {
		t.Errorf("expected StatusOK after CheckIn, got %v", j.LastStatus)
	}
	if j.LastCheckin.Before(before) || j.LastCheckin.After(after) {
		t.Errorf("LastCheckin %v not within expected range", j.LastCheckin)
	}
}

func TestJob_Evaluate_Missed(t *testing.T) {
	j := &Job{
		ID:          "job-2",
		Name:        "report",
		GracePeriod: 5 * time.Minute,
		LastCheckin: time.Now().Add(-2 * time.Hour),
	}
	expectedAt := time.Now().Add(-1 * time.Hour)
	status := j.Evaluate(expectedAt)
	if status != StatusMissed {
		t.Errorf("expected StatusMissed, got %v", status)
	}
}

func TestJob_Evaluate_WithinGrace(t *testing.T) {
	j := &Job{
		ID:          "job-3",
		Name:        "cleanup",
		GracePeriod: 10 * time.Minute,
		LastCheckin: time.Now().Add(-30 * time.Second),
		LastStatus:  StatusOK,
	}
	// expected run was 1 minute ago, grace is 10 min — not yet missed
	expectedAt := time.Now().Add(-1 * time.Minute)
	status := j.Evaluate(expectedAt)
	if status != StatusOK {
		t.Errorf("expected StatusOK within grace period, got %v", status)
	}
}
