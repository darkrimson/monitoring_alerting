package worker

import (
	"context"
	"log"
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/httpclient"
	"github.com/darkrimson/monitoring_alerting/internal/scheduler"
)

type Worker struct {
	scheduler  *scheduler.Scheduler
	httpClient *httpclient.Client
	checksRepo ChecksRepository
	stateRepo  MonitorStateRepository
	tick       time.Duration
}

func New(
	scheduler *scheduler.Scheduler,
	httpClient *httpclient.Client,
	checksRepo ChecksRepository,
	stateRepo MonitorStateRepository,
	tick time.Duration,
) *Worker {
	return &Worker{
		scheduler:  scheduler,
		httpClient: httpClient,
		checksRepo: checksRepo,
		stateRepo:  stateRepo,
		tick:       tick,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.tick)
	defer ticker.Stop()

	log.Println("worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("worker stopped")
			return

		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

func (w *Worker) runOnce(ctx context.Context) {
	now := time.Now()

	due, err := w.scheduler.DueMonitors(ctx, now)
	if err != nil {
		log.Println("scheduler error:", err)
		return
	}

	log.Printf("scheduler selected %d monitors\n", len(due))

	for _, m := range due {
		result := w.httpClient.Check(ctx, m)

		if err := w.checksRepo.Insert(ctx, result); err != nil {
			log.Println("insert checks error:", err)
			continue
		}

		if err := w.stateRepo.UpdateStatus(
			ctx,
			result.MonitorID,
			string(result.Status),
			result.CheckedAt,
		); err != nil {
			log.Println("update status error:", err)
		}
	}
}
