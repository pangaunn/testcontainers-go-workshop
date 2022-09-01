package book

import "context"

type NewBookRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Author      string  `json:"author"`
	Description string  `json:"description"`
	ImageURL    string  `json:"imageUrl"`
}

type BookResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Author      string  `json:"author"`
	Description string  `json:"description"`
	ImageURL    string  `json:"imageUrl"`
}

type BookService interface {
	NewBook(ctx context.Context, book NewBookRequest) (*BookResponse, error)
	GetBookByID(ctx context.Context, ID int64) (*BookResponse, error)
	DeleteByID(ctx context.Context, ID int64) error
	UpdateByID(ctx context.Context, ID int64, book NewBookRequest) (*BookResponse, error)
	GetBookByKeyword(ctx context.Context, keyword string) ([]BookResponse, error)
	GetCache(ctx context.Context, keyword string) ([]BookResponse, error)
	SetCache(ctx context.Context, keyword string, books []BookResponse) error
}
