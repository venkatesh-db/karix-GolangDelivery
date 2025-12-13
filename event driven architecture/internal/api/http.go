package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/venkatesh/order-service/internal/app"
	"github.com/venkatesh/order-service/internal/domain"
	"github.com/venkatesh/order-service/internal/readmodel"
)

// Server wires HTTP handlers to the application service.
type Server struct {
	svc   *app.OrderService
	views *readmodel.OrdersProjection
}

// NewServer builds an HTTP server wrapper.
func NewServer(svc *app.OrderService, views *readmodel.OrdersProjection) *Server {
	return &Server{svc: svc, views: views}
}

// Router returns mux with all endpoints registered.
func (s *Server) Router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	r.HandleFunc("/orders", s.placeOrder).Methods(http.MethodPost)
	r.HandleFunc("/orders", s.listOrders).Methods(http.MethodGet)
	r.HandleFunc("/orders/{id}", s.getOrder).Methods(http.MethodGet)
	r.HandleFunc("/orders/{id}/payment", s.authorizePayment).Methods(http.MethodPost)
	r.HandleFunc("/orders/{id}/reserve", s.reserveInventory).Methods(http.MethodPost)
	r.HandleFunc("/orders/{id}/ship", s.shipOrder).Methods(http.MethodPost)
	r.HandleFunc("/orders/{id}/cancel", s.cancelOrder).Methods(http.MethodPost)
	return r
}

type placeOrderRequest struct {
	OrderID    string            `json:"order_id"`
	CustomerID string            `json:"customer_id"`
	Items      []domain.LineItem `json:"items"`
}

type paymentRequest struct {
	PaymentID string `json:"payment_id"`
	Amount    int64  `json:"amount_cents"`
}

type reserveRequest struct {
	ReservationID string `json:"reservation_id"`
}

type shipRequest struct {
	TrackingNumber string `json:"tracking_number"`
	Carrier        string `json:"carrier"`
}

type cancelRequest struct {
	Reason string `json:"reason"`
}

func (s *Server) placeOrder(w http.ResponseWriter, r *http.Request) {
	var req placeOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}
	cmd := app.PlaceOrder{OrderID: req.OrderID, CustomerID: req.CustomerID, Items: req.Items}
	if err := s.svc.HandlePlaceOrder(r.Context(), cmd); err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}
	respondOK(w, map[string]string{"order_id": req.OrderID})
}

func (s *Server) authorizePayment(w http.ResponseWriter, r *http.Request) {
	var req paymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}
	cmd := app.AuthorizePayment{
		OrderID:   mux.Vars(r)["id"],
		PaymentID: req.PaymentID,
		Amount:    req.Amount,
	}
	if err := s.svc.HandleAuthorizePayment(r.Context(), cmd); err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}
	respondAccepted(w)
}

func (s *Server) reserveInventory(w http.ResponseWriter, r *http.Request) {
	var req reserveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}
	cmd := app.ReserveInventory{OrderID: mux.Vars(r)["id"], ReservationID: req.ReservationID}
	if err := s.svc.HandleReserveInventory(r.Context(), cmd); err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}
	respondAccepted(w)
}

func (s *Server) shipOrder(w http.ResponseWriter, r *http.Request) {
	var req shipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}
	cmd := app.ShipOrder{OrderID: mux.Vars(r)["id"], TrackingNumber: req.TrackingNumber, Carrier: req.Carrier}
	if err := s.svc.HandleShipOrder(r.Context(), cmd); err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}
	respondAccepted(w)
}

func (s *Server) cancelOrder(w http.ResponseWriter, r *http.Request) {
	var req cancelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}
	cmd := app.CancelOrder{OrderID: mux.Vars(r)["id"], Reason: req.Reason}
	if err := s.svc.HandleCancelOrder(r.Context(), cmd); err != nil {
		respondErr(w, http.StatusBadRequest, err)
		return
	}
	respondAccepted(w)
}

func (s *Server) listOrders(w http.ResponseWriter, r *http.Request) {
	respondOK(w, s.views.List())
}

func (s *Server) getOrder(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]
	view, ok := s.views.Get(orderID)
	if !ok {
		http.NotFound(w, r)
		return
	}
	respondOK(w, view)
}

func respondAccepted(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte(`{"status":"accepted"}`))
}

func respondOK(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondErr(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
