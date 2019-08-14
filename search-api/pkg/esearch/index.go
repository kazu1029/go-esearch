package esearch

import (
	"context"
	"errors"
	"time"

	"github.com/olivere/elastic"
)

type IndexService struct {
	Client *elastic.Client
}

type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewIndexService(Client *elastic.Client) *IndexService {
	return &IndexService{Client: Client}
}

func (s *IndexService) CreateMapping(ctx context.Context, indexName string, mapping interface{}) (string, error) {
	exists, err := s.Client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return "", err
	}

	if !exists {
		_, err := s.Client.
			CreateIndex(indexName).
			BodyJson(mapping).
			Do(ctx)
		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New("Index already exists")
	}
	return "Mapping Created", nil
}

func (s *IndexService) CreateIndexTemplate(ctx context.Context, templateName string, template interface{}) (string, error) {
	temp := s.Client.
		IndexPutTemplate(templateName).
		BodyJson(template)

	if _, err := temp.Do(ctx); err != nil {
		return "", err
	}
	return "Template Created", nil
}

func (s *IndexService) BulkInsert(ctx context.Context, docs []interface{}, indexName string, typeName string) (string, error) {
	bulk := s.Client.
		Bulk().
		Index(indexName).
		Type(typeName)

	for _, d := range docs {
		bulk.Add(elastic.NewBulkIndexRequest().Doc(d))
	}
	if _, err := bulk.Do(ctx); err != nil {
		return "", err
	}
	return "Document Bulk Created", nil
}
