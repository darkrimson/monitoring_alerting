package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/darkrimson/monitoring_alerting/internal/handler"
	"github.com/darkrimson/monitoring_alerting/internal/monitor"
	"github.com/darkrimson/monitoring_alerting/internal/repository/postgres"
	"github.com/darkrimson/monitoring_alerting/internal/router"
)

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	// --- DB ---
	pool, err := postgres.NewPool(ctx, dbURL)
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
		Monitor: monitorHandler, // 👈 реализация интерфейса router.MonitorHandler
	})

	log.Println("API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
