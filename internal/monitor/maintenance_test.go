package monitor

import (
	"testing"
	"time"
)

func TestMaintenanceWindow_IsActive(t *testing.T) {
	now := time.Now()
	w := MaintenanceWindow{
		Start: now.Add(-time.Hour),
		End:   now.Add(time.Hour),
	}
	if !w.IsActive(now) {
		t.Error("expected window to be active")
	}
	if w.IsActive(now.Add(2 * time.Hour)) {
		t.Error("expected window to be inactive after end")
	}
	if w.IsActive(now.Add(-2 * time.Hour)) {
		t.Error("expected window to be inactive before start")
	}
}

func TestMaintenanceStore_IsInMaintenance(t *testing.T) {
	s := NewMaintenanceStore()
	now := time.Now()
	s.Add("backup-job", now.Add(-time.Hour), now.Add(time.Hour))

	if !s.IsInMaintenance("backup-job", now) {
		t.Error("expected job to be in maintenance")
	}
	if s.IsInMaintenance("other-job", now) {
		t.Error("expected other-job to not be in maintenance")
	}
	if s.IsInMaintenance("backup-job", now.Add(2*time.Hour)) {
		t.Error("expected job to be outside maintenance window")
	}
}

func TestMaintenanceStore_All(t *testing.T) {
	s := NewMaintenanceStore()
	now := time.Now()
	s.Add("job-a", now, now.Add(time.Hour))
	s.Add("job-b", now, now.Add(2*time.Hour))

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 windows, got %d", len(all))
	}
}

func TestMaintenanceStore_Remove(t *testing.T) {
	s := NewMaintenanceStore()
	now := time.Now()
	s.Add("job-a", now, now.Add(time.Hour))
	s.Add("job-a", now.Add(time.Hour), now.Add(2*time.Hour))
	s.Add("job-b", now, now.Add(time.Hour))

	removed := s.Remove("job-a")
	if removed != 2 {
		t.Fatalf("expected 2 removed, got %d", removed)
	}
	all := s.All()
	if len(all) != 1 || all[0].JobName != "job-b" {
		t.Errorf("expected only job-b to remain, got %+v", all)
	}
}

func TestMaintenanceStore_Remove_UnknownJob(t *testing.T) {
	s := NewMaintenanceStore()
	now := time.Now()
	s.Add("job-a", now, now.Add(time.Hour))

	removed := s.Remove("nonexistent")
	if removed != 0 {
		t.Errorf("expected 0 removed, got %d", removed)
	}
	if len(s.All()) != 1 {
		t.Error("expected store to remain unchanged")
	}
}
