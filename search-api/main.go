package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kazu1029/gin-elastic/search-api/handlers"
	"github.com/olivere/elastic/v6"
	"github.com/teris-io/shortid"
)

var (
	elasticClient *elastic.Client
)

func main() {
	var err error
	for {
		elasticClient, err = elastic.NewClient(
			elastic.SetURL("http://elasticsearch:9200"),
			elastic.SetSniff(false),
		)
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	r := gin.Default()
	r.POST("/documents", handlers.CreateDocumentsEndpoint)
	r.GET("/search", handlers.SearchEndpoint)
	// r.POST("/documents", createDocumentsEndpoint)
	// r.GET("/search", searchEndpoint)
	if err = r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
