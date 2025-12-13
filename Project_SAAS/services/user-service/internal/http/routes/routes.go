package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"project_saas/services/user-service/internal/data/migrations"
	"project_saas/services/user-service/internal/users"
	"project_saas/shared/pkg/config"
	"project_saas/shared/pkg/postgres"
	"project_saas/shared/pkg/postgres/migrate"
)

// Register wires user-service HTTP endpoints.
func Register(r chi.Router, cfg config.ServiceConfig, log *zap.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := postgres.Pool(ctx, cfg.PostgresURL, 32)
	if err != nil {
		log.Fatal("failed to connect to postgres", zap.Error(err))
	}
	if err := migrate.Run(ctx, pool, migrations.Files, "."); err != nil {
		log.Fatal("failed to apply migrations", zap.Error(err))
	}
	h := &handler{
		log: log.Named("http"),
		svc: users.NewService(users.NewRepository(pool)),
	}
	h.log.Info("registering routes", zap.String("port", cfg.HTTPPort))
	r.Get("/health", health)
	r.Route("/tenants/{tenantID}", func(r chi.Router) {
		r.Get("/users", h.listUsers)
		r.Post("/users", h.createUser)
	})
}

type handler struct {
	svc *users.Service
	log *zap.Logger
}

func health(w http.ResponseWriter, _ *http.Request) {
	respond(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *handler) listUsers(w http.ResponseWriter, r *http.Request) {
	tenantID := chi.URLParam(r, "tenantID")
	result, err := h.svc.List(r.Context(), tenantID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	respond(w, http.StatusOK, map[string]interface{}{
		"tenant": tenantID,
		"users":  result,
	})
}

type createUserPayload struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

func (h *handler) createUser(w http.ResponseWriter, r *http.Request) {
	var payload createUserPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.handleError(w, errBadRequest("invalid json"))
		return
	}
	user, err := h.svc.Create(r.Context(), users.CreateInput{
		TenantID: chi.URLParam(r, "tenantID"),
		Email:    payload.Email,
		FullName: payload.FullName,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}
	respond(w, http.StatusCreated, user)
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
	case errors.Is(err, users.ErrInvalidTenant), errors.Is(err, users.ErrInvalidEmail), errors.Is(err, users.ErrInvalidName):
		respond(w, http.StatusBadRequest, apiError{Message: err.Error(), Code: "validation"})
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
