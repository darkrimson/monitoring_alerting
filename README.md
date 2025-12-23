# Monitoring & Alerting (MVP)

Minimal monitoring and alerting service for HTTP endpoints.

This project is an MVP implementation of a monitoring system that:
- periodically checks HTTP endpoints
- detects incidents based on consecutive failures
- creates alerts on state changes
- delivers notifications via Telegram

The focus of this project is **correct domain modeling**, not UI or auth.

---

## ✨ Features

- HTTP endpoint monitoring
- Configurable check interval and timeout
- Incident lifecycle (OPEN → RESOLVED)
- Failure threshold to avoid flapping
- Alert creation only on state change
- Telegram notifications
- Docker-based setup
- PostgreSQL storage
- Clean architecture with explicit boundaries

---

## 🧠 Core Concepts

### Monitor
A monitored HTTP endpoint with:
- URL
- check interval
- timeout
- expected HTTP status

### Check
A single execution result of a monitor:
- status (UP / DOWN)
- latency
- status code
- timestamp

### Incident
Represents a **stable failure**, not a single error.
An incident is created only after *N consecutive failures*.

Only one OPEN incident can exist per monitor.

### Alert
Created when an incident is:
- opened
- resolved

Alerts are delivered asynchronously via notifiers.

---

## 🏗 Architecture Overview

- **API**: manages monitors (CRUD)
- **Worker**:
    - schedules checks
    - performs HTTP requests
    - evaluates incidents
    - creates alerts
    - sends notifications
- **PostgreSQL**: single source of truth
- **Telegram**: notification channel

---

🚀 Running the project
Start services

```dockerfile
docker-compose up --build
```
Add monitor

```shell
 curl -X POST http://localhost:8080/api/v1/monitors   -H "Content-Type: application/json"   -d '{
    "name": "Test",
    "url": "https://example.com",
    "interval_seconds": 30,
    "timeout_seconds": 5,
    "expected_status": 200,
    "enabled": true
  }'
```

🧪 Example workflow

- Create a monitor via API
- Worker starts checking the endpoint
- Endpoint fails N times consecutively
- Incident is opened
- Alert is created and sent to Telegram
- Endpoint recovers
- Incident is resolved
- Recovery alert is sent
---

📌 Project status

MVP completed

Implemented features cover the core monitoring and alerting loop.
The project is intentionally limited in scope.

- Possible future improvements:
- multiple notifiers (email, Slack)
- UI dashboard
- parallel workers with locking
- SLA / SLO metrics
- authentication