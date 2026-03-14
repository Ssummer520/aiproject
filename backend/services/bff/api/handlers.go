package api

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	// In a real microservice, we might call other services passing 'lang'
	// and they would return translated content.
	// For now, let's keep it simple as a proof of concept.

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *BFFHandler) ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	if id == 0 {
		// Try parsing from body if not in query
		type req struct {
			ID int `json:"id"`
		}
		var rr req
		if err := json.NewDecoder(r.Body).Decode(&rr); err == nil {
			id = rr.ID
		}
	}

	if id == 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	isFav := h.service.ToggleFavorite(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "is_favorite": isFav})
}

func (h *BFFHandler) RecordView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	type req struct {
		ID int `json:"id"`
	}
	var rr req
	if err := json.NewDecoder(r.Body).Decode(&rr); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.service.AddToHistory(rr.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
}
