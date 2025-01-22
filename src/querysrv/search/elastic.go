package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/SemyonTolkachyov/amb/src/common/schema"
	elastic "github.com/elastic/go-elasticsearch/v8"
)

type ElasticRepository struct {
	client *elastic.Client
}

func NewElastic(url string) (*ElasticRepository, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, err
	}
	_, err = client.Info()
	if err != nil {
		return nil, err
	}
	return &ElasticRepository{client}, nil
}

func (r *ElasticRepository) Close() {
}

func (r *ElasticRepository) InsertMessage(ctx context.Context, message schema.Message) error {
	body, _ := json.Marshal(message)
	_, err := r.client.Index(
		"messages",
		bytes.NewReader(body),
		r.client.Index.WithDocumentID(message.Id),
		r.client.Index.WithRefresh("wait_for"),
	)
	return err
}

func (r *ElasticRepository) SearchMessages(ctx context.Context, query string, skip uint64, take uint64) (result []schema.Message, err error) {
	var buf bytes.Buffer
	reqBody := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     query,
				"fields":    []string{"body"},
				"fuzziness": "AUTO",
			},
		},
	}
	if err = json.NewEncoder(&buf).Encode(reqBody); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("messages"),
		r.client.Search.WithFrom(int(skip)),
		r.client.Search.WithSize(int(take)),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			result = nil
		}
	}()
	if res.IsError() {
		return nil, errors.New("search failed")
	}

	type Response struct {
		Took int64
		Hits struct {
			Total struct {
				Value int64
			}
			Hits []*struct {
				Source schema.Message `json:"_source"`
			}
		}
	}
	resBody := Response{}
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, err
	}
	var messages []schema.Message
	for _, hit := range resBody.Hits.Hits {
		messages = append(messages, hit.Source)
	}
	return messages, nil
}
