-- +goose Up
ALTER TABLE monitors
    ADD COLUMN failure_streak INT NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE monitors
DROP COLUMN failure_streak;