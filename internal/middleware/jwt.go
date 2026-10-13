package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"simple-auction-system/internal/services"
)

type contextKey string

const userIDKey contextKey = "userID"

func RequireAuth(
	authService services.AuthService,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			// 1. Lấy Authorization header
			authHeader := r.Header.Get("Authorization")

			// 2. Header phải có dạng: Bearer <token>
			if !strings.HasPrefix(authHeader, "Bearer ") {
				unauthorized(w)
				return
			}

			// 3. Bỏ phần "Bearer " để lấy JWT
			token := strings.TrimPrefix(authHeader, "Bearer ")

			if token == "" {
				unauthorized(w)
				return
			}

			// 4. Kiểm tra JWT
			claims, err := authService.ParseToken(token)
			if err != nil {
				unauthorized(w)
				return
			}

			// 5. Lưu user ID vào request
			ctx := context.WithValue(
				r.Context(),
				userIDKey,
				claims.UserID,
			)

			// 6. JWT hợp lệ, cho request chạy tiếp
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) (uint, bool) {
	userID, ok := r.Context().Value(userIDKey).(uint)
	return userID, ok
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": "unauthorized",
	})
}
