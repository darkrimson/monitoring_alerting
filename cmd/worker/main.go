package main

import (
	"context"
	"log"
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/alerts"
	"github.com/darkrimson/monitoring_alerting/internal/config"
	"github.com/darkrimson/monitoring_alerting/internal/httpclient"
	"github.com/darkrimson/monitoring_alerting/internal/incidents"
	"github.com/darkrimson/monitoring_alerting/internal/repository/postgres"
	"github.com/darkrimson/monitoring_alerting/internal/scheduler"
	"github.com/darkrimson/monitoring_alerting/internal/worker"
)

func main() {
	ctx := context.Background()

	// ---------- DB ----------
	dbCfg := config.LoadDB()
	if dbCfg.DSN == "" {
		log.Fatal("DB_URL is not set")
	}

	pool, err := postgres.NewPool(ctx, dbCfg.DSN)
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

	workerCfg := config.LoadWorker()
	tgCfg := config.LoadTelegram()

	evaluator := incidents.NewEvaluator(workerCfg.FailureThreshold) // failure threshold
	notifier := alerts.NewTelegramNotifier(tgCfg)

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
		time.Duration(workerCfg.TickSeconds)*time.Second,
	)

	log.Println("worker starting")
	w.Run(ctx)
}
