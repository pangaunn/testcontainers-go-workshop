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

	"github.com/joho/godotenv"
	"github.com/pangaunn/testcontainers-go-workshop/cmd/api/handler"
	"github.com/pangaunn/testcontainers-go-workshop/pkg/datastore"
	logger "github.com/sirupsen/logrus"
)

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
	dbCredential := datastore.DatabaseCredential{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	engine := handler.InitHandler(dbCredential, os.Getenv("ELASTICSEARCH_API_ENDPOINT"))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		Handler: engine,
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
