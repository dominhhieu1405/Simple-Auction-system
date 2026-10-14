package handlers

import (
	"errors"
	"log"
	"net/http"
	"simple-auction-system/internal/middleware"
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

// Register đăng ký tài khoản mới.
//
// @Summary Đăng ký tài khoản
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Thông tin đăng ký"
// @Success 201 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
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

// Login đăng nhập tài khoản.
//
// @Summary Đăng nhập
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Thông tin đăng nhập"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
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

// Me lấy thông tin người dùng hiện tại.
//
// @Summary Thông tin người dùng
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/me [get]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Lấy user ID mà JWT middleware đã lưu
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Tìm user trong database
	user, err := h.service.GetUserByID(userID)

	switch {
	case err == nil:
		writeJSON(w, http.StatusOK, user)

	case errors.Is(err, services.ErrUserNotFound):
		writeError(w, http.StatusNotFound, err.Error())

	default:
		log.Printf("get current user: %v", err)
		writeError(
			w,
			http.StatusInternalServerError,
			"internal server error",
		)
	}
}
