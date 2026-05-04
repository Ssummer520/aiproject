package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"travel-api/internal/contextkeys"
	"travel-api/services/user/application"
	"travel-api/services/user/domain"
)

type UserHandler struct {
	service *application.UserService
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) HandleMe(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/v1/users/me"), "/")
	switch {
	case path == "profile":
		h.handleProfile(w, r)
	case path == "travelers":
		h.handleTravelers(w, r)
	case strings.HasPrefix(path, "travelers/"):
		h.handleTravelerAction(w, r, strings.TrimPrefix(path, "travelers/"))
	case path == "membership":
		h.handleMembership(w, r)
	case path == "roles":
		h.handleRoles(w, r)
	default:
		writeError(w, http.StatusNotFound, "not_found")
	}
}

func (h *UserHandler) handleProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	switch r.Method {
	case http.MethodGet:
		profile, err := h.service.GetProfile(userID, "")
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, profile)
	case http.MethodPatch:
		var req domain.UserProfile
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "bad_request")
			return
		}
		profile, err := h.service.UpsertProfile(userID, req)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, profile)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

func (h *UserHandler) handleTravelers(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	switch r.Method {
	case http.MethodGet:
		travelers, err := h.service.ListTravelers(userID)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"travelers": travelers})
	case http.MethodPost:
		var req domain.TravelerProfileInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "bad_request")
			return
		}
		traveler, err := h.service.CreateTraveler(userID, req)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, traveler)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

func (h *UserHandler) handleTravelerAction(w http.ResponseWriter, r *http.Request, suffix string) {
	userID, ok := currentUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	parts := strings.Split(strings.Trim(suffix, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		writeError(w, http.StatusNotFound, "not_found")
		return
	}
	travelerID, err := strconv.Atoi(parts[0])
	if err != nil || travelerID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid_traveler_id")
		return
	}
	if len(parts) == 2 && parts[1] == "default" {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
			return
		}
		traveler, err := h.service.SetDefaultTraveler(userID, travelerID)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, traveler)
		return
	}
	if len(parts) != 1 {
		writeError(w, http.StatusNotFound, "not_found")
		return
	}
	switch r.Method {
	case http.MethodGet:
		traveler, err := h.service.GetTraveler(userID, travelerID)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, traveler)
	case http.MethodPatch:
		var req domain.TravelerProfileInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "bad_request")
			return
		}
		traveler, err := h.service.UpdateTraveler(userID, travelerID, req)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, traveler)
	case http.MethodDelete:
		if err := h.service.DeleteTraveler(userID, travelerID); err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

func (h *UserHandler) handleMembership(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	userID, ok := currentUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	membership, err := h.service.GetMembership(userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, membership)
}

func (h *UserHandler) handleRoles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	userID, ok := currentUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	roles, err := h.service.GetRoles(userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"roles": roles})
}

func currentUserID(r *http.Request) (string, bool) {
	userID, _ := r.Context().Value(contextkeys.UserID).(string)
	userID = strings.TrimSpace(userID)
	return userID, userID != ""
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, application.ErrInvalidUserRequest):
		writeError(w, http.StatusBadRequest, "invalid_user_request")
	case errors.Is(err, application.ErrTravelerNotFound):
		writeError(w, http.StatusNotFound, "traveler_not_found")
	case errors.Is(err, application.ErrDocumentDuplicate):
		writeError(w, http.StatusConflict, "document_duplicate")
	default:
		writeError(w, http.StatusInternalServerError, "server_error")
	}
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, map[string]string{"error": code})
}
