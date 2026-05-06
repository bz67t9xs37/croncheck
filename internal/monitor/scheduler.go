package monitor

import (
	"sync"
	"time"
)

// Scheduler tracks expected check-in intervals and marks jobs as missed
// when they exceed their schedule plus grace period.
type Scheduler struct {
	registry *Registry
	history  *History
	alertLog *AlertLog
	notifier Notifier
	ticker   *time.Ticker
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

// Notifier is a small interface so Scheduler doesn't import the webhook package directly.
type Notifier interface {
	Send(jobName, status, message string) error
}

// NewScheduler creates a Scheduler that evaluates jobs on the given interval.
func NewScheduler(registry *Registry, history *History, alertLog *AlertLog, notifier Notifier, interval time.Duration) *Scheduler {
	return &Scheduler{
		registry: registry,
		history:  history,
		alertLog: alertLog,
		notifier: notifier,
		ticker:   time.NewTicker(interval),
		stopCh:   make(chan struct{}),
	}
}

// Start begins the background evaluation loop.
func (s *Scheduler) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.ticker.C:
				s.evaluate()
			case <-s.stopCh:
				return
			}
		}
	}()
}

// Stop gracefully shuts down the scheduler.
func (s *Scheduler) Stop() {
	s.ticker.Stop()
	close(s.stopCh)
	s.wg.Wait()
}

func (s *Scheduler) evaluate() {
	now := time.Now()
	for _, job := range s.registry.All() {
		if job.Evaluate(now) == StatusMissed {
			s.alertLog.Record(AlertEntry{
				JobName:   job.Name,
				Status:    StatusMissed,
				Timestamp: now,
				Message:   "job missed expected check-in",
			})
			if s.notifier != nil {
				_ = s.notifier.Send(job.Name, StatusMissed.String(), "job missed expected check-in")
			}
		}
	}
}
