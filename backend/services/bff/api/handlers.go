package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"travel-api/internal/contextkeys"
	"travel-api/services/bff/application"
)

type BFFHandler struct {
	service *application.BFFService
}

func NewBFFHandler() *BFFHandler {
	return &BFFHandler{
		service: application.NewBFFService(),
	}
}

func userIDFromRequest(r *http.Request) string {
	v := r.Context().Value(contextkeys.UserID)
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func (h *BFFHandler) GetHomePage(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	if lang == "" {
		lang = "en"
	}
	userID := userIDFromRequest(r)

	data := h.service.GetHomePageData(lang, userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *BFFHandler) Search(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	if lang == "" {
		lang = "en"
	}
	userID := userIDFromRequest(r)

	query := r.URL.Query().Get("q")
	city := r.URL.Query().Get("city")
	category := r.URL.Query().Get("category")
	minPrice, _ := strconv.Atoi(r.URL.Query().Get("min_price"))
	maxPrice, _ := strconv.Atoi(r.URL.Query().Get("max_price"))

	data := h.service.SearchDestinations(lang, userID, query, city, category, minPrice, maxPrice)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *BFFHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	if lang == "" {
		lang = "en"
	}
	userID := userIDFromRequest(r)

	category := strings.TrimPrefix(r.URL.Path, "/api/v1/category/")

	data := h.service.GetCategoryData(lang, userID, category)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *BFFHandler) GetCity(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	if lang == "" {
		lang = "en"
	}
	userID := userIDFromRequest(r)

	city := strings.TrimPrefix(r.URL.Path, "/api/v1/city/")

	data := h.service.GetCityData(lang, userID, city)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *BFFHandler) HandleDestinations(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/destinations/"), "/")

	// Single destination detail endpoint
	if len(parts) == 1 {
		id, _ := strconv.Atoi(parts[0])
		if id > 0 {
			lang := r.Header.Get("Accept-Language")
			if lang == "" {
				lang = "en"
			}
			userID := userIDFromRequest(r)

			data := h.service.GetDestinationDetail(lang, userID, id)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
			return
		}
	}

	if len(parts) < 2 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	id, _ := strconv.Atoi(parts[0])
	action := parts[1]

	if id == 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := userIDFromRequest(r)
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "login_required"})
		return
	}

	switch action {
	case "favorite":
		isFav := h.service.ToggleFavorite(userID, id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "is_favorite": isFav})
	case "view":
		h.service.AddToHistory(userID, id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
	default:
		http.Error(w, "Action not found", http.StatusNotFound)
	}
}

func (h *BFFHandler) HandleBookings(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "login_required"})
		return
	}

	if r.Method == "GET" {
		bookings := h.service.GetUserBookings(userID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bookings)
		return
	}

	if r.Method == "POST" {
		var req struct {
			DestinationID int    `json:"destination_id"`
			CheckIn       string `json:"check_in"`
			CheckOut      string `json:"check_out"`
			Guests        int    `json:"guests"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		booking := h.service.CreateBooking(userID, req.DestinationID, req.CheckIn, req.CheckOut, req.Guests)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(booking)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *BFFHandler) HandleNotifications(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "login_required"})
		return
	}

	if r.Method == "GET" {
		notifications := h.service.GetNotifications(userID)
		unreadCount := h.service.GetUnreadNotificationCount(userID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"notifications": notifications,
			"unread_count":  unreadCount,
		})
		return
	}

	if r.Method == "POST" {
		var req struct {
			NotificationID int `json:"notification_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		h.service.MarkNotificationAsRead(userID, req.NotificationID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "marked_as_read"})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
