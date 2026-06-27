package middleware

import (
	"net/http"
	"strings"

	"ticket-system/utils"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.AbortWithError(c, http.StatusUnauthorized, "Authorization header required")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.AbortWithError(c, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		claims, err := utils.ParseToken(parts[1], secret)
		if err != nil {
			utils.AbortWithError(c, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}
