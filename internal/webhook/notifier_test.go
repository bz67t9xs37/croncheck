package webhook_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/croncheck/internal/webhook"
)

func TestSend_Success(t *testing.T) {
	var received webhook.Payload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json content-type, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := webhook.NewNotifier(server.URL)
	p := webhook.Payload{
		JobName:   "backup",
		Alert:     webhook.AlertMissed,
		Message:   "job did not run within expected window",
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}

	if err := n.Send(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.JobName != p.JobName {
		t.Errorf("job_name: got %q, want %q", received.JobName, p.JobName)
	}
	if received.Alert != p.Alert {
		t.Errorf("alert: got %q, want %q", received.Alert, p.Alert)
	}
}

func TestSend_NonSuccessStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := webhook.NewNotifier(server.URL)
	err := n.Send(webhook.Payload{JobName: "test", Alert: webhook.AlertFailed})
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestSend_DefaultTimestamp(t *testing.T) {
	before := time.Now().UTC()
	var received webhook.Payload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	n := webhook.NewNotifier(server.URL)
	_ = n.Send(webhook.Payload{JobName: "ts-test", Alert: webhook.AlertMissed})

	if received.Timestamp.Before(before) {
		t.Error("expected auto-set timestamp to be >= time before send")
	}
}
