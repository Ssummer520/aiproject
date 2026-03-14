package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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

func (h *BFFHandler) GetHomePage(w http.ResponseWriter, r *http.Request) {
	// Support multi-language from header
	lang := r.Header.Get("Accept-Language")
	if lang == "" {
		lang = "en"
	}

	data := h.service.GetHomePageData(lang)

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

	switch action {
	case "favorite":
		isFav := h.service.ToggleFavorite(id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "is_favorite": isFav})
	case "view":
		h.service.AddToHistory(id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
	default:
		http.Error(w, "Action not found", http.StatusNotFound)
	}
}
