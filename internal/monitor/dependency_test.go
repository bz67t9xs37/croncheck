package monitor

import (
	"testing"
)

func TestDependencyStore_AddAndQuery(t *testing.T) {
	ds := NewDependencyStore()
	if err := ds.Add("job-a", "job-b"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ups := ds.UpstreamsOf("job-b")
	if len(ups) != 1 || ups[0] != "job-a" {
		t.Errorf("expected [job-a], got %v", ups)
	}
	downs := ds.DownstreamsOf("job-a")
	if len(downs) != 1 || downs[0] != "job-b" {
		t.Errorf("expected [job-b], got %v", downs)
	}
}

func TestDependencyStore_DuplicateRejected(t *testing.T) {
	ds := NewDependencyStore()
	_ = ds.Add("job-a", "job-b")
	if err := ds.Add("job-a", "job-b"); err == nil {
		t.Error("expected error for duplicate dependency")
	}
}

func TestDependencyStore_SelfDependencyRejected(t *testing.T) {
	ds := NewDependencyStore()
	if err := ds.Add("job-a", "job-a"); err == nil {
		t.Error("expected error for self-dependency")
	}
}

func TestDependencyStore_Remove(t *testing.T) {
	ds := NewDependencyStore()
	_ = ds.Add("job-a", "job-b")
	if err := ds.Remove("job-a", "job-b"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ds.All()) != 0 {
		t.Error("expected empty store after remove")
	}
}

func TestDependencyStore_Remove_NotFound(t *testing.T) {
	ds := NewDependencyStore()
	if err := ds.Remove("job-x", "job-y"); err == nil {
		t.Error("expected error when removing non-existent dependency")
	}
}

func TestDependencyStore_All(t *testing.T) {
	ds := NewDependencyStore()
	_ = ds.Add("job-a", "job-b")
	_ = ds.Add("job-a", "job-c")
	all := ds.All()
	if len(all) != 2 {
		t.Errorf("expected 2 links, got %d", len(all))
	}
}

func TestDependencyStore_UnknownJob_ReturnsEmpty(t *testing.T) {
	ds := NewDependencyStore()
	if ups := ds.UpstreamsOf("ghost"); len(ups) != 0 {
		t.Errorf("expected empty, got %v", ups)
	}
	if downs := ds.DownstreamsOf("ghost"); len(downs) != 0 {
		t.Errorf("expected empty, got %v", downs)
	}
}
