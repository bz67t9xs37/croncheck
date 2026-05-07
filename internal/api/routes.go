package api

import (
	"net/http"

	"github.com/croncheck/internal/monitor"
	"github.com/croncheck/internal/webhook"
)

// RouterDeps groups all dependencies needed to register API routes.
type RouterDeps struct {
	Registry     *monitor.Registry
	History      *monitor.History
	AlertLog     *monitor.AlertLog
	Scheduler    *monitor.Scheduler
	Silences     *monitor.SilenceStore
	Maintenance  *monitor.MaintenanceStore
	Retry        *monitor.RetryTracker
	Tags         *monitor.TagStore
	Dependencies *monitor.DependencyStore
	Notifier     *webhook.Notifier
}

// RegisterRoutes attaches all API handlers to mux.
func RegisterRoutes(mux *http.ServeMux, deps RouterDeps) {
	mux.Handle("/checkin/", NewHandler(deps.Registry, deps.History))
	mux.Handle("/status", NewHandler(deps.Registry, deps.History))
	mux.Handle("/history", NewHistoryHandler(deps.History))
	mux.Handle("/alerts", NewAlertLogHandler(deps.AlertLog))
	mux.Handle("/scheduler", NewSchedulerHandler(deps.Scheduler))
	mux.Handle("/silences", NewSilenceHandler(deps.Silences))
	mux.Handle("/maintenance", NewMaintenanceHandler(deps.Maintenance))
	mux.Handle("/retry", NewRetryHandler(deps.Retry))
	mux.Handle("/tags", NewTagsHandler(deps.Tags))
	mux.Handle("/dependencies", NewDependencyHandler(deps.Dependencies))
}
