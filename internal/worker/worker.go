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

		// INCIDENTS LOGIC

		openIncident, err := w.incidentRepo.GetOpenByMonitor(ctx, result.MonitorID)
		if err != nil {
			log.Println("get open incident error:", err)
			continue
		}

		if openIncident == nil {
			log.Println("NO OPEN INCIDENT")
		} else {
			log.Printf("OPEN INCIDENT id=%s failures=%d",
				openIncident.ID,
				openIncident.FailureCount,
			)
		}

		input := incidents.EvaluateInput{
			HasOpenIncident: openIncident != nil,
			FailureCount:    0,
			CheckStatus:     string(result.Status),
		}

		if openIncident != nil {
			input.FailureCount = openIncident.FailureCount
		}

		decision := w.evaluator.Evaluate(input)

		log.Printf(
			"EVAL monitor=%s status=%s hasOpen=%v failures=%d",
			result.MonitorID,
			result.Status,
			input.HasOpenIncident,
			input.FailureCount,
		)

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

			log.Println("TRY CREATE INCIDENT")

			if err := w.incidentRepo.CreateIncident(ctx, incident); err != nil {
				log.Println("create incident error:", err)
			}

			log.Println("INCIDENT CREATED")

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
