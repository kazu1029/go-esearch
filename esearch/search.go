package esearch

import (
	"context"
	"encoding/json"
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

func (s *SearchService) SearchMultiMatchQuery(ctx context.Context, indexName string, skip int, take int, text interface{}, sortField string, ascending bool, fields ...string) (SearchResponse, error) {
	res := SearchResponse{}
	// TODO: check fields are not empty
	esQuery := elastic.NewMultiMatchQuery(text, fields...).
		Fuzziness("AUTO").
		MinimumShouldMatch("1")
	result, err := s.Client.Search().
		Index(indexName).
		Query(esQuery).
		Sort(sortField, ascending).
		From(skip).Size(take).
		Do(ctx)

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
