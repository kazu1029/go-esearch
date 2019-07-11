package handlers

import (
	"github.com/olivere/elastic"
)

var elasticClient *elastic.Client

func InitElastic() (elasticClient *elastic.Client, err error) {
	elasticClient, err = elastic.NewClient(
		elastic.SetURL("http://elasticsearch:9200"),
		elastic.SetSniff(false),
	)
	return
}
