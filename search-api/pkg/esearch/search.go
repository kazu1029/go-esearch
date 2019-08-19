package esearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/olivere/elastic"
)

type SearchService struct {
	Client *elastic.Client
}

type SearchResponse struct {
	Time    string        `json:"time"`
	Hits    string        `json:"hits"`
	Results []interface{} `json:"results"`
}

func NewSearchService(Client *elastic.Client) *SearchService {
	return &SearchService{Client: Client}
}

func (s *SearchService) SearchMultiMatchQuery(ctx context.Context, indexName string, skip int, take int, text interface{}, fields ...string) (SearchResponse, error) {
	res := SearchResponse{}
	// TODO: check fields are not empty
	esQuery := elastic.NewMultiMatchQuery(text, fields...).
		Fuzziness("AUTO").
		MinimumShouldMatch("1")
	result, err := s.Client.Search().
		Index(indexName).
		Query(esQuery).
		From(skip).Size(take).
		Do(ctx)

	res.Time = fmt.Sprintf("%d", result.TookInMillis)
	res.Hits = fmt.Sprintf("%d", result.Hits.TotalHits)
	if err != nil {
		return res, err
	}

	docs := make([]interface{}, 0)
	for _, hit := range result.Hits.Hits {
		var doc interface{}
		json.Unmarshal(*hit.Source, &doc)
		docs = append(docs, doc)
	}
	res.Results = docs
	return res, nil
}
