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

func (h *BFFHandler) HandleDestinations(w http.ResponseWriter, r *http.Request) {
	// Pattern: /api/v1/destinations/{id}/action
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/destinations/"), "/")
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
