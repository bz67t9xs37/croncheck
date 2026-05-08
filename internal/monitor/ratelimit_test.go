package monitor

import (
	"testing"
	"time"
)

func TestRateLimitStore_Allow_FirstAlert(t *testing.T) {
	store := NewRateLimitStore(time.Minute)
	if !store.Allow("job-a") {
		t.Fatal("expected first alert to be allowed")
	}
}

func TestRateLimitStore_Allow_BlockedWithinCooldown(t *testing.T) {
	store := NewRateLimitStore(time.Minute)
	store.Allow("job-a")
	if store.Allow("job-a") {
		t.Fatal("expected second alert within cooldown to be blocked")
	}
}

func TestRateLimitStore_Allow_PermittedAfterCooldown(t *testing.T) {
	store := NewRateLimitStore(10 * time.Millisecond)
	store.Allow("job-a")
	time.Sleep(20 * time.Millisecond)
	if !store.Allow("job-a") {
		t.Fatal("expected alert to be allowed after cooldown expired")
	}
}

func TestRateLimitStore_Allow_IndependentJobs(t *testing.T) {
	store := NewRateLimitStore(time.Minute)
	store.Allow("job-a")
	if !store.Allow("job-b") {
		t.Fatal("expected different job to be allowed independently")
	}
}

func TestRateLimitStore_Reset(t *testing.T) {
	store := NewRateLimitStore(time.Minute)
	store.Allow("job-a")
	store.Reset("job-a")
	if !store.Allow("job-a") {
		t.Fatal("expected alert to be allowed after reset")
	}
}

func TestRateLimitStore_LastAlert(t *testing.T) {
	store := NewRateLimitStore(time.Minute)
	_, ok := store.LastAlert("job-a")
	if ok {
		t.Fatal("expected no last alert for unknown job")
	}
	before := time.Now()
	store.Allow("job-a")
	last, ok := store.LastAlert("job-a")
	if !ok {
		t.Fatal("expected last alert to be recorded")
	}
	if last.Before(before) {
		t.Errorf("expected last alert time >= %v, got %v", before, last)
	}
}

func TestRateLimitStore_All(t *testing.T) {
	store := NewRateLimitStore(time.Minute)
	store.Allow("job-a")
	store.Allow("job-b")
	all := store.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 records, got %d", len(all))
	}
}

func TestRateLimitStore_DefaultCooldown(t *testing.T) {
	store := NewRateLimitStore(0)
	if store.cooldown != 15*time.Minute {
		t.Errorf("expected default cooldown of 15m, got %v", store.cooldown)
	}
}
