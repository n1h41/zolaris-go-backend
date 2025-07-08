package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/afreedicp/zolaris-backend-app/internal/transport/dto"
)

// UserIDKey is the context key for storing the user ID
type contextKey string

const UserIDKey contextKey = "userID"

// sendJSONError sends a JSON error response
func sendJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := dto.ErrorResponse{
		Success: false,
		Error:   message,
		Code:    "UNAUTHORIZED",
	}
	json.NewEncoder(w).Encode(response)
}

// AuthMiddleware checks for user authentication
// This is a simplified version - in a real app, use JWT or OAuth2
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from header
		userID := r.Header.Get("X-User-ID")

		// For demo purposes only - in production, never use a default user
		if userID == "" {
			sendJSONError(w, http.StatusUnauthorized, "Unauthorized: Missing user ID")
			return
		}

		// Log authentication
		log.Printf("Authenticated request for user: %s", userID)

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts the user ID from the request context
func GetUserID(r *http.Request) string {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}
