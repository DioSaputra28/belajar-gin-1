package middleware

import (
	"net/http"
	"strings"

	"github.com/DioSaputra28/belajar-gin-1/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authRepo auth.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		token = strings.TrimSpace(token)

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}

		user, err := authRepo.FindUserByToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Next()
	}
}