package handlers

import (
	"github.com/gin-gonic/gin"
)

func errorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, gin.H{
		"error": err,
	})
}
