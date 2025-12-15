package config

import "os"

type DBConfig struct {
	DSN string
}

func LoadDB() DBConfig {
	return DBConfig{
		DSN: os.Getenv("DB_URL"),
	}
}
