package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"project_saas/services/billing-service/internal/engine"
	"project_saas/shared/pkg/config"
)

// Register exposes billing aggregation endpoints.
func Register(r chi.Router, cfg config.ServiceConfig, log *zap.Logger) {
	proc := engine.NewProcessor(cfg, log)
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		respond(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	r.Route("/billing", func(r chi.Router) {
		r.Post("/tenants/{tenantID}/run", func(w http.ResponseWriter, req *http.Request) {
			countParam := req.URL.Query().Get("records")
			total := 1_000_000
			if countParam != "" {
				if parsed, err := strconv.Atoi(countParam); err == nil {
					total = parsed
				}
			}
			result, err := proc.Run(req.Context(), chi.URLParam(req, "tenantID"), total)
			if err != nil {
				log.Error("billing run failed", zap.Error(err))
				respond(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			respond(w, http.StatusAccepted, result)
		})
	})
}

func respond(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
