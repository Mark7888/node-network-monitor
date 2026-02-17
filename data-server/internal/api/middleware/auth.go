package middleware

import (
	"net/http"
	"strings"

	"mark7888/speedtest-data-server/internal/auth"
	"mark7888/speedtest-data-server/internal/db"
	"mark7888/speedtest-data-server/internal/logger"
	"mark7888/speedtest-data-server/pkg/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// APIKeyAuth middleware validates API keys for node endpoints
func APIKeyAuth(database db.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Verify API key
		apiKey, err := database.VerifyAPIKey(token)
		if err != nil {
			logger.Log.Warn("Invalid API key attempt", zap.Error(err))
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid API key",
			})
			c.Abort()
			return
		}

		// Store API key in context for later use
		c.Set("api_key", apiKey)
		c.Next()
	}
}

// JWTAuth middleware validates JWT tokens for admin endpoints
func JWTAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate JWT
		claims, err := jwtManager.Validate(token)
		if err != nil {
			logger.Log.Warn("Invalid JWT attempt", zap.Error(err))
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("username", claims.Username)
		c.Next()
	}
}
