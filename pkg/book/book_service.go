package book

import (
	"context"
	"encoding/json"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/go-redis/redis"
	"github.com/pangaunn/testcontainers-go-workshop/pkg/repository"
	logger "github.com/sirupsen/logrus"
)

type SearchResult struct {
	Hits SearchResultHits `json:"hits"`
}

type SearchResultHits struct {
	Hits ESHits `json:"hits"`
}

type ESHit struct {
	Source BookResponse `json:"_source"`
}

type bookService struct {
	bookRepo    repository.BookRepo
	bookESRepo  repository.BookESRepo
	redisClient *redis.Client
}

type ESHits []ESHit

func (esh ESHits) ToParseBookReponseFromES() []BookResponse {
	bs := []BookResponse{}
	for _, val := range esh {
		bs = append(bs, val.Source)
	}
	return bs
}

func NewBookService(bookRepo repository.BookRepo, bookESRepo repository.BookESRepo, redisClient *redis.Client) BookService {
	return &bookService{
		bookRepo:    bookRepo,
		bookESRepo:  bookESRepo,
		redisClient: redisClient,
	}
}

func (b *bookService) captureToESDataStore(ctx context.Context, book *repository.Book) {
	_, err := b.bookESRepo.Index(ctx, *book)
	if err != nil {
		logger.Warn("cannot index book to elasticsearch")
	}
}

func (b *bookService) NewBook(ctx context.Context, data NewBookRequest) (*BookResponse, error) {
	newBook := repository.Book{
		Name:        data.Name,
		Price:       data.Price,
		Author:      data.Author,
		Description: data.Description,
		ImageURL:    data.ImageURL,
	}

	bookCreated, err := b.bookRepo.Create(ctx, newBook)
	if err != nil {
		return nil, err
	}

	b.captureToESDataStore(ctx, bookCreated)
	return &BookResponse{
		ID:          bookCreated.ID,
		Name:        bookCreated.Name,
		Price:       bookCreated.Price,
		Author:      bookCreated.Author,
		Description: bookCreated.Description,
		ImageURL:    bookCreated.ImageURL,
	}, nil
}

func (b *bookService) GetBookByID(ctx context.Context, ID int64) (*BookResponse, error) {
	book, err := b.bookRepo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return &BookResponse{
		ID:          book.ID,
		Name:        book.Name,
		Price:       book.Price,
		Author:      book.Author,
		Description: book.Description,
		ImageURL:    book.ImageURL,
	}, nil
}

func (b *bookService) DeleteByID(ctx context.Context, ID int64) error {
	err := b.bookRepo.DeleteByID(ctx, ID)
	if err != nil {
		return err
	}

	err = b.bookESRepo.Delete(ctx, ID)
	if err != nil {
		return err
	}

	return nil
}

func (b *bookService) UpdateByID(ctx context.Context, ID int64, book NewBookRequest) (*BookResponse, error) {
	newBookUpdated := repository.Book{
		ID:          ID,
		Name:        book.Name,
		Price:       book.Price,
		Author:      book.Author,
		Description: book.Description,
		ImageURL:    book.ImageURL,
	}
	bookUpdated, err := b.bookRepo.Update(ctx, newBookUpdated)
	if err != nil {
		return nil, err
	}

	b.captureToESDataStore(ctx, bookUpdated)
	return &BookResponse{
		ID:          ID,
		Name:        bookUpdated.Name,
		Price:       bookUpdated.Price,
		Author:      bookUpdated.Author,
		Description: bookUpdated.Description,
		ImageURL:    bookUpdated.ImageURL,
	}, nil
}

func (b *bookService) GetBookByKeyword(ctx context.Context, keyword string) ([]BookResponse, error) {
	result, err := b.bookESRepo.Search(ctx, keyword)
	if err != nil {
		logger.Warn("cannot search book to elasticsearch")
	}

	books := ParseESToBookResponse(result)

	return books, err
}

func (b *bookService) GetCache(ctx context.Context, keyword string) ([]BookResponse, error) {
	cache, err := b.redisClient.Get(keyword).Result()
	if err != nil {
		if err == redis.Nil {
			return []BookResponse{}, nil
		}
		logger.Warn("error b.redisClient.Get", err)
	}

	var br []BookResponse
	err = json.Unmarshal([]byte(cache), &br)
	if err != nil {
		logger.Warn("error json unmarshal", err)
	}

	return br, err
}

func (b *bookService) SetCache(ctx context.Context, keyword string, books []BookResponse) error {
	data, err := json.Marshal(books)
	if err != nil {
		logger.Warn("cache redis marshal ", err)
	}
	err = b.redisClient.Set(keyword, string(data), time.Second*3000).Err()
	if err != nil {
		logger.Warn("error b.redisClient.Set", err)
	}

	return err
}

func ParseESToBookResponse(es *esapi.Response) []BookResponse {
	var s SearchResult
	if err := json.NewDecoder(es.Body).Decode(&s); err != nil {
		logger.Warnf("Error parsing the response body from elasticsearch: %s", err)
	}

	books := s.Hits.Hits.ToParseBookReponseFromES()
	return books
}
