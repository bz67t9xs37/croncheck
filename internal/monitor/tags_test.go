package monitor

import (
	"testing"
)

func TestTagStore_SetAndGet(t *testing.T) {
	s := NewTagStore()
	if err := s.Set("job1", "env", "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := s.Set("job1", "team", "platform"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags := s.Get("job1")
	if tags["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", tags["env"])
	}
	if tags["team"] != "platform" {
		t.Errorf("expected team=platform, got %q", tags["team"])
	}
}

func TestTagStore_GetUnknownJob(t *testing.T) {
	s := NewTagStore()
	if tags := s.Get("ghost"); tags != nil {
		t.Errorf("expected nil for unknown job, got %v", tags)
	}
}

func TestTagStore_Delete(t *testing.T) {
	s := NewTagStore()
	_ = s.Set("job1", "env", "staging")
	_ = s.Set("job1", "tier", "critical")
	s.Delete("job1", "env")
	tags := s.Get("job1")
	if _, ok := tags["env"]; ok {
		t.Error("expected env tag to be deleted")
	}
	if tags["tier"] != "critical" {
		t.Errorf("expected tier=critical, got %q", tags["tier"])
	}
}

func TestTagStore_Delete_RemovesJobWhenEmpty(t *testing.T) {
	s := NewTagStore()
	_ = s.Set("job1", "env", "dev")
	s.Delete("job1", "env")
	if tags := s.Get("job1"); tags != nil {
		t.Errorf("expected nil after all tags deleted, got %v", tags)
	}
}

func TestTagStore_All(t *testing.T) {
	s := NewTagStore()
	_ = s.Set("jobA", "env", "prod")
	_ = s.Set("jobB", "env", "dev")
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(all))
	}
	if all["jobA"]["env"] != "prod" {
		t.Errorf("expected jobA env=prod")
	}
}

func TestTagStore_Set_ValidationError(t *testing.T) {
	s := NewTagStore()
	if err := s.Set("", "key", "val"); err == nil {
		t.Error("expected error for empty jobID")
	}
	if err := s.Set("job1", "", "val"); err == nil {
		t.Error("expected error for empty key")
	}
}
