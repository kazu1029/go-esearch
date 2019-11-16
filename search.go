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
	Index        string
	searchSource *elastic.SearchSource
}

type SearchServiceInput struct {
	Ctx          context.Context
	Typ          string
	Skip         int // Skip starts from 0
	Take         int
	SearchText   interface{}
	EsQuery      elastic.Query
	SortField    string
	Ascending    bool
	TargetFields []string
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
	res := SearchResponse{}
	i.EsQuery = elastic.NewMultiMatchQuery(i.SearchText, i.TargetFields...).
		Type(i.Typ).
		Fuzziness("AUTO").
		MinimumShouldMatch("1")

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
	if hits < 50 {
		length = hits
	} else {
		length = i.Take
	}
	docs := make([]interface{}, length)

	for i, hit := range result.Hits.Hits {
		var doc interface{}
		if err := json.Unmarshal(*hit.Source, &doc); err != nil {
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
		Query(i.EsQuery).
		SortBy(sortQuery).
		From(i.Skip).Size(i.Take).
		Do(i.Ctx)

	return res, err
}

func (s *SearchService) SearchWithoutSort(i *SearchServiceInput) (res *elastic.SearchResult, err error) {
	res, err = s.Client.Search().
		Query(i.EsQuery).
		From(i.Skip).Size(i.Take).
		Do(i.Ctx)

	return res, err
}
