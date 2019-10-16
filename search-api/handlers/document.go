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
	pipeline := c.Param("pipeline")
	ctx := context.Background()
	// TODO: need to dynamic variables
	var docs []interface{}
	if err := c.BindJSON(&docs); err != nil {
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
	}

	index := NewElasticIndex(elasticClient)
	res, err := index.BulkInsert(ctx, docs, indexName, typeName, pipeline)
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
	var query, sortField string
	var targetTypes []string
	var ascending bool
	elasticClient, err := InitElastic()
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}
	ctx := context.Background()
	queries := c.Request.URL.Query()
	query = queries["query"][0]
	if len(queries["sort_field"]) > 0 {
		sortField = queries["sort_field"][0]
	} else {
		sortField = ""
	}
	ascending, err = strconv.ParseBool(c.Query("ascending"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
	}
	// TODO: accept the other symbols
	targetTypes = strings.Split(queries["target_types"][0], ",")
	indexName := queries["index_name"][0]
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
	search.Index = indexName
	searchInput := &esearch.SearchServiceInput{
		Ctx:          ctx,
		Skip:         skip,
		Take:         take,
		SearchText:   query,
		SortField:    sortField,
		Ascending:    ascending,
		TargetFields: targetTypes,
	}
	res, err := search.SearchMultiMatchQuery(searchInput)
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
