package api

import (
	"net/http"

	"github.com/croncheck/internal/monitor"
)

// RegisterRoutes wires all API handlers to the provided ServeMux.
func RegisterRoutes(
	mux *http.ServeMux,
	reg *monitor.Registry,
	hist *monitor.History,
	alertLog *monitor.AlertLog,
	sched *monitor.Scheduler,
	silence *monitor.SilenceStore,
	maint *monitor.MaintenanceStore,
	retryTracker *monitor.RetryTracker,
	tags *monitor.TagStore,
	deps *monitor.DependencyStore,
	heartbeat *monitor.HeartbeatStore,
	snapshot *monitor.SnapshotStore,
	rateLimit *monitor.RateLimitStore,
	runbook *monitor.RunbookStore,
	escalation *monitor.EscalationStore,
	annotation *monitor.AnnotationStore,
	metric *monitor.MetricStore,
	sla *monitor.SLAStore,
	ownership *monitor.OwnershipStore,
) {
	mux.Handle("/checkin", NewHandler(reg, hist))
	mux.Handle("/status", NewHandler(reg, hist))
	mux.Handle("/history", NewHistoryHandler(hist))
	mux.Handle("/alerts", NewAlertLogHandler(alertLog))
	mux.Handle("/scheduler", NewSchedulerHandler(sched))
	mux.Handle("/silence", NewSilenceHandler(silence))
	mux.Handle("/maintenance", NewMaintenanceHandler(maint))
	mux.Handle("/retry", NewRetryHandler(retryTracker))
	mux.Handle("/tags", NewTagsHandler(tags))
	mux.Handle("/dependencies", NewDependencyHandler(deps))
	mux.Handle("/snapshots", NewSnapshotHandler(snapshot))
	mux.Handle("/ratelimit", NewRateLimitHandler(rateLimit))
	mux.Handle("/runbooks", NewRunbookHandler(runbook))
	mux.Handle("/escalations", NewEscalationHandler(escalation))
	mux.Handle("/annotations", NewAnnotationHandler(annotation))
	mux.Handle("/metrics", NewMetricHandler(metric))
	mux.Handle("/sla", NewSLAHandler(sla))
	mux.Handle("/ownership", NewOwnershipHandler(ownership))
}
