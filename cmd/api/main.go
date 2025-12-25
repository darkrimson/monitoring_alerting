package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/config"
	"github.com/darkrimson/monitoring_alerting/internal/handler"
	"github.com/darkrimson/monitoring_alerting/internal/monitor"
	"github.com/darkrimson/monitoring_alerting/internal/repository/postgres"
	"github.com/darkrimson/monitoring_alerting/internal/router"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Println("starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down server")

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("could not shutdown server gracefully: %v", err)
	}
}
