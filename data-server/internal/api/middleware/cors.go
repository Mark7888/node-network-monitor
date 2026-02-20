package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a configured CORS middleware.
// allowedOrigins should contain the list of permitted origins; pass ["*"] to allow all.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	// AllowCredentials cannot be true when origin is wildcard
	allowCredentials := !(len(allowedOrigins) == 1 && allowedOrigins[0] == "*")

	config := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: allowCredentials,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(config)
}
