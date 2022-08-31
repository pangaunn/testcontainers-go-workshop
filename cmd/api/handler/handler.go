package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/pangaunn/testcontainers-go-workshop/pkg/book"
	"github.com/pangaunn/testcontainers-go-workshop/pkg/datastore"
	"github.com/pangaunn/testcontainers-go-workshop/pkg/repository"
	logger "github.com/sirupsen/logrus"
)

type handler struct {
	bookSvc book.BookService
}

func NewHandler(b book.BookService) handler {
	return handler{
		bookSvc: b,
	}
}

func InitHandler(cre datastore.DatabaseCredential, esURL string, redisURL string) *gin.Engine {
	connStr := datastore.GenerateMysqlConnectionString(cre)
	sqlConn := datastore.InitMySQL(connStr)

	cfg := elasticsearch.Config{Addresses: []string{esURL}}
	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Fatal("elasticsearch.NewClient Error: ", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	bookRepo := repository.NewBookRepo(sqlConn)
	bookESRepo := repository.NewBookESRepo(esClient, time.Second*5)
	bookSvc := book.NewBookService(bookRepo, bookESRepo, redisClient)
	bookHandler := NewHandler(bookSvc)

	r := gin.Default()
	r.GET("/healthcheck", bookHandler.Healthcheck)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/book/:id", bookHandler.GetBookByID)
		v1.POST("/book", bookHandler.NewBook)
		v1.PUT("/book/:id", bookHandler.UpdateBookByID)
		v1.DELETE("/book/:id", bookHandler.DeleteBookByID)
		v1.GET("/book/search", bookHandler.SearchBook)
	}

	return r
}

func (h handler) Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}

func (h handler) NewBook(c *gin.Context) {
	var bookRequest book.NewBookRequest
	if err := c.ShouldBindJSON(&bookRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	response, err := h.bookSvc.NewBook(c, bookRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h handler) GetBookByID(c *gin.Context) {
	ID := c.Param("id")
	idInt, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	response, err := h.bookSvc.GetBookByID(c, int64(idInt))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, "book not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h handler) DeleteBookByID(c *gin.Context) {
	ID := c.Param("id")
	idInt, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = h.bookSvc.DeleteByID(c, int64(idInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "OK")
}

func (h handler) UpdateBookByID(c *gin.Context) {
	ID := c.Param("id")
	idInt, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	var bookRequest book.NewBookRequest
	if err := c.ShouldBindJSON(&bookRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	response, err := h.bookSvc.UpdateByID(c, int64(idInt), bookRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h handler) SearchBook(c *gin.Context) {
	keyword := c.Query("keyword")

	cacheBooks, err := h.bookSvc.GetCache(c, keyword)
	if err != nil {
		logger.Warn("cannot search keyword from redis")
	}

	if len(cacheBooks) > 0 {
		c.JSON(http.StatusOK, cacheBooks)
		return
	}

	books, err := h.bookSvc.GetBookByKeyword(c, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = h.bookSvc.SetCache(c, keyword, books)
	if err != nil {
		logger.Warn("cannot set keyword from redis")
	}

	c.JSON(http.StatusOK, books)
}
