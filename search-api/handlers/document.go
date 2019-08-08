package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kazu1029/gin-elastic/search-api/pkg/esearch"
	"github.com/olivere/elastic"
)

const (
	elasticIndexName = "documents"
	elasticTypeName  = "document"
)

type DocumentResponse struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SearchResponse struct {
	Time      string             `json:"time"`
	Hits      string             `json:"hits"`
	Documents []DocumentResponse `json:"documents"`
}

func NewElasticIndex(client *elastic.Client) *esearch.IndexService {
	return esearch.NewIndexService(client)
}

func NewElasticSearch(client *elastic.Client) *esearch.SearchService {
	return esearch.NewSearchService(client)
}

func CreateDocumentsEndpoint(c *gin.Context) {
	elasticClient, err := InitElastic()
	if err != nil {
		log.Println(err)
		errorResponse(c, http.StatusBadRequest, err.Error())
	}
	ctx := context.Background()
	// TODO: need to dynamic variables
	var docs []esearch.DocumentRequest
	if err := c.BindJSON(&docs); err != nil {
		log.Printf("err is %+v\n", err)
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
		return
	}

	index := NewElasticIndex(elasticClient)
	res, err := index.BulkInsert(ctx, docs, elasticIndexName, elasticTypeName)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": res,
	})
}

func CreateMapping(c *gin.Context) {
	elasticClient, err := InitElastic()
	ctx := context.Background()
	if err != nil {
		log.Println(err)
	}

	indexName := c.Param("index_name")
	var mapping interface{}
	if err := c.BindJSON(&mapping); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	index := NewElasticIndex(elasticClient)
	res, err := index.CreateMapping(ctx, indexName, mapping)
	if err != nil {
		log.Printf("err is %+v\n", err)
		errorResponse(c, http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": res,
	})
}

func SearchEndpoint(c *gin.Context) {
	elasticClient, err := InitElastic()
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}
	ctx := context.Background()
	queries := c.Request.URL.Query()
	query := queries["query"][0]
	if query == "" {
		errorResponse(c, http.StatusBadRequest, "Query not specified")
		return
	}

	skip := 0
	take := 10
	if i, err := strconv.Atoi(c.Query("skip")); err != nil {
		skip = i
	}
	if i, err := strconv.Atoi(c.Query("take")); err != nil {
		take = i
	}
	types := []string{"title", "content"}

	search := NewElasticSearch(elasticClient)
	res, err := search.SearchMultiMatchQuery(ctx, elasticIndexName, skip, take, query, types...)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, res)
}

func CreateIndexTemplate(c *gin.Context) {
	elasticClient, err := InitElastic()
	ctx := context.Background()
	if err != nil {
		log.Fatal(err)
	}

	templateName := c.Param("template_name")
	var template interface{}
	if err := c.BindJSON(&template); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	index := NewElasticIndex(elasticClient)
	res, err := index.CreateIndexTemplate(ctx, templateName, template)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": res,
	})
}
