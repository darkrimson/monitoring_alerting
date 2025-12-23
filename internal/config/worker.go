package config

import (
	"os"
	"strconv"
)

type WorkerConfig struct {
	TickSeconds      int
	FailureThreshold int
}

func LoadWorker() WorkerConfig {
	return WorkerConfig{
		TickSeconds:      getEnvInt("WORKER_TICK_SECONDS", 1),
		FailureThreshold: getEnvInt("FAILURE_THRESHOLD", 3),
	}
}

func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}

	return i
}
