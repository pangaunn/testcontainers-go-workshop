package book

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pangaunn/testcontainers-go-workshop/pkg/repository"
	logger "github.com/sirupsen/logrus"
)

type bookService struct {
	bookRepo   repository.BookRepo
	bookESRepo repository.BookESRepo
}

func NewBookService(bookRepo repository.BookRepo, bookESRepo repository.BookESRepo) BookService {
	return &bookService{
		bookRepo:   bookRepo,
		bookESRepo: bookESRepo,
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

func (b *bookService) GetBookByID(ctx context.Context, id int64) (*BookResponse, error) {
	book, err := b.bookRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("book not found")
		}
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

func (b *bookService) DeleteByID(ctx context.Context, id int64) error {
	err := b.bookRepo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}

	err = b.bookESRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (b *bookService) UpdateByID(ctx context.Context, id int64, book NewBookRequest) (*BookResponse, error) {
	newBookUpdated := repository.Book{
		Name:        book.Name,
		Price:       book.Price,
		Author:      book.Author,
		Description: book.Description,
		ImageURL:    book.ImageURL,
	}
	bookUpdated, err := b.bookRepo.Update(ctx, id, newBookUpdated)
	if err != nil {
		return nil, err
	}

	b.captureToESDataStore(ctx, bookUpdated)
	return &BookResponse{
		ID:          id,
		Name:        bookUpdated.Name,
		Price:       bookUpdated.Price,
		Author:      bookUpdated.Author,
		Description: bookUpdated.Description,
		ImageURL:    bookUpdated.ImageURL,
	}, nil
}

func (b *bookService) DeleteToESDataStore(ctx context.Context, id int64) error {
	err := b.bookESRepo.Delete(ctx, id)
	if err != nil {
		logger.Warn("cannot delete book to elasticsearch")
	}
	return err
}
