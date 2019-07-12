package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kazu1029/gin-elastic/search-api/handlers"
)

func main() {
	for {
		_, err := handlers.InitElastic()
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	r := gin.Default()
	r.POST("/documents", func(c *gin.Context) { handlers.CreateDocumentsEndpoint(c) })
	r.GET("/search", func(c *gin.Context) { handlers.SearchEndpoint(c) })
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
