package middleware

import (
	"log"
	"net/http"

	"ticket-system/utils"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				utils.JSONError(c, http.StatusInternalServerError, "Internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
