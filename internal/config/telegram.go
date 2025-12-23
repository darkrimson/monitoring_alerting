package config

import "os"

type TelegramConfig struct {
	Token  string
	ChatID string
}

func LoadTelegram() TelegramConfig {
	return TelegramConfig{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatID: os.Getenv("TELEGRAM_CHAT_ID"),
	}
}
