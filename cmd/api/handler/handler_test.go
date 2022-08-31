package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pangaunn/testcontainers-go-workshop/pkg/book"
)

var API_URL = "http://localhost:3000"
var _ = Describe("Handler", Label("integration"), func() {

	It("Should return OK 200", func() {
		req, _ := http.NewRequest(http.MethodGet, API_URL+"/healthcheck", nil)
		w := httptest.NewRecorder()
		Engine.ServeHTTP(w, req)

		res, _ := io.ReadAll(w.Body)

		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(string(res)).To(Equal("\"OK\""))
	})

	It("Get Book By ID Should return OK 200", func() {
		req, _ := http.NewRequest(http.MethodGet, API_URL+"/api/v1/book/1", nil)
		w := httptest.NewRecorder()
		Engine.ServeHTTP(w, req)

		res, _ := io.ReadAll(w.Body)
		var b book.BookResponse
		json.Unmarshal(res, &b)
		expected := book.BookResponse{
			ID:          1,
			Name:        "Book 1: Harry Potter and the Sorcerer's Stone",
			Price:       530,
			Author:      "JK rowling",
			Description: "Harry Potter and the Sorcerer's Stone",
			ImageURL:    "http://www.adviceforyou.co.th/blog/wp-content/uploads/2011/12/harry-potter.jpeg",
		}

		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(b).To(Equal(expected))
	})

	It("Create Book Should return OK 201", func() {
		payload := book.NewBookRequest{
			Name:        "the snowman (harry hole)",
			Price:       320,
			Author:      "Jo Nesbø",
			Description: "the snowman",
			ImageURL:    "https://images-na.ssl-images-amazon.com/images/I/51kgjrXdYKL._SX325_BO1,204,203,200_.jpg",
		}

		data, _ := json.Marshal(payload)
		reader := bytes.NewReader(data)

		req, _ := http.NewRequest(http.MethodPost, API_URL+"/api/v1/book", reader)
		w := httptest.NewRecorder()
		Engine.ServeHTTP(w, req)

		expected := book.BookResponse{
			ID:          3,
			Name:        "the snowman (harry hole)",
			Price:       320,
			Author:      "Jo Nesbø",
			Description: "the snowman",
			ImageURL:    "https://images-na.ssl-images-amazon.com/images/I/51kgjrXdYKL._SX325_BO1,204,203,200_.jpg",
		}

		res, _ := io.ReadAll(w.Body)
		var b book.BookResponse
		json.Unmarshal(res, &b)

		Expect(w.Code).To(Equal(http.StatusCreated))
		Expect(b).To(Equal(expected))
	})

	It("Update Book Should return OK 200", func() {
		payload := book.NewBookRequest{
			Name:        "the snowman (harry hole)",
			Price:       320,
			Author:      "Jo Nesbø",
			Description: "the snowman",
			ImageURL:    "https://images-na.ssl-images-amazon.com/images/I/51kgjrXdYKL._SX325_BO1,204,203,200_.jpg",
		}

		data, _ := json.Marshal(payload)
		reader := bytes.NewReader(data)

		req, _ := http.NewRequest(http.MethodPut, API_URL+"/api/v1/book/2", reader)
		w := httptest.NewRecorder()
		Engine.ServeHTTP(w, req)

		expected := book.BookResponse{
			ID:          2,
			Name:        "the snowman (harry hole)",
			Price:       320,
			Author:      "Jo Nesbø",
			Description: "the snowman",
			ImageURL:    "https://images-na.ssl-images-amazon.com/images/I/51kgjrXdYKL._SX325_BO1,204,203,200_.jpg",
		}

		res, _ := io.ReadAll(w.Body)
		var b book.BookResponse
		json.Unmarshal(res, &b)

		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(b).To(Equal(expected))
	})

	It("Delete Book Should return OK 200", func() {
		req, _ := http.NewRequest(http.MethodDelete, API_URL+"/api/v1/book/1", nil)
		w := httptest.NewRecorder()
		Engine.ServeHTTP(w, req)

		res, _ := io.ReadAll(w.Body)

		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(string(res)).To(Equal("\"OK\""))

		getReq, _ := http.NewRequest(http.MethodGet, API_URL+"/api/v1/book/1", nil)
		w2 := httptest.NewRecorder()
		Engine.ServeHTTP(w2, getReq)

		resGet, _ := io.ReadAll(w2.Body)
		var b book.BookResponse
		json.Unmarshal(resGet, &b)

		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(b).To(Equal(book.BookResponse{}))

	})

	It("Create book then search for it should found", func() {
		payload := book.NewBookRequest{
			Name:        "the snowman (harry hole)",
			Price:       320,
			Author:      "Jo Nesbø",
			Description: "the snowman",
			ImageURL:    "https://images-na.ssl-images-amazon.com/images/I/51kgjrXdYKL._SX325_BO1,204,203,200_.jpg",
		}
		data, _ := json.Marshal(payload)
		reader := bytes.NewReader(data)

		req, _ := http.NewRequest(http.MethodPost, API_URL+"/api/v1/book", reader)
		w := httptest.NewRecorder()
		Engine.ServeHTTP(w, req)

		searchReq, _ := http.NewRequest(http.MethodGet, API_URL+"/api/v1/book/search?keyword=snowman", nil)
		w2 := httptest.NewRecorder()
		Engine.ServeHTTP(w2, searchReq)

		res, _ := io.ReadAll(w2.Body)

		var b []book.BookResponse
		json.Unmarshal(res, &b)

		expected := []book.BookResponse{
			{
				ID:          3,
				Name:        "the snowman (harry hole)",
				Price:       320,
				Author:      "Jo Nesbø",
				Description: "the snowman",
				ImageURL:    "https://images-na.ssl-images-amazon.com/images/I/51kgjrXdYKL._SX325_BO1,204,203,200_.jpg",
			},
		}

		Expect(b[0]).To(Equal(expected[0]))
	})
})