package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pangaunn/testcontainers-go-workshop/pkg/book"
)

type handler struct {
	bookSvc book.BookService
}

func NewHandler(b book.BookService) handler {
	return handler{
		bookSvc: b,
	}
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
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	response, err := h.bookSvc.GetBookByID(c, int64(idInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h handler) DeleteBookByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = h.bookSvc.DeleteByID(c, int64(idInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, "")
}

func (h handler) UpdateBookByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
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
