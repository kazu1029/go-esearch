package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kazu1029/gin-elastic/search-api/handlers"
	"github.com/olivere/elastic"
	"github.com/teris-io/shortid"
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
	// r.POST("/documents", func(c *gin.Context) { handlers.CreateDocumentsEndpoint(c) })
	r.POST("/documents", handlers.CreateDocumentsEndpoint)
	r.POST("/document", createDocumentsEndpoint)
	r.GET("/search", func(c *gin.Context) { handlers.SearchEndpoint(c) })
	if err = r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func errorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, gin.H{
		"error": err,
	})
}

func createDocumentsEndpoint(c *gin.Context) {
	var docs []DocumentRequest
	if err := c.BindJSON(&docs); err != nil {
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
		return
	}

	bulk := elasticClient.
		Bulk().
		Index(elasticIndexName).
		Type(elasticTypeName)

	for _, d := range docs {
		doc := Document{
			ID:        shortid.MustGenerate(),
			Title:     d.Title,
			CreatedAt: time.Now().UTC(),
			Content:   d.Content,
		}
		bulk.Add(elastic.NewBulkIndexRequest().Id(doc.ID).Doc(doc))
	}
	if _, err := bulk.Do(c.Request.Context()); err != nil {
		log.Printf("err is %+v\n", err)
		errorResponse(c, http.StatusInternalServerError, "Failed to create documents")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document Created",
	})
}
