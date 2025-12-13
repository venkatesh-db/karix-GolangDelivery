package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"project_saas/services/subscription-service/internal/data/migrations"
	"project_saas/services/subscription-service/internal/subscriptions"
	"project_saas/shared/pkg/config"
	"project_saas/shared/pkg/postgres"
	"project_saas/shared/pkg/postgres/migrate"
)

// Register wires subscription catalog endpoints backed by Postgres persistence.
func Register(r chi.Router, cfg config.ServiceConfig, log *zap.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := postgres.Pool(ctx, cfg.PostgresURL, 16)
	if err != nil {
		log.Fatal("failed to connect to postgres", zap.Error(err))
	}
	if err := migrate.Run(ctx, pool, migrations.Files, "."); err != nil {
		log.Fatal("failed to apply migrations", zap.Error(err))
	}
	h := &handler{
		log: log.Named("http"),
		svc: subscriptions.NewService(subscriptions.NewRepository(pool)),
	}
	h.log.Info("subscription routes ready", zap.String("port", cfg.HTTPPort))
	r.Get("/health", health)
	r.Route("/subscriptions", func(r chi.Router) {
		r.Get("/plans", h.listPlans)
		r.Get("/tenants/{tenantID}", h.getSubscription)
		r.Post("/tenants/{tenantID}", h.activatePlan)
	})
}

type handler struct {
	svc *subscriptions.Service
	log *zap.Logger
}

func health(w http.ResponseWriter, _ *http.Request) {
	respond(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *handler) listPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.svc.Plans(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}
	respond(w, http.StatusOK, map[string]interface{}{"plans": plans})
}

type activatePayload struct {
	PlanID string `json:"plan_id"`
	Seats  int    `json:"seats"`
}

func (h *handler) activatePlan(w http.ResponseWriter, r *http.Request) {
	var payload activatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.handleError(w, errBadRequest("invalid json payload"))
		return
	}
	sub, err := h.svc.Activate(r.Context(), subscriptions.ActivateInput{
		TenantID: chi.URLParam(r, "tenantID"),
		PlanID:   payload.PlanID,
		Seats:    payload.Seats,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}
	respond(w, http.StatusAccepted, sub)
}

func (h *handler) getSubscription(w http.ResponseWriter, r *http.Request) {
	sub, err := h.svc.Subscription(r.Context(), chi.URLParam(r, "tenantID"))
	if err != nil {
		h.handleError(w, err)
		return
	}
	respond(w, http.StatusOK, sub)
}

type apiError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func errBadRequest(msg string) error {
	return &apiError{Message: msg, Code: "bad_request"}
}

func (e *apiError) Error() string { return e.Message }

func respond(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, subscriptions.ErrInvalidTenantID),
		errors.Is(err, subscriptions.ErrInvalidPlanID),
		errors.Is(err, subscriptions.ErrInvalidSeats):
		respond(w, http.StatusBadRequest, apiError{Message: err.Error(), Code: "validation"})
	case errors.Is(err, subscriptions.ErrSeatLimit):
		respond(w, http.StatusBadRequest, apiError{Message: err.Error(), Code: "seat_limit"})
	case errors.Is(err, subscriptions.ErrPlanNotFound):
		respond(w, http.StatusNotFound, apiError{Message: err.Error(), Code: "plan_not_found"})
	case errors.Is(err, subscriptions.ErrSubscriptionNotFound):
		respond(w, http.StatusNotFound, apiError{Message: err.Error(), Code: "subscription_not_found"})
	default:
		var apiErr *apiError
		if errors.As(err, &apiErr) {
			respond(w, http.StatusBadRequest, apiErr)
			return
		}
		h.log.Error("request failed", zap.Error(err))
		respond(w, http.StatusInternalServerError, apiError{Message: "internal error", Code: "internal"})
	}
}
