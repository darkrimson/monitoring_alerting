package main

import (
	"context"
	"log"
	"net/http"

	"github.com/darkrimson/monitoring_alerting/internal/config"
	"github.com/darkrimson/monitoring_alerting/internal/handler"
	"github.com/darkrimson/monitoring_alerting/internal/monitor"
	"github.com/darkrimson/monitoring_alerting/internal/repository/postgres"
	"github.com/darkrimson/monitoring_alerting/internal/router"
)

func main() {
	ctx := context.Background()

	// --- DB ---
	dbCfg := config.LoadDB()
	if dbCfg.DSN == "" {
		log.Fatal("DB_URL is not set")
	}

	pool, err := postgres.NewPool(ctx, dbCfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// --- repositories ---
	monitorRepo := postgres.NewMonitorRepository(pool)

	// --- services ---
	monitorService := monitor.NewMonitorService(monitorRepo)

	// --- handlers ---
	monitorHandler := handler.NewMonitorHandler(monitorService)

	// --- router ---
	r := router.New(router.Handlers{
		Monitor: monitorHandler,
	})

	log.Println("API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
