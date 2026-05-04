package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"travel-api/internal/contextkeys"
	"travel-api/services/review/application"
	"travel-api/services/review/domain"
)

type ReviewHandler struct {
	service *application.ReviewService
}

func NewReviewHandler(service *application.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

func (h *ReviewHandler) HandleProductReviews(w http.ResponseWriter, r *http.Request) {
	productID, ok := productIDFromPath(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		result, err := h.service.List(productID, r.URL.Query().Get("language"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "server_error"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	case http.MethodPost:
		userID, _ := r.Context().Value(contextkeys.UserID).(string)
		if userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "login_required"})
			return
		}
		var req domain.CreateReviewRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "bad_request"})
			return
		}
		review, err := h.service.Create(userID, productID, req)
		if err != nil {
			status := http.StatusBadRequest
			if err == application.ErrReviewNotAllowed {
				status = http.StatusForbidden
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(review)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func productIDFromPath(path string) (int, bool) {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/api/v1/products/"), "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) != 2 || parts[1] != "reviews" {
		return 0, false
	}
	id, _ := strconv.Atoi(parts[0])
	return id, id > 0
}
