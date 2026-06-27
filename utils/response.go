package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func AbortWithError(c *gin.Context, status int, message string) {
	JSONError(c, status, message)
	c.Abort()
}

func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		JSONError(c, http.StatusBadRequest, "Invalid request body")
		return false
	}
	return true
}
