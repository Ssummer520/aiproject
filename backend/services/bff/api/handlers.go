package api

import (
	"encoding/json"
	"net/http"
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
