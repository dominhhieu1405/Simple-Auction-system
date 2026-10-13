package handlers

import (
	"errors"
	"log"
	"net/http"
	"simple-auction-system/internal/services"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Cấu trúc payload nhận từ request
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	err := h.service.Register(req.Username, req.Password)
	switch {
	case err == nil:
		writeJSON(w, http.StatusCreated, map[string]string{"message": "user registered successfully"})
	case errors.Is(err, services.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrUserExists):
		writeError(w, http.StatusConflict, err.Error())
	default:
		log.Printf("register: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	token, err := h.service.Login(req.Username, req.Password)
	switch {
	case err == nil:
		writeJSON(w, http.StatusOK, map[string]string{"token": token})
	case errors.Is(err, services.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, err.Error())
	default:
		log.Printf("login: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
