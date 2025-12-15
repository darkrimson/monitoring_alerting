package main

import (
	"context"
	"log"
	"os"

	"github.com/darkrimson/monitoring_alerting/internal/repository/postgres"
)

func main() {
	pool, err := postgres.NewPool(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
}
