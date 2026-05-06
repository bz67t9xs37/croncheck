package monitor

import (
	"testing"
	"time"
)

func TestAlertLog_RecordAndAll(t *testing.T) {
	log := NewAlertLog(10)
	log.Record(AlertEntry{JobName: "job1", Status: StatusMissed, FiredAt: time.Now(), Message: "missed"})
	log.Record(AlertEntry{JobName: "job2", Status: StatusMissed, FiredAt: time.Now(), Message: "missed"})

	entries := log.All()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestAlertLog_BoundedSize(t *testing.T) {
	log := NewAlertLog(3)
	for i := 0; i < 5; i++ {
		log.Record(AlertEntry{JobName: "job", Status: StatusMissed, FiredAt: time.Now()})
	}
	if len(log.All()) != 3 {
		t.Fatalf("expected 3 entries after overflow, got %d", len(log.All()))
	}
}

func TestAlertLog_DefaultMaxSize(t *testing.T) {
	log := NewAlertLog(0)
	if log.maxSize != defaultAlertLogSize {
		t.Fatalf("expected default max size %d, got %d", defaultAlertLogSize, log.maxSize)
	}
}

func TestAlertLog_ForJob(t *testing.T) {
	log := NewAlertLog(10)
	log.Record(AlertEntry{JobName: "alpha", Status: StatusMissed, FiredAt: time.Now()})
	log.Record(AlertEntry{JobName: "beta", Status: StatusMissed, FiredAt: time.Now()})
	log.Record(AlertEntry{JobName: "alpha", Status: StatusMissed, FiredAt: time.Now()})

	alpha := log.ForJob("alpha")
	if len(alpha) != 2 {
		t.Fatalf("expected 2 entries for alpha, got %d", len(alpha))
	}

	beta := log.ForJob("beta")
	if len(beta) != 1 {
		t.Fatalf("expected 1 entry for beta, got %d", len(beta))
	}

	unknown := log.ForJob("unknown")
	if len(unknown) != 0 {
		t.Fatalf("expected 0 entries for unknown, got %d", len(unknown))
	}
}
