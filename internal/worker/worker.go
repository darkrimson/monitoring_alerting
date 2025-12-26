package worker

import (
	"context"
	"log"
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/alerts"
	"github.com/google/uuid"

	"github.com/darkrimson/monitoring_alerting/internal/httpclient"
	"github.com/darkrimson/monitoring_alerting/internal/incidents"
	"github.com/darkrimson/monitoring_alerting/internal/models"
	"github.com/darkrimson/monitoring_alerting/internal/scheduler"
)

type Worker struct {
	scheduler    *scheduler.Scheduler
	httpClient   *httpclient.Client
	checksRepo   ChecksRepository
	stateRepo    MonitorStateRepository
	incidentRepo incidents.Repository
	alertRepo    alerts.Repository
	notifier     alerts.Notifier
	evaluator    *incidents.Evaluator
	tick         time.Duration
}

func New(
	scheduler *scheduler.Scheduler,
	httpClient *httpclient.Client,
	checksRepo ChecksRepository,
	stateRepo MonitorStateRepository,
	incidentRepo incidents.Repository,
	alertRepo alerts.Repository,
	notifier alerts.Notifier,
	evaluator *incidents.Evaluator,
	tick time.Duration,
) *Worker {
	return &Worker{
		scheduler:    scheduler,
		httpClient:   httpClient,
		checksRepo:   checksRepo,
		stateRepo:    stateRepo,
		incidentRepo: incidentRepo,
		alertRepo:    alertRepo,
		notifier:     notifier,
		evaluator:    evaluator,
		tick:         tick,
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

		checkID, err := w.checksRepo.Insert(ctx, result)
		if err != nil {
			log.Println("insert check error:", err)
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

		if result.Status == "DOWN" {
			_ = w.stateRepo.IncrementFailureStreak(ctx, result.MonitorID)
		} else {
			_ = w.stateRepo.ResetFailureStreak(ctx, result.MonitorID)
		}

		// INCIDENTS LOGIC

		openIncident, err := w.incidentRepo.GetOpenByMonitor(ctx, result.MonitorID)
		if err != nil {
			log.Println("get open incident error:", err)
			continue
		}

		if openIncident == nil {
			log.Println("NO OPEN INCIDENT")
		}

		input := incidents.EvaluateInput{
			HasOpenIncident: openIncident != nil,
			FailureCount:    m.FailureStreak,
			CheckStatus:     string(result.Status),
		}

		if openIncident != nil {
			input.FailureCount = openIncident.FailureCount
		}

		currentFailures := input.FailureCount
		if result.Status == "DOWN" {
			currentFailures++
		}

		log.Printf(
			"EVAL monitor=%s status=%s hasOpen=%v failures=%d",
			result.MonitorID,
			result.Status,
			openIncident != nil,
			currentFailures,
		)

		decision := w.evaluator.Evaluate(input)

		switch decision.Type {

		case incidents.DecisionNoop:

		case incidents.DecisionOpen:
			incident := &models.Incident{
				ID:           uuid.New(),
				MonitorID:    result.MonitorID,
				Status:       "OPEN",
				StartedAt:    result.CheckedAt,
				FailureCount: w.evaluator.FailureThreshold,
				LastCheckID:  &checkID,
			}

			if err := w.incidentRepo.CreateIncident(ctx, incident); err != nil {
				log.Println("create incident error:", err)
				break
			}

			log.Printf(
				"OPEN INCIDENT id=%s failures=%d",
				incident.ID,
				currentFailures,
			)

			_ = w.stateRepo.ResetFailureStreak(ctx, result.MonitorID)

			alert := &models.Alert{
				ID:         uuid.New(),
				IncidentID: incident.ID,
				Type:       "INCIDENT_OPENED",
				Channel:    "TELEGRAM",
				Payload:    alerts.BuildIncidentPayload(incident),
			}

			if err := w.alertRepo.Create(ctx, alert); err != nil {
				log.Println("create alert error:", err)
			}

		case incidents.DecisionUpdate:
			if openIncident == nil {
				break
			}

			if err := w.incidentRepo.UpdateFailure(
				ctx,
				openIncident.ID,
				checkID,
			); err != nil {
				log.Println("update incident error:", err)
			}

			log.Printf(
				"UPDATE INCIDENT id=%s failures=%d",
				openIncident.ID,
				currentFailures,
			)

		case incidents.DecisionResolve:
			if openIncident == nil {
				break
			}

			if err := w.incidentRepo.ResolveIncident(
				ctx,
				openIncident.ID,
				checkID,
				result.CheckedAt,
			); err != nil {
				log.Println("resolve incident error:", err)
				break
			}

			log.Printf(
				"RESOLVE INCIDENT id=%s",
				openIncident.ID,
			)

			alert := &models.Alert{
				ID:         uuid.New(),
				IncidentID: openIncident.ID,
				Type:       "INCIDENT_RESOLVED",
				Channel:    "TELEGRAM",
				Payload:    alerts.BuildIncidentPayload(openIncident),
			}

			if err := w.alertRepo.Create(ctx, alert); err != nil {
				log.Println("create alert error:", err)
			}
		}
	}

	pending, err := w.alertRepo.GetPending(ctx)
	if err != nil {
		log.Println("get pending alerts error:", err)
		return
	}

	for _, alert := range pending {
		if err := w.notifier.Send(ctx, alert); err != nil {
			log.Println("send alert error:", err)
			continue
		}

		if err := w.alertRepo.MarkSent(ctx, alert.ID); err != nil {
			log.Println("mark alert sent error:", err)
		}
	}
}
