package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"travel-api/internal/contextkeys"
	"travel-api/services/order/application"
	"travel-api/services/order/domain"
)

type OrderHandler struct {
	service *application.OrderService
}

func NewOrderHandler(service *application.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func userIDFromRequest(r *http.Request) string {
	value := r.Context().Value(contextkeys.UserID)
	if value == nil {
		return ""
	}
	userID, _ := value.(string)
	return userID
}

func (h *OrderHandler) HandleOrders(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "login_required"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		orders, err := h.service.List(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "server_error"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	case http.MethodPost:
		var req domain.CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "bad_request"})
			return
		}
		order, err := h.service.Create(userID, req)
		if err != nil {
			status := http.StatusBadRequest
			errorCode := err.Error()
			if err == application.ErrProductNotFound || err == application.ErrPackageNotFound {
				status = http.StatusNotFound
			}
			if err != application.ErrInvalidOrderRequest && err != application.ErrProductNotFound && err != application.ErrPackageNotFound && err != application.ErrAvailabilityClosed {
				status = http.StatusInternalServerError
				errorCode = "server_error"
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": errorCode})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *OrderHandler) HandleOrderActions(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "login_required"})
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/orders/"), "/")
	if len(parts) < 2 || parts[1] != "cancel" {
		http.NotFound(w, r)
		return
	}
	orderID, _ := strconv.Atoi(parts[0])
	if orderID <= 0 {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, ok, err := h.service.Cancel(userID, orderID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "server_error"})
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "order_not_found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "order": order})
}
