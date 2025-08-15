package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	telegram_bot "vpn-tg-bot/internal/bot/telegram"
	xui_service "vpn-tg-bot/internal/service/subscription/xui"
	"vpn-tg-bot/internal/service/vpnserver"
	"vpn-tg-bot/internal/storage/sqlite"

	"github.com/joho/godotenv"
	"gopkg.in/telebot.v4"
)

func main() {
	envPathPtr := flag.String("env", ".env", "path to .env file")
	storagePathPtr := flag.String("storage", "internal/storage/sqlite/test_data/db.db", "path to storage file")
	i18nPathPtr := flag.String("i18n", "./internal/bot/telegram/i18n", "path to i18n files")
	flag.Parse()

	envPath := *envPathPtr
	log.Printf("env path: %s", envPath)

	storagePath := *storagePathPtr
	log.Printf("storage path: %s", storagePath)

	i18nPath := *i18nPathPtr
	log.Printf("i18n path: %s", i18nPath)

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file: %v \n env path: %s", err, envPath)
	}

	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		log.Fatalf("[ERR]: Storage not found by given path: %s", storagePath)
	} else if err != nil {
		log.Fatalf("[ERR]: Unexpected error os.Stat: %v", err)
	}

	storage, err := sqlite.New(storagePath)
	if err != nil {
		log.Fatalf("[ERR]: Can't create sqlite storage instance: %v", err)
	}

	auth, err := sqlite.New(storagePath)
	if err != nil {
		log.Fatalf("[ERR]: Can't create sqlite storage instance: %v", err)
	}
	vpnSeverManager := vpnserver.New(storage, auth)

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatalf("[ERR]: TELEGRAM_TOKEN is not set")
	}
	// fmt.Println("Token:", token)
	subsService := xui_service.NewXUIService(xui_service.TokenKey_3x_ui, storage, storage)

	settings := telegram_bot.Settings{
		Settings: telebot.Settings{
			Token:  token,
			Poller: &telebot.LongPoller{Timeout: 10},
			// Verbose: true,
		},
	}

	bot, err := telegram_bot.New(settings, storage, subsService, vpnSeverManager, i18nPath)
	if err != nil {
		log.Fatalf("[ERR]: Can't create bot instance: %v", err)
	}
	log.Println("Bot started.")
	bot.Run()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)

	<-exitChan

	bot.Stop()
	log.Println("Bot stopped.")
	os.Exit(0)

}
