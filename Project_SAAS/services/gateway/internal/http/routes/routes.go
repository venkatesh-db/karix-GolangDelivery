package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"project_saas/shared/pkg/auth"
	"project_saas/shared/pkg/config"
)

// Register sets up API gateway routes used by external clients.
func Register(r chi.Router, cfg config.ServiceConfig, log *zap.Logger) {
	logger := log.Named("http")
	logger.Info("gateway ready", zap.String("port", cfg.HTTPPort))
	validator := auth.NewValidator(cfg.AuthSecret)
	mw := auth.Middleware(validator, log.Named("auth"))
	r.Get("/health", health)
	r.Route("/api", func(r chi.Router) {
		r.Use(mw)
		r.Get("/status", aggregateStatus)
		r.Get("/me", me)
	})
}

func health(w http.ResponseWriter, _ *http.Request) {
	respond(w, http.StatusOK, map[string]string{"status": "ok"})
}

func aggregateStatus(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	payload := map[string]interface{}{"user": "ok", "billing": "ok"}
	if claims != nil {
		payload["tenant_id"] = claims.TenantID
		payload["roles"] = claims.Roles
	}
	respond(w, http.StatusOK, payload)
}

func me(w http.ResponseWriter, r *http.Request) {
	if claims, ok := auth.ClaimsFromContext(r.Context()); ok {
		respond(w, http.StatusOK, claims)
		return
	}
	respond(w, http.StatusInternalServerError, map[string]string{"error": "claims missing"})
}

func respond(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
