package monitor

import (
	"testing"
)

func TestOwnershipStore_SetAndGet(t *testing.T) {
	s := NewOwnershipStore()
	o := Owner{Job: "backup", Name: "Alice", Email: "alice@example.com", Team: "ops"}
	if err := s.Set(o); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := s.Get("backup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != o {
		t.Errorf("expected %+v, got %+v", o, got)
	}
}

func TestOwnershipStore_GetUnknownJob(t *testing.T) {
	s := NewOwnershipStore()
	_, err := s.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown job, got nil")
	}
}

func TestOwnershipStore_Set_EmptyJob(t *testing.T) {
	s := NewOwnershipStore()
	err := s.Set(Owner{Job: "", Email: "x@example.com"})
	if err == nil {
		t.Fatal("expected error for empty job name")
	}
}

func TestOwnershipStore_Set_EmptyEmail(t *testing.T) {
	s := NewOwnershipStore()
	err := s.Set(Owner{Job: "sync", Email: ""})
	if err == nil {
		t.Fatal("expected error for empty email")
	}
}

func TestOwnershipStore_Delete(t *testing.T) {
	s := NewOwnershipStore()
	_ = s.Set(Owner{Job: "cleanup", Name: "Bob", Email: "bob@example.com"})
	if err := s.Delete("cleanup"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := s.Get("cleanup")
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

func TestOwnershipStore_Delete_NotFound(t *testing.T) {
	s := NewOwnershipStore()
	if err := s.Delete("ghost"); err == nil {
		t.Fatal("expected error when deleting unknown job")
	}
}

func TestOwnershipStore_All(t *testing.T) {
	s := NewOwnershipStore()
	_ = s.Set(Owner{Job: "job1", Email: "a@example.com"})
	_ = s.Set(Owner{Job: "job2", Email: "b@example.com"})
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 owners, got %d", len(all))
	}
}
