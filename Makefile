include .env
export

SHELL := /bin/bash

.PHONY: build run migrate test lint help

help:
	@echo "Targets: build, run, migrate, test, lint"

build:
	@echo "build (not implemented yet)"

run:
	@echo "setup"

migrate:
	@echo "Running database migrations..."
	@goose -dir migrations postgres "$(DB_URL)" up

test:
	@echo "test (not implemented yet)"

lint:
	@echo "lint (not implemented yet)"
