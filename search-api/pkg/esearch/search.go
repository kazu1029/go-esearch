package esearch

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/olivere/elastic"
)

type SearchService struct {
	Client *elastic.Client
}

type DocumentResponse struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SearchResponse struct {
	Time    string             `json:"time"`
	Hits    string             `json:"hits"`
	Results []DocumentResponse `json:"results"`
}

func NewSearchService(Client *elastic.Client) *SearchService {
	return &SearchService{Client: Client}
}

func (s *SearchService) SearchMultiMatchQuery(ctx context.Context, indexName string, skip int, take int, text interface{}, fields ...string) (SearchResponse, error) {
	esQuery := elastic.NewMultiMatchQuery(text, fields...).
		Fuzziness("AUTO").
		MinimumShouldMatch("1")
	result, err := s.Client.Search().
		Index(indexName).
		Query(esQuery).
		From(skip).Size(take).
		Do(ctx)

	res := SearchResponse{
		Time: fmt.Sprintf("%d", result.TookInMillis),
		Hits: fmt.Sprintf("%d", result.Hits.TotalHits),
	}
	if err != nil {
		return res, err
	}
	docs := make([]DocumentResponse, 0)
	for _, hit := range result.Hits.Hits {
		var doc DocumentResponse
		json.Unmarshal(*hit.Source, &doc)
		docs = append(docs, doc)
	}
	res.Results = docs
	return res, nil
}
