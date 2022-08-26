package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pangaunn/testcontainers-go-workshop/cmd/handler"
)

func main() {
	port := ":3000"

	r := gin.Default()
	r.GET("/healthcheck", handler.Test)

	r.Run(port)
}
