package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"travel-api/services/inbound/application"
	"travel-api/services/inbound/domain"
)

type InboundHandler struct{ service *application.InboundService }

func NewInboundHandler(service *application.InboundService) *InboundHandler {
	return &InboundHandler{service: service}
}
func (h *InboundHandler) HandleInbound(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	data, err := h.service.Snapshot()
	respond(w, data, err, http.StatusOK)
}
func (h *InboundHandler) HandleActions(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/v1/inbound/"), "/")
	if strings.HasPrefix(path, "cities/") && strings.HasSuffix(path, "/guide") && r.Method == http.MethodGet {
		city := strings.TrimSuffix(strings.TrimPrefix(path, "cities/"), "/guide")
		guide, ok, err := h.service.Guide(city)
		if err != nil {
			respond(w, nil, err, statusForErr(err))
			return
		}
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "guide_not_found"})
			return
		}
		writeJSON(w, http.StatusOK, guide)
		return
	}
	if path == "concierge" && r.Method == http.MethodPost {
		var req domain.ConciergeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
			return
		}
		data, err := h.service.Concierge(req)
		respond(w, data, err, http.StatusOK)
		return
	}
	http.NotFound(w, r)
}
func respond(w http.ResponseWriter, value interface{}, err error, status int) {
	if err != nil {
		writeJSON(w, statusForErr(err), map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, status, value)
}
func statusForErr(err error) int {
	if err == application.ErrInvalidInboundRequest {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
