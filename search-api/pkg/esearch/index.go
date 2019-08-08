package esearch

import (
	"context"
	"errors"
	"time"

	"github.com/olivere/elastic"
	"github.com/teris-io/shortid"
)

type IndexService struct {
	Client    *elastic.Client
	IndexName string
	TypeName  string
}

type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DocumentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
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

func (s *IndexService) BulkInsert(ctx context.Context, docs []DocumentRequest, indexName string, typeName string) (string, error) {
	bulk := s.Client.
		Bulk().
		Index(indexName).
		Type(typeName)

	for _, d := range docs {
		doc := Document{
			ID:        shortid.MustGenerate(),
			Title:     d.Title,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Content:   d.Content,
		}
		bulk.Add(elastic.NewBulkIndexRequest().Id(doc.ID).Doc(doc))
	}
	if _, err := bulk.Do(ctx); err != nil {
		return "", err
	}
	return "Document Bulk Created", nil
}
