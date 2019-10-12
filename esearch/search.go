package esearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/olivere/elastic"
)

type SearchService struct {
	Client    *elastic.Client
	ascending bool
}

type SearchResponse struct {
	Time    string        `json:"time"`
	Hits    string        `json:"hits"`
	Results []interface{} `json:"results"`
}

func NewSearchService(Client *elastic.Client) *SearchService {
	return &SearchService{Client: Client}
}

// TODO: indexName should move to SearchService
func (s *SearchService) SearchMultiMatchQuery(ctx context.Context, indexName string, skip int, take int, text interface{}, sortField string, ascending bool, fields ...string) (SearchResponse, error) {
	res := SearchResponse{}
	s.SetAsc(ascending)
	esQuery := elastic.NewMultiMatchQuery(text, fields...).
		Fuzziness("AUTO").
		MinimumShouldMatch("1")

	result, err := s.Client.Search().
		Index(indexName).
		Query(esQuery).
		SortBy(elastic.NewFieldSort(sortField).Asc()).
		From(skip).Size(take).
		Do(ctx)
	fmt.Printf("result: %v\n", result)

	if result == nil {
		return res, nil
	}

	res.Time = fmt.Sprintf("%d", result.TookInMillis)
	res.Hits = fmt.Sprintf("%d", result.Hits.TotalHits)
	if err != nil {
		return res, err
	}

	hits, _ := strconv.Atoi(res.Hits)
	docs := make([]interface{}, hits)
	for i, hit := range result.Hits.Hits {
		var doc interface{}
		json.Unmarshal(*hit.Source, &doc)
		docs[i] = doc
	}
	res.Results = docs
	return res, nil
}

func (s *SearchService) SetAsc(ascending bool) *SearchService {
	if ascending {
		return s.Asc()
	} else {
		return s.Desc()
	}
}

func (s *SearchService) Asc() *SearchService {
	s.ascending = true
	return s
}

func (s *SearchService) Desc() *SearchService {
	s.ascending = false
	return s
}
