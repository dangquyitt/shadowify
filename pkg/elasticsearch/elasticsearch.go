package elasticsearch

import (
	"shadowify/pkg/config"

	"github.com/elastic/go-elasticsearch/v9"
)

func NewElasticsearch(cfg *config.Config) (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			cfg.Elasticsearch.URL,
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = es.Ping()
	if err != nil {
		return nil, err
	}

	return es, nil
}
