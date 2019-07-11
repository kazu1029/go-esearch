package main

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kazu1029/gin-elastic/search-api/handlers"
)

const (
	elasticIndexName = "documents"
	elasticTypeName  = "document"
)

type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}

type DocumentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type DocumentResponse struct {
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}

type SearchResponse struct {
	Time      string             `json:"time"`
	Hits      string             `json:"hits"`
	Documents []DocumentResponse `json:"documents"`
}

var elasticClient *elastic.Client

func main() {
	var err error
	for {
		// elasticClient, err = elastic.NewClient(
		// 	elastic.SetURL("http://elasticsearch:9200"),
		// 	elastic.SetSniff(false),
		// )
		elasticClient, err = handlers.InitElastic()
		fmt.Printf("elasticClient is %+v\n", elasticClient)
		fmt.Printf("elasticClient type is %+v\n", reflect.TypeOf(elasticClient))
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
	if err = r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
