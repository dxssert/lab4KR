package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dxssert/lab4KR/internal/httpapi"
	"github.com/dxssert/lab4KR/internal/service"
	"github.com/dxssert/lab4KR/internal/storage/postgres"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("DATABASE_URL is not set")
	}

	pool, err := postgres.New(databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	repo := postgres.NewDefectRepository(pool)

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = "127.0.0.1:9095"
	}
	server := &http.Server{
		Addr:              addr,
		Handler:           httpapi.NewRouter(service.NewDefects(repo)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
	log.Printf("API запущен на http://%s", addr)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
