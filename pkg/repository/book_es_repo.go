package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strconv"
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
func (bes *bookESRepo) Search(ctx context.Context, text string) (interface{}, error) {
	panic("unimplemented")
}
