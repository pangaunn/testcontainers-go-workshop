package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type bookRepo struct {
	db *sqlx.DB
}

func NewBookRepo(db *sqlx.DB) BookRepo {
	return &bookRepo{
		db: db,
	}
}

func (br *bookRepo) Create(ctx context.Context, b Book) (*Book, error) {
	query := `insert into book (name, price, author, description, image_url) values (?,?,?,?,?)`
	result, err := br.db.ExecContext(ctx, query, b.Name, b.Price, b.Author, b.Description, b.ImageURL)
	if err != nil {
		return nil, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	b.ID = ID
	return &b, nil

}

func (br *bookRepo) GetByID(ctx context.Context, ID int64) (*Book, error) {
	var book Book
	err := br.db.Get(&book, `select id, name, price, author, description, image_url from book where id= ? `, ID)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (br *bookRepo) DeleteByID(ctx context.Context, ID int64) error {
	_, err := br.db.ExecContext(ctx, `delete from book where id=?`, ID)
	if err != nil {
		return err
	}
	return nil
}

func (br *bookRepo) Update(ctx context.Context, b Book) (*Book, error) {
	_, err := br.db.NamedExecContext(ctx, `update book SET name=:name, price=:price, author=:author, description=:description, image_url=:image_url where id=:id`, &b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
