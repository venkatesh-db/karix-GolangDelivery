package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"project_saas/shared/pkg/config"
)

// Register exposes notification fan-out endpoints.
func Register(r chi.Router, cfg config.ServiceConfig, log *zap.Logger) {
	log.Named("http").Info("notification routes ready", zap.String("port", cfg.HTTPPort))
	r.Get("/health", health)
	r.Route("/notifications", func(r chi.Router) {
		r.Post("/tenants/{tenantID}", enqueueNotification)
	})
}

func health(w http.ResponseWriter, _ *http.Request) {
	respond(w, http.StatusOK, map[string]string{"status": "ok"})
}

func enqueueNotification(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusAccepted, map[string]string{"tenant": chi.URLParam(r, "tenantID"), "notification": "queued"})
}

func respond(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
