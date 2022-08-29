package book

import (
	"database/sql"
	"errors"

	"github.com/pangaunn/testcontainers-go-workshop/pkg/repository"
)

type bookService struct {
	bookRepo repository.BookRepo
}

func NewBookService(b repository.BookRepo) BookService {
	return &bookService{
		bookRepo: b,
	}
}

func (b *bookService) NewBook(data NewBookRequest) (*BookResponse, error) {
	newBook := repository.Book{
		Name:        data.Name,
		Price:       data.Price,
		Author:      data.Author,
		Description: data.Description,
		ImageURL:    data.ImageURL,
	}

	bookCreated, err := b.bookRepo.Create(newBook)
	if err != nil {
		return nil, err
	}

	return &BookResponse{
		ID:          bookCreated.ID,
		Name:        bookCreated.Name,
		Price:       bookCreated.Price,
		Author:      bookCreated.Author,
		Description: bookCreated.Description,
		ImageURL:    bookCreated.ImageURL,
	}, nil
}

func (b *bookService) GetBookByID(id int64) (*BookResponse, error) {
	book, err := b.bookRepo.GetByID(id)
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

func (b *bookService) DeleteByID(id int64) error {
	err := b.bookRepo.DeleteByID(id)
	if err != nil {
		return err
	}
	return nil
}

func (b *bookService) UpdateByID(id int64, book NewBookRequest) (*BookResponse, error) {
	bookUpdated, err := b.bookRepo.Update(id, repository.Book{
		Name:        book.Name,
		Price:       book.Price,
		Author:      book.Author,
		Description: book.Description,
		ImageURL:    book.ImageURL,
	})

	if err != nil {
		return nil, err
	}

	return &BookResponse{
		ID:          id,
		Name:        bookUpdated.Name,
		Price:       bookUpdated.Price,
		Author:      bookUpdated.Author,
		Description: bookUpdated.Description,
		ImageURL:    bookUpdated.ImageURL,
	}, nil
}
