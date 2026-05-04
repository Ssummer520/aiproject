package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"travel-api/services/product/application"
	"travel-api/services/product/domain"
)

type ProductHandler struct {
	service *application.ProductService
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{service: application.NewProductService()}
}

func NewProductHandlerWithService(service *application.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	if destinationID, _ := strconv.Atoi(query.Get("destination_id")); destinationID > 0 {
		product, err := h.service.GetByDestinationID(destinationID)
		if err == application.ErrProductNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "product_not_found"})
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "server_error"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
		return
	}

	minPrice, _ := strconv.ParseFloat(firstNonEmpty(query.Get("price_min"), query.Get("min_price")), 64)
	maxPrice, _ := strconv.ParseFloat(firstNonEmpty(query.Get("price_max"), query.Get("max_price")), 64)
	ratingMin, _ := strconv.ParseFloat(query.Get("rating_min"), 64)

	result, err := h.service.Search(domain.SearchFilters{
		Query:          query.Get("q"),
		City:           query.Get("city"),
		Category:       query.Get("category"),
		Type:           query.Get("type"),
		MinPrice:       minPrice,
		MaxPrice:       maxPrice,
		RatingMin:      ratingMin,
		InstantConfirm: application.ParseBoolPointer(query.Get("instant_confirm")),
		FreeCancel:     application.ParseBoolPointer(query.Get("free_cancel")),
		Sort:           query.Get("sort"),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "server_error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *ProductHandler) HandleProductDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/products/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	id, _ := strconv.Atoi(parts[0])
	if id <= 0 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if len(parts) > 1 && parts[1] == "availability" {
		availability, err := h.service.Availability(id, r.URL.Query().Get("date"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "server_error"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"availability": availability})
		return
	}

	product, err := h.service.Get(id)
	if err == application.ErrProductNotFound {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "product_not_found"})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "server_error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
