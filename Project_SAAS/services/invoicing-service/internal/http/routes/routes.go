package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"project_saas/shared/pkg/config"
)

// Register exposes invoice generation endpoints.
func Register(r chi.Router, cfg config.ServiceConfig, log *zap.Logger) {
	log.Named("http").Info("invoicing routes ready", zap.String("port", cfg.HTTPPort))
	r.Get("/health", health)
	r.Route("/invoices", func(r chi.Router) {
		r.Post("/tenants/{tenantID}/generate", generateInvoice)
	})
}

func health(w http.ResponseWriter, _ *http.Request) {
	respond(w, http.StatusOK, map[string]string{"status": "ok"})
}

func generateInvoice(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusAccepted, map[string]string{"tenant": chi.URLParam(r, "tenantID"), "invoice": "queued"})
}

func respond(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
