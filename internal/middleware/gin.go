package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/afreedicp/zolaris-backend-app/internal/services"
)

// GinAuthMiddleware checks for user authentication
// This is a simplified version - in a real app, use JWT or OAuth2
func GinAuthMiddleware(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		cogntioId := c.GetHeader("X-Cognito-ID")

		userID, err := userService.GetUserIdByCognitoId(c.Request.Context(), cogntioId)
		if err != nil {
			log.Printf("Error retrieving user ID by Cognito ID: %v", err)
			c.JSON(500, gin.H{"status": false, "message": "Internal server error"})
			c.Abort()
			return
		}

		// For demo purposes only - in production, never use a default user
		if userID == "" {
			c.JSON(401, gin.H{"status": false, "message": "Unauthorized: Invalid user ID. Could not retrieve user by Cognito ID."})
			c.Abort()
			return
		}

		// Log authentication
		log.Printf("Authenticated request for user: %s", userID)

		// Add user ID to request context
		c.Set(string(UserIDKey), userID)

		c.Next()
	}
}

// GinLoggerMiddleware logs request details
func GinLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Log request details
		log.Printf("%s | %s | %d | %s | %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
			c.ClientIP(),
		)
	}
}

// GetUserIDFromGin extracts the user ID from the Gin context
func GetUserIDFromGin(c *gin.Context) string {
	userID, exists := c.Get(string(UserIDKey))
	if !exists {
		return ""
	}
	return userID.(string)
}
