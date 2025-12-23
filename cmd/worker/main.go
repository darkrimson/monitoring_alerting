package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/alerts"
	"github.com/darkrimson/monitoring_alerting/internal/httpclient"
	"github.com/darkrimson/monitoring_alerting/internal/incidents"
	"github.com/darkrimson/monitoring_alerting/internal/repository/postgres"
	"github.com/darkrimson/monitoring_alerting/internal/scheduler"
	"github.com/darkrimson/monitoring_alerting/internal/worker"
)

func main() {
	ctx := context.Background()

	// ---------- DB ----------
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	pool, err := postgres.NewPool(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// ---------- repositories ----------
	schedulerRepo := postgres.NewSchedulerRepository(pool)
	checksRepo := postgres.NewChecksRepository(pool)
	stateRepo := postgres.NewMonitorStateRepository(pool)

	incidentRepo := postgres.NewIncidentRepository(pool)
	alertRepo := postgres.NewAlertRepository(pool)

	// ---------- core components ----------
	sched := scheduler.NewScheduler(schedulerRepo)
	httpClient := httpclient.NewClient()

	evaluator := incidents.NewEvaluator(1) // failure threshold
	notifier := alerts.NewTelegramNotifier()

	// ---------- worker ----------
	w := worker.New(
		sched,
		httpClient,
		checksRepo,
		stateRepo,
		incidentRepo,
		alertRepo,
		notifier,
		evaluator,
		1*time.Second,
	)

	log.Println("worker starting")
	w.Run(ctx)
}
