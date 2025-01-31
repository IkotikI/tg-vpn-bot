package app

import (
	"log"
	"vpn-tg-bot/internal/service/bot/telegram"
	"vpn-tg-bot/internal/storage/sqlite"
	"vpn-tg-bot/pkg/env"

	"github.com/joho/godotenv"
	telebot "gopkg.in/telebot.v4"
)

func Run() {

	envPath := "../../.env"

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect storage
	s, err := sqlite.New(env.MustEnv("DB_PATH"))
	if err != nil {
		log.Fatal("[ERR] can't create sqlite storage:", err)
	}

	// Run Telegram bot
	bot, err := telegram.NewBot(
		telegram.Settings{
			Settings: telebot.Settings{
				Token: env.MustEnv("TELEGRAM_TOKEN"),
				Poller: &telebot.LongPoller{
					Timeout: 10,
				},
			},
		},
	)
	if err != nil {
		log.Fatal("[ERR] can't create telegram bot:", err)
	}
	processor := telegram.New(bot, s)
	processor.RegisterHandlers()

	// Run HTTP REST server
}
