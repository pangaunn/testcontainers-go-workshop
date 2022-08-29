package repository

type Book struct {
	ID          int64   `db:"id"`
	Name        string  `db:"name"`
	Price       float64 `db:"price"`
	Author      string  `db:"author"`
	Description string  `db:"description"`
	ImageURL    string  `db:"image_url"`
}

type BookRepo interface {
	Create(b Book) (*Book, error)
	Update(id int64, b Book) (*Book, error)
	DeleteByID(id int64) error
	GetByID(id int64) (*Book, error)
}
