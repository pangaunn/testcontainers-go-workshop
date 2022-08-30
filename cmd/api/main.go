package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pangaunn/testcontainers-go-workshop/cmd/api/handler"
	"github.com/pangaunn/testcontainers-go-workshop/pkg/book"
	"github.com/pangaunn/testcontainers-go-workshop/pkg/datastore"
	"github.com/pangaunn/testcontainers-go-workshop/pkg/repository"
	logger "github.com/sirupsen/logrus"
)

type databaseCredential struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

func init() {
	if currentEnvironment, ok := os.LookupEnv("ENV"); ok {
		if currentEnvironment == "dev" {
			err := godotenv.Load("./.env")
			if err != nil {
				logger.Info("Can't load .env", err)
			}
		}
	}
}

func main() {
	dbCredential := databaseCredential{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	connStr := generateMysqlConnectionString(dbCredential)
	sqlConn := datastore.InitMySQL(connStr)

	cfg := elasticsearch.Config{Addresses: []string{os.Getenv("ELASTICSEARCH_API_ENDPOINT")}}
	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Fatal("elasticsearch.NewClient Error: ", err)
	}

	bookRepo := repository.NewBookRepo(sqlConn)
	bookESRepo := repository.NewBookESRepo(esClient, time.Second*5)
	bookSvc := book.NewBookService(bookRepo, bookESRepo)
	bookHandler := handler.NewHandler(bookSvc)

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

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	logger.Info("Server exiting")

}

func generateMysqlConnectionString(cred databaseCredential) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cred.Username, cred.Password, cred.Host, cred.Port, cred.Name)
}
