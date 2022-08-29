package book

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
	NewBook(NewBookRequest) (*BookResponse, error)
	GetBookByID(id int64) (*BookResponse, error)
	DeleteByID(id int64) error
	UpdateByID(id int64, book NewBookRequest) (*BookResponse, error)
}
