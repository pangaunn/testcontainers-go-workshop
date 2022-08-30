package repository

import (
	"context"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type Book struct {
	ID          int64   `db:"id" json:"id,omitempty"`
	Name        string  `db:"name" json:"name,omitempty"`
	Price       float64 `db:"price" json:"price,omitempty"`
	Author      string  `db:"author" json:"author,omitempty"`
	Description string  `db:"description" json:"description,omitempty"`
	ImageURL    string  `db:"image_url" json:"imageUrl,omitempty"`
}

type BookRepo interface {
	Create(ctx context.Context, b Book) (*Book, error)
	Update(ctx context.Context, b Book) (*Book, error)
	DeleteByID(ctx context.Context, ID int64) error
	GetByID(ctx context.Context, ID int64) (*Book, error)
}

type BookESRepo interface {
	Search(ctx context.Context, keyword string) (*esapi.Response, error)
	Index(ctx context.Context, b Book) (*Book, error)
	Delete(ctx context.Context, ID int64) error
}
