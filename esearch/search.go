package esearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/olivere/elastic"
)

type SearchService struct {
	Client       *elastic.Client
	searchSource *elastic.SearchSource
}

type SearchServiceInput struct {
	Ctx          context.Context
	Index        string
	Typ          string
	Skip         int // Skip starts from 0
	Take         int
	SearchText   interface{}
	EsQuery      elastic.Query
	SortField    string
	Ascending    bool
	TargetFields []string
	TargetTerms  map[string]string
	TargetBools  map[string]bool
}

type SearchResponse struct {
	Time    string        `json:"time"`
	Hits    string        `json:"hits"`
	Results []interface{} `json:"results"`
}

func NewSearchService(Client *elastic.Client) *SearchService {
	return &SearchService{Client: Client, searchSource: elastic.NewSearchSource()}
}

func (s *SearchService) SearchMultiMatchQuery(i *SearchServiceInput) (SearchResponse, error) {
	var result *elastic.SearchResult
	var err error
	var res SearchResponse

	query := elastic.NewMultiMatchQuery(i.SearchText, i.TargetFields...).
		Type(i.Typ).
		Operator("OR")

	i.EsQuery = elastic.NewBoolQuery().Should(query)

	if len(i.SortField) > 0 {
		result, err = s.SearchWithSort(i)
	} else {
		result, err = s.SearchWithoutSort(i)
	}

	if err != nil {
		return res, err
	}
	res.Time = fmt.Sprintf("%d", result.TookInMillis)
	res.Hits = fmt.Sprintf("%d", result.Hits.TotalHits)

	hits, _ := strconv.Atoi(res.Hits)
	var length int
	if hits < 100 {
		length = hits
	} else {
		length = i.Take
	}
	docs := make([]interface{}, length)

	for i, doc := range docs {
		err := json.Unmarshal(*result.Hits.Hits[i].Source, &doc)
		if err != nil {
			return res, err
		}
		docs[i] = doc
	}

	res.Results = docs
	return res, nil
}

func (s *SearchService) SearchWithSort(i *SearchServiceInput) (res *elastic.SearchResult, err error) {
	var sortQuery *elastic.FieldSort
	if i.Ascending {
		sortQuery = elastic.NewFieldSort(i.SortField).Asc()
	} else {
		sortQuery = elastic.NewFieldSort(i.SortField).Desc()
	}

	res, err = s.Client.Search().
		Index(i.Index).
		Query(i.EsQuery).
		SortBy(sortQuery).
		From(i.Skip).Size(i.Take).
		Do(i.Ctx)

	return res, err
}

func (s *SearchService) SearchWithoutSort(i *SearchServiceInput) (res *elastic.SearchResult, err error) {
	res, err = s.Client.Search().
		Index(i.Index).
		Query(i.EsQuery).
		From(i.Skip).Size(i.Take).
		Do(i.Ctx)

	return res, err
}
