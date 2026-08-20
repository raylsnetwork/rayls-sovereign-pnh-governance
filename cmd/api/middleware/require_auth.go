package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/api/core"
)

// RequireAuth returns a middleware that validates JWT tokens using AuthService
func RequireAuth(authService core.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		pn, err := authService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("privateNetwork", pn)
		c.Next()
	}
}
