-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE monitors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    url TEXT NOT NULL,

    interval_seconds INT NOT NULL CHECK (interval_seconds >= 30),
    timeout_seconds INT NOT NULL CHECK (timeout_seconds > 0),
    expected_status INT NOT NULL,

    enabled BOOLEAN NOT NULL DEFAULT true,

    last_status TEXT,
    last_checked_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE checks (
    id BIGSERIAL PRIMARY KEY,
    monitor_id UUID NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,

    ts TIMESTAMP NOT NULL DEFAULT now(),

    status TEXT NOT NULL,
    status_code INT,
    latency_ms INT,
    error TEXT
);

CREATE INDEX idx_checks_monitor_ts
    ON checks (monitor_id, ts DESC);

CREATE TABLE incidents (
    id UUID PRIMARY KEY,
    monitor_id UUID NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,

    started_at TIMESTAMP NOT NULL,
    resolved_at TIMESTAMP,

    current_status TEXT NOT NULL
);

CREATE TABLE alert_targets (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    payload JSONB NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE monitor_alerts (
    monitor_id UUID NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    alert_target_id UUID NOT NULL REFERENCES alert_targets(id) ON DELETE CASCADE,
    PRIMARY KEY (monitor_id, alert_target_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS monitor_alerts;
DROP TABLE IF EXISTS alert_targets;
DROP TABLE IF EXISTS incidents;
DROP TABLE IF EXISTS checks;
DROP TABLE IF EXISTS monitors;

-- +goose StatementEnd
