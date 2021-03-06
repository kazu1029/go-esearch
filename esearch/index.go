package esearch

import (
	"context"
	"errors"

	"github.com/olivere/elastic"
)

type IndexService struct {
	Client *elastic.Client
	// TODO: fix later
	// Index string
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

func (s *IndexService) BulkInsert(ctx context.Context, docs []interface{}, indexName, typeName, pipeline string) (string, error) {
	bulk := s.Client.
		Bulk().
		Index(indexName).
		Type(typeName).
		Pipeline(pipeline)

	for _, d := range docs {
		bulk.Add(elastic.NewBulkIndexRequest().Doc(d))
	}
	if _, err := bulk.Do(ctx); err != nil {
		return "", err
	}
	return "Bulk Updated", nil
}
