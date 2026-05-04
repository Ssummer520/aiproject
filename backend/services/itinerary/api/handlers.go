package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"travel-api/internal/contextkeys"
	"travel-api/services/itinerary/application"
	"travel-api/services/itinerary/domain"
)

type ItineraryHandler struct {
	service *application.ItineraryService
}

func NewItineraryHandler(service *application.ItineraryService) *ItineraryHandler {
	return &ItineraryHandler{service: service}
}

func userIDFromRequest(r *http.Request) string {
	value, _ := r.Context().Value(contextkeys.UserID).(string)
	return value
}

func (h *ItineraryHandler) HandleItineraries(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "login_required"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		items, err := h.service.List(userID)
		if err != nil {
			writeJSON(w, statusForErr(err), map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, items)
	case http.MethodPost:
		var req domain.CreateItineraryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
			return
		}
		item, err := h.service.Create(userID, req)
		if err != nil {
			writeJSON(w, statusForErr(err), map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, item)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ItineraryHandler) HandleItineraryActions(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "login_required"})
		return
	}
	parts := strings.Split(strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/v1/itineraries/"), "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	itineraryID, _ := strconv.Atoi(parts[0])
	if itineraryID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
		return
	}
	if len(parts) == 1 && r.Method == http.MethodGet {
		item, err := h.service.Get(userID, itineraryID)
		if err != nil {
			writeJSON(w, statusForErr(err), map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}
	if len(parts) >= 2 && parts[1] == "items" && r.Method == http.MethodPost {
		var req domain.AddItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
			return
		}
		item, err := h.service.AddItem(userID, itineraryID, req)
		if err != nil {
			writeJSON(w, statusForErr(err), map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, item)
		return
	}
	if len(parts) >= 4 && parts[1] == "items" && parts[3] == "move" && r.Method == http.MethodPost {
		itemID, _ := strconv.Atoi(parts[2])
		var req struct {
			Direction string `json:"direction"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		item, err := h.service.MoveItem(userID, itineraryID, itemID, req.Direction)
		if err != nil {
			writeJSON(w, statusForErr(err), map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}
	http.NotFound(w, r)
}

func (h *ItineraryHandler) HandleGenerate(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "login_required"})
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req domain.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
		return
	}
	item, err := h.service.Generate(userID, req)
	if err != nil {
		writeJSON(w, statusForErr(err), map[string]string{"error": err.Error()})
		return
	}
	status := http.StatusOK
	if req.Save {
		status = http.StatusCreated
	}
	writeJSON(w, status, item)
}

func statusForErr(err error) int {
	if err == application.ErrItineraryNotFound {
		return http.StatusNotFound
	}
	if err == application.ErrInvalidItineraryRequest {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
