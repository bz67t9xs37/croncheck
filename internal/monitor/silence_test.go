package monitor

import (
	"testing"
	"time"
)

func TestSilenceStore_IsSilenced_Active(t *testing.T) {
	s := NewSilenceStore()
	now := time.Now()
	s.Silence("backup", now.Add(10*time.Minute), "maintenance")

	if !s.IsSilenced("backup", now) {
		t.Error("expected job to be silenced")
	}
}

func TestSilenceStore_IsSilenced_Expired(t *testing.T) {
	s := NewSilenceStore()
	now := time.Now()
	s.Silence("backup", now.Add(-1*time.Minute), "old window")

	if s.IsSilenced("backup", now) {
		t.Error("expected silence to be expired")
	}
}

func TestSilenceStore_IsSilenced_UnknownJob(t *testing.T) {
	s := NewSilenceStore()
	if s.IsSilenced("nonexistent", time.Now()) {
		t.Error("expected unknown job to not be silenced")
	}
}

func TestSilenceStore_Unsilence(t *testing.T) {
	s := NewSilenceStore()
	now := time.Now()
	s.Silence("deploy", now.Add(1*time.Hour), "deploy window")
	s.Unsilence("deploy")

	if s.IsSilenced("deploy", now) {
		t.Error("expected job to be unsilenced after removal")
	}
}

func TestSilenceStore_All(t *testing.T) {
	s := NewSilenceStore()
	now := time.Now()
	s.Silence("job-a", now.Add(1*time.Hour), "reason a")
	s.Silence("job-b", now.Add(2*time.Hour), "reason b")

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestSilenceStore_OverwriteSilence(t *testing.T) {
	s := NewSilenceStore()
	now := time.Now()
	s.Silence("job-a", now.Add(-1*time.Minute), "expired")
	s.Silence("job-a", now.Add(1*time.Hour), "extended")

	if !s.IsSilenced("job-a", now) {
		t.Error("expected overwritten silence to be active")
	}
}
