package api

import (
	"net/http"

	"github.com/example/croncheck/internal/monitor"
)

// RegisterRoutes wires all API handlers onto the given mux.
func RegisterRoutes(
	mux *http.ServeMux,
	reg *monitor.Registry,
	history *monitor.History,
	alertLog *monitor.AlertLog,
	scheduler *monitor.Scheduler,
	silences *monitor.SilenceStore,
	maintenance *monitor.MaintenanceStore,
	retryTracker *monitor.RetryTracker,
	tags *monitor.TagStore,
	deps *monitor.DependencyStore,
	heartbeats *monitor.HeartbeatStore,
	snapshots *monitor.SnapshotStore,
	rateLimits *monitor.RateLimitStore,
	runbooks *monitor.RunbookStore,
	escalations *monitor.EscalationStore,
	annotations *monitor.AnnotationStore,
	metrics *monitor.MetricStore,
	sla *monitor.SLAStore,
) {
	mux.Handle("/checkin", NewHandler(reg, history))
	mux.Handle("/status", NewHandler(reg, history))
	mux.Handle("/history", NewHistoryHandler(history))
	mux.Handle("/alerts", NewAlertLogHandler(alertLog))
	mux.Handle("/scheduler", NewSchedulerHandler(scheduler))
	mux.Handle("/silences", NewSilenceHandler(silences))
	mux.Handle("/maintenance", NewMaintenanceHandler(maintenance))
	mux.Handle("/retry", NewRetryHandler(retryTracker))
	mux.Handle("/tags", NewTagsHandler(tags))
	mux.Handle("/dependencies", NewDependencyHandler(deps))
	mux.Handle("/snapshots", NewSnapshotHandler(snapshots))
	mux.Handle("/ratelimits", NewRateLimitHandler(rateLimits))
	mux.Handle("/runbooks", NewRunbookHandler(runbooks))
	mux.Handle("/escalations", NewEscalationHandler(escalations))
	mux.Handle("/annotations", NewAnnotationHandler(annotations))
	mux.Handle("/metrics", NewMetricHandler(metrics))
	mux.Handle("/sla", NewSLAHandler(sla))
	mux.Handle("/sla/violations", NewSLAHandler(sla))
}
