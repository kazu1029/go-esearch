package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kazu1029/go-elastic/search-api/handlers"
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
	r.POST("/bulk/:index_name/:type_name", func(c *gin.Context) { handlers.CreateDocumentsEndpoint(c) })
	r.POST("/index/:index_name/template/:template_name", func(c *gin.Context) { handlers.CreateIndexTemplate(c) })
	r.POST("/index/:index_name/mapping", func(c *gin.Context) { handlers.CreateMapping(c) })
	r.GET("/search", func(c *gin.Context) { handlers.SearchEndpoint(c) })

	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 & time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
