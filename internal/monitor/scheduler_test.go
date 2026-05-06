package monitor

import (
	"sync"
	"testing"
	"time"
)

type mockNotifier struct {
	mu   sync.Mutex
	calls []string
}

func (m *mockNotifier) Send(jobName, status, message string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, jobName+":"+status)
	return nil
}

func (m *mockNotifier) CallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.calls)
}

func newSchedulerDeps(t *testing.T) (*Registry, *History, *AlertLog) {
	t.Helper()
	reg := NewRegistry()
	hist := NewHistory(0)
	alog := NewAlertLog(0)
	return reg, hist, alog
}

func TestScheduler_StartStop(t *testing.T) {
	reg, hist, alog := newSchedulerDeps(t)
	n := &mockNotifier{}
	s := NewScheduler(reg, hist, alog, n, 50*time.Millisecond)
	s.Start()
	time.Sleep(20 * time.Millisecond)
	s.Stop() // should not block or panic
}

func TestScheduler_AlertsFiredForMissedJob(t *testing.T) {
	reg, hist, alog := newSchedulerDeps(t)
	n := &mockNotifier{}

	reg.Register(Job{
		Name:        "nightly",
		Schedule:    "@daily",
		GracePeriod: time.Second,
		LastCheckIn: time.Now().Add(-25 * time.Hour),
		ExpectedEvery: 24 * time.Hour,
	})

	s := NewScheduler(reg, hist, alog, n, 30*time.Millisecond)
	s.Start()
	time.Sleep(80 * time.Millisecond)
	s.Stop()

	if n.CallCount() == 0 {
		t.Error("expected at least one alert notification for missed job")
	}

	alerts := alog.All()
	if len(alerts) == 0 {
		t.Error("expected alert log entries for missed job")
	}
}

func TestScheduler_NoAlertForHealthyJob(t *testing.T) {
	reg, hist, alog := newSchedulerDeps(t)
	n := &mockNotifier{}

	reg.Register(Job{
		Name:          "frequent",
		Schedule:      "@hourly",
		GracePeriod:   5 * time.Minute,
		LastCheckIn:   time.Now().Add(-30 * time.Minute),
		ExpectedEvery: time.Hour,
	})

	s := NewScheduler(reg, hist, alog, n, 30*time.Millisecond)
	s.Start()
	time.Sleep(80 * time.Millisecond)
	s.Stop()

	if n.CallCount() != 0 {
		t.Errorf("expected no alerts for healthy job, got %d", n.CallCount())
	}
	_ = hist
}
