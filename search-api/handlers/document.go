package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kazu1029/go-elastic/esearch"
	"github.com/olivere/elastic"
)

func NewElasticIndex(client *elastic.Client) *esearch.IndexService {
	return esearch.NewIndexService(client)
}

func NewElasticSearch(client *elastic.Client) *esearch.SearchService {
	return esearch.NewSearchService(client)
}

func CreateDocumentsEndpoint(c *gin.Context) {
	elasticClient, err := InitElastic()
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}
	indexName := c.Param("index_name")
	typeName := c.Param("type_name")
	ctx := context.Background()
	// TODO: need to dynamic variables
	var docs []interface{}
	if err := c.BindJSON(&docs); err != nil {
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
	}

	index := NewElasticIndex(elasticClient)
	res, err := index.BulkInsert(ctx, docs, indexName, typeName)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"message": res,
	})
}

func CreateMapping(c *gin.Context) {
	elasticClient, err := InitElastic()
	ctx := context.Background()
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}

	indexName := c.Param("index_name")
	var mapping interface{}
	if err := c.BindJSON(&mapping); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}

	index := NewElasticIndex(elasticClient)
	res, err := index.CreateMapping(ctx, indexName, mapping)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": res,
	})
}

func SearchEndpoint(c *gin.Context) {
	var query string
	var targetTypes []string
	elasticClient, err := InitElastic()
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}
	ctx := context.Background()
	queries := c.Request.URL.Query()
	query = queries["query"][0]
	// TODO: accept the other symbols
	targetTypes = strings.Split(queries["target_types"][0], ",")
	indexName := c.Param("index_name")
	if query == "" {
		errorResponse(c, http.StatusBadRequest, "Query not specified")
	}

	skip := 0
	take := 50
	if i, err := strconv.Atoi(c.Query("skip")); err == nil {
		skip = i
	}
	if i, err := strconv.Atoi(c.Query("take")); err == nil {
		take = i
	}

	search := NewElasticSearch(elasticClient)
	res, err := search.SearchMultiMatchQuery(ctx, indexName, skip, take, query, targetTypes...)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, res)
}

func CreateIndexTemplate(c *gin.Context) {
	elasticClient, err := InitElastic()
	ctx := context.Background()
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}

	templateName := c.Param("template_name")
	var template interface{}
	if err := c.BindJSON(&template); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
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
