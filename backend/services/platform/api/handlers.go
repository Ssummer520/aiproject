package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"travel-api/internal/contextkeys"
	"travel-api/services/platform/application"
	"travel-api/services/platform/domain"
)

type PlatformHandler struct{ service *application.PlatformService }

func NewPlatformHandler(service *application.PlatformService) *PlatformHandler {
	return &PlatformHandler{service: service}
}

func userIDFromRequest(r *http.Request) string {
	value, _ := r.Context().Value(contextkeys.UserID).(string)
	return value
}
func userEmailFromRequest(r *http.Request) string { return r.Header.Get("X-Demo-User-Email") }

func (h *PlatformHandler) HandlePlatform(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "login_required"})
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	snapshot, err := h.service.Snapshot(userID, userEmailFromRequest(r))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, snapshot)
}

func (h *PlatformHandler) HandlePlatformActions(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromRequest(r)
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "login_required"})
		return
	}
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/v1/platform/"), "/")
	switch path {
	case "merchants":
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		items, err := h.service.ListMerchants()
		respond(w, items, err)
	case "inventory":
		switch r.Method {
		case http.MethodGet:
			items, err := h.service.ListInventory()
			respond(w, items, err)
		case http.MethodPost:
			var req domain.InventoryUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
				return
			}
			item, err := h.service.UpdateInventory(req)
			respondCreated(w, item, err)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "orders":
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		items, err := h.service.ListOrders()
		respond(w, items, err)
	case "refunds":
		switch r.Method {
		case http.MethodGet:
			items, err := h.service.Snapshot(userID, "")
			if err != nil {
				respond(w, nil, err)
				return
			}
			writeJSON(w, http.StatusOK, items.Refunds)
		case http.MethodPost:
			var req domain.CreateRefundRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
				return
			}
			item, err := h.service.CreateRefund(req)
			respondCreated(w, item, err)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "profile":
		switch r.Method {
		case http.MethodGet:
			profile, err := h.service.Profile(userID, userEmailFromRequest(r))
			respond(w, profile, err)
		case http.MethodPost:
			var req domain.UserProfile
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
				return
			}
			profile, err := h.service.UpdateProfile(userID, req)
			respondCreated(w, profile, err)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "cms":
		switch r.Method {
		case http.MethodGet:
			items, err := h.service.ListCMS()
			respond(w, items, err)
		case http.MethodPost:
			var req domain.CMSArticle
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
				return
			}
			item, err := h.service.CreateCMS(req)
			respondCreated(w, item, err)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "metrics":
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		metrics, err := h.service.Metrics()
		respond(w, metrics, err)
	default:
		http.NotFound(w, r)
	}
}

func respond(w http.ResponseWriter, value interface{}, err error) {
	if err != nil {
		writeJSON(w, statusForErr(err), map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, value)
}
func respondCreated(w http.ResponseWriter, value interface{}, err error) {
	if err != nil {
		writeJSON(w, statusForErr(err), map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, value)
}
func statusForErr(err error) int {
	if err == application.ErrInvalidPlatformRequest {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
