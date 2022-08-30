package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	logger "github.com/sirupsen/logrus"
)

type bookESRepo struct {
	esClient      *elasticsearch.Client
	timeoutSecond time.Duration
}

var BOOKS_TEMP_INDEX = "books_temp"

func NewBookESRepo(esClient *elasticsearch.Client, timeoutSecond time.Duration) BookESRepo {
	return &bookESRepo{
		esClient:      esClient,
		timeoutSecond: timeoutSecond,
	}
}

// es index api will fully replace an existing document
func (bes *bookESRepo) Index(ctx context.Context, b Book) (*Book, error) {
	data, err := json.Marshal(b)
	if err != nil {
		return nil, errors.New("cannot marshal with book")
	}

	bookID := strconv.Itoa(int(b.ID))
	req := esapi.IndexRequest{
		Index:      BOOKS_TEMP_INDEX,
		DocumentID: bookID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, bes.esClient)
	if err != nil {
		logger.Warnf("Can't index elasticsearch")
		return nil, err
	}
	defer res.Body.Close()
	return nil, err

}
func (bes *bookESRepo) Search(ctx context.Context, keyword string) (*esapi.Response, error) {
	query := fmt.Sprintf(`
	{
		"query": {
				"multi_match": {
					"query": "%s",
					"operator" : "and",
					"fields": [
						"name",
						"author",
						"description"
					]
				}
			}
	}
	`, escapeChar(keyword))

	res, err := bes.esClient.Search(
		bes.esClient.Search.WithContext(context.Background()),
		bes.esClient.Search.WithIndex(BOOKS_TEMP_INDEX),
		bes.esClient.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil || res.IsError() {
		logger.Warn("bes.esClient.Search Error:", err)
		return nil, err
	}

	return res, err
}

func (bes *bookESRepo) Delete(ctx context.Context, ID int64) error {
	bookID := strconv.Itoa(int(ID))
	req := esapi.DeleteRequest{
		Index:      BOOKS_TEMP_INDEX,
		DocumentID: bookID,
		Refresh:    "true",
	}

	_, err := req.Do(ctx, bes.esClient)

	if err != nil {
		logger.Warnf("Can't delete elasticsearch")
	}

	return err
}

func escapeChar(str string) string {
	bs, err := json.Marshal(str)
	if err != nil {
		logger.Warnf("cannot marshal string: %s", str)
		return str
	}
	str = string(bs)
	return str[1 : len(str)-1]
}
