package api

import (
	"encoding/json"
	"net/http"

	"travel-api/services/coupon/application"
	"travel-api/services/coupon/domain"
)

type CouponHandler struct {
	service *application.CouponService
}

func NewCouponHandler(service *application.CouponService) *CouponHandler {
	return &CouponHandler{service: service}
}

func (h *CouponHandler) HandleCoupons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	coupons, err := h.service.ListActive()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "server_error"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"coupons": coupons})
}

func (h *CouponHandler) HandleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req domain.ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "bad_request"})
		return
	}
	result, err := h.service.Validate(req.Code, req.Amount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": result.Error, "result": result})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
