package esearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	Time      string             `json:"time"`
	Hits      string             `json:"hits"`
	Documents []DocumentResponse `json:"documents"`
}

func NewSearchService(Client *elastic.Client) *SearchService {
	return &SearchService{Client: Client}
}

func (s *SearchService) SearchMultiMatchQuery(ctx context.Context, indexName string, skip int, take int, text interface{}, fields ...string) (SearchResponse, error) {
	esQuery := elastic.NewMultiMatchQuery(text, fields...).
		Fuzziness("2").
		MinimumShouldMatch("2")
	log.Printf("esQuery is %+v\n", esQuery)
	result, err := s.Client.Search().
		Index(indexName).
		Query(esQuery).
		From(skip).Size(take).
		Do(ctx)
	log.Printf("result are %+\n", result)

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
	log.Printf("text is %+v\n", text)
	log.Printf("docs are %+v\n", docs)
	log.Printf("fields are %+v\n", fields)
	res.Documents = docs
	return res, nil
}
