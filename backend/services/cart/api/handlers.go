package api

import (
	"encoding/json"
	"net/http"

	"travel-api/internal/contextkeys"
	"travel-api/services/cart/application"
	"travel-api/services/cart/domain"
)

type CartHandler struct {
	service *application.CartService
}

func NewCartHandler(service *application.CartService) *CartHandler {
	return &CartHandler{service: service}
}

func userIDFromRequest(r *http.Request) string {
	value, _ := r.Context().Value(contextkeys.UserID).(string)
	return value
}

func (h *CartHandler) HandleCart(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "login_required"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		summary, err := h.service.List(userID)
		if err != nil {
			writeJSON(w, statusForErr(err), map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, summary)
	case http.MethodPost:
		var req domain.AddCartItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
			return
		}
		summary, err := h.service.Add(userID, req)
		if err != nil {
			writeJSON(w, statusForErr(err), map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, summary)
	case http.MethodDelete:
		if err := h.service.Clear(userID); err != nil {
			writeJSON(w, statusForErr(err), map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CartHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "login_required"})
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req domain.CheckoutRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	orders, err := h.service.Checkout(userID, req)
	if err != nil {
		writeJSON(w, statusForErr(err), map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"ok": true, "orders": orders})
}

func statusForErr(err error) int {
	if err == application.ErrInvalidCartRequest || err == application.ErrCartEmpty {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
