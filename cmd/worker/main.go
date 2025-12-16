package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/httpclient"
	"github.com/darkrimson/monitoring_alerting/internal/repository/postgres"
	"github.com/darkrimson/monitoring_alerting/internal/scheduler"
	"github.com/darkrimson/monitoring_alerting/internal/worker"
)

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	pool, err := postgres.NewPool(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// repositories
	schedulerRepo := postgres.NewSchedulerRepository(pool)
	checksRepo := postgres.NewChecksRepository(pool)
	stateRepo := postgres.NewMonitorStateRepository(pool)

	// core components
	sched := scheduler.NewScheduler(schedulerRepo)
	httpClient := httpclient.NewClient()

	w := worker.New(
		sched,
		httpClient,
		checksRepo,
		stateRepo,
		1*time.Second, // tick
	)

	log.Println("worker starting")
	w.Run(ctx)
}
