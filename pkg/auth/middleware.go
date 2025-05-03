package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// Middleware is a JWT authentication middleware
func Middleware(jwtService JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			// Check if the Authorization header has the correct format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
				return
			}

			// Validate the token
			userID, err := jwtService.ValidateToken(parts[1])
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Add user ID to the request context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts the user ID from the request context
func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	return userID, ok
}