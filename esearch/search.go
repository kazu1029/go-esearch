package esearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

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

func (s *SearchService) SearchMultiMatchQuery(ctx context.Context, indexName string, skip int, take int, text interface{}, sortField string, ascending bool, fields ...string) (SearchResponse, err error) {
	res := SearchResponse{}
	// TODO: check fields are not empty
	fmt.Printf("ascending: %v\n", ascending)
	fmt.Printf("sortField: %v\n", sortField)
	result := elastic.SearchResult{}
	esQuery := elastic.NewMultiMatchQuery(text, fields...).
		Fuzziness("AUTO").
		MinimumShouldMatch("1")
	if ascending {
		result, err = s.Client.Search().
			Index(indexName).
			Query(esQuery).
			SortBy(elastic.NewFieldSort(sortField).Asc(), elastic.NewScoreSort()).
			From(skip).Size(take).
			Do(ctx)
	} else {
		result, err = s.Client.Search().
			Index(indexName).
			Query(esQuery).
			SortBy(elastic.NewFieldSort(sortField).Desc(), elastic.NewScoreSort()).
			From(skip).Size(take).
			Do(ctx)
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
