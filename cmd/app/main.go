package main

import (
	"context"
	"log"
	"os"

	"github.com/dxssert/lab4KR/internal/storage/postgres"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := postgres.New(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := postgres.NewDefectRepository(pool)

	defects, err := repo.List(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Подключение к БД успешно. Найдено записей брака: %d", len(defects))
}
