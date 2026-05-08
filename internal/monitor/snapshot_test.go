package monitor

import (
	"testing"
	"time"
)

func TestSnapshotStore_RecordAndGet(t *testing.T) {
	store := NewSnapshotStore()

	now := time.Now()
	snap := JobSnapshot{
		JobID:       "job-a",
		Status:      StatusHealthy,
		LastCheckIn: &now,
	}
	store.Record(snap)

	got, ok := store.Get("job-a")
	if !ok {
		t.Fatal("expected snapshot to exist")
	}
	if got.JobID != "job-a" {
		t.Errorf("expected job_id=job-a, got %s", got.JobID)
	}
	if got.Status != StatusHealthy {
		t.Errorf("expected status healthy, got %s", got.Status)
	}
	if got.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set")
	}
}

func TestSnapshotStore_GetUnknownJob(t *testing.T) {
	store := NewSnapshotStore()
	_, ok := store.Get("nonexistent")
	if ok {
		t.Error("expected no snapshot for unknown job")
	}
}

func TestSnapshotStore_All(t *testing.T) {
	store := NewSnapshotStore()
	store.Record(JobSnapshot{JobID: "job-1", Status: StatusHealthy})
	store.Record(JobSnapshot{JobID: "job-2", Status: StatusMissed})

	all := store.All()
	if len(all) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(all))
	}
}

func TestSnapshotStore_Delete(t *testing.T) {
	store := NewSnapshotStore()
	store.Record(JobSnapshot{JobID: "job-x", Status: StatusHealthy})
	store.Delete("job-x")

	_, ok := store.Get("job-x")
	if ok {
		t.Error("expected snapshot to be deleted")
	}
}

func TestSnapshotStore_OverwritesExisting(t *testing.T) {
	store := NewSnapshotStore()
	store.Record(JobSnapshot{JobID: "job-z", Status: StatusHealthy})
	store.Record(JobSnapshot{JobID: "job-z", Status: StatusMissed})

	got, ok := store.Get("job-z")
	if !ok {
		t.Fatal("expected snapshot to exist")
	}
	if got.Status != StatusMissed {
		t.Errorf("expected overwritten status=missed, got %s", got.Status)
	}
}
