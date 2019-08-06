package main

import (
	"log"
	"net/http"
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
	r.GET("/index", func(c *gin.Context) { handlers.CreateMapping(c) })
	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 & time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
