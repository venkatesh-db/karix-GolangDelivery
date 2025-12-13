package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"log/slog"

	"github.com/google/uuid"

	appcustomer "github.com/helrachar/banking/internal/application/customer"
	"github.com/helrachar/banking/internal/domain"
	"github.com/helrachar/banking/internal/server"
)

const handlerTimeout = 3 * time.Second

// CustomerHandler exposes HTTP endpoints for customer-centric use cases.
type CustomerHandler struct {
	service *appcustomer.Service
	logger  *slog.Logger
}

// NewCustomerHandler wires dependencies for HTTP delivery.
func NewCustomerHandler(service *appcustomer.Service, logger *slog.Logger) *CustomerHandler {
	return &CustomerHandler{service: service, logger: logger}
}

// Register attaches endpoints to the provided mux.
func (h *CustomerHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", h.health)
	mux.HandleFunc("/customers", h.routeCustomers)
	mux.HandleFunc("/customers/", h.routeCustomerByID)
}

func (h *CustomerHandler) health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok", "requestId": server.RequestIDFromContext(r.Context())})
}

func (h *CustomerHandler) routeCustomers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createCustomer(w, r)
	case http.MethodGet:
		h.listCustomers(w, r)
	default:
		respondError(w, r, http.StatusMethodNotAllowed, "method not allowed", nil)
	}
}

func (h *CustomerHandler) routeCustomerByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, r, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/customers/")
	if id == "" {
		respondError(w, r, http.StatusNotFound, "customer not found", nil)
		return
	}
	h.getCustomer(w, r, domain.CustomerID(id))
}

type createCustomerRequest struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	PAN      string `json:"pan"`
}

type customerResponse struct {
	ID        string `json:"id"`
	FullName  string `json:"fullName"`
	Email     string `json:"email"`
	PAN       string `json:"pan"`
	CreatedAt string `json:"createdAt"`
}

type errorResponse struct {
	Error     string            `json:"error"`
	Details   map[string]string `json:"details,omitempty"`
	RequestID string            `json:"requestId"`
}

func (h *CustomerHandler) createCustomer(w http.ResponseWriter, r *http.Request) {
	var req createCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid payload", "error", err, "request_id", server.RequestIDFromContext(r.Context()))
		respondError(w, r, http.StatusBadRequest, "invalid payload", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), handlerTimeout)
	defer cancel()

	id := domain.CustomerID(uuid.NewString())
	customer, err := h.service.RegisterCustomer(ctx, id, req.FullName, req.Email, req.PAN)
	if err != nil {
		var validationErr *domain.ValidationError
		if errors.As(err, &validationErr) {
			respondError(w, r, http.StatusBadRequest, "validation failed", validationErr.Fields())
			return
		}
		respondError(w, r, http.StatusInternalServerError, "unable to create customer", nil)
		h.logger.Error("register customer failed", "error", err, "request_id", server.RequestIDFromContext(r.Context()))
		return
	}

	respondJSON(w, http.StatusCreated, mapCustomer(customer))
}

func (h *CustomerHandler) getCustomer(w http.ResponseWriter, r *http.Request, id domain.CustomerID) {
	ctx, cancel := context.WithTimeout(r.Context(), handlerTimeout)
	defer cancel()

	customer, err := h.service.FetchCustomer(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, r, http.StatusNotFound, "customer not found", nil)
			return
		}
		respondError(w, r, http.StatusInternalServerError, "unable to fetch customer", nil)
		h.logger.Error("fetch customer failed", "error", err, "request_id", server.RequestIDFromContext(r.Context()), "customer_id", id)
		return
	}

	respondJSON(w, http.StatusOK, mapCustomer(customer))
}

func (h *CustomerHandler) listCustomers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), handlerTimeout)
	defer cancel()

	customers, err := h.service.ListCustomers(ctx)
	if err != nil {
		respondError(w, r, http.StatusInternalServerError, "unable to list customers", nil)
		h.logger.Error("list customers failed", "error", err, "request_id", server.RequestIDFromContext(r.Context()))
		return
	}

	resp := make([]customerResponse, 0, len(customers))
	for _, c := range customers {
		resp = append(resp, mapCustomer(c))
	}
	respondJSON(w, http.StatusOK, resp)
}

func mapCustomer(c *domain.Customer) customerResponse {
	return customerResponse{
		ID:        string(c.ID()),
		FullName:  c.FullName(),
		Email:     c.Email(),
		PAN:       c.PAN(),
		CreatedAt: c.CreatedAt().Format(time.RFC3339),
	}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, r *http.Request, status int, message string, details map[string]string) {
	respondJSON(w, status, errorResponse{
		Error:     message,
		Details:   details,
		RequestID: server.RequestIDFromContext(r.Context()),
	})
}
