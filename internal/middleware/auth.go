package middleware

import (
	"book-lending-api/internal/domain"
	"book-lending-api/pkg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the Authorization header for bearer tokens.
// If the token is valid the user id and email are injected into the
// context.  Otherwise the request is aborted with 401.
func AuthMiddleware(jwtUtil *pkg.JWTUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Authorization header required",
			})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid authorization header format",
			})
			c.Abort()
			return
		}
		claims, err := jwtUtil.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Next()
	}
}

// GetUserIDFromContext extracts the user id from the context.
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	if id, ok := userID.(uint); ok {
		return id, true
	}
	return 0, false
}
