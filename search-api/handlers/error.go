package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
)

func errorResponse(c *gin.Context, code int, err string) {
	log.Println(err)
	c.JSON(code, gin.H{
		"error": err,
	})
}
