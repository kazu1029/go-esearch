package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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

func ElasticIndex(client *elastic.Client) *esearch.IndexService {
	return esearch.NewIndexService(client)
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

	index := ElasticIndex(elasticClient)
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

	index := ElasticIndex(elasticClient)
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

	esQuery := elastic.NewMultiMatchQuery(query, "title", "content").
		Fuzziness("2").
		MinimumShouldMatch("2")
	result, err := elasticClient.Search().
		Index(elasticIndexName).
		Query(esQuery).
		From(skip).Size(take).
		Do(c.Request.Context())
	if err != nil {
		log.Printf("err is %v\n", err)
		errorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	res := SearchResponse{
		Time: fmt.Sprintf("%d", result.TookInMillis),
		Hits: fmt.Sprintf("%d", result.Hits.TotalHits),
	}
	docs := make([]DocumentResponse, 0)
	for _, hit := range result.Hits.Hits {
		var doc DocumentResponse
		json.Unmarshal(*hit.Source, &doc)
		docs = append(docs, doc)
	}
	res.Documents = docs

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

	index := ElasticIndex(elasticClient)
	res, err := index.CreateIndexTemplate(ctx, templateName, template)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": res,
	})
}
