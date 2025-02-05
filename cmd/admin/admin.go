package main

import (
	"flag"
	"log"
	"os"

	"vpn-tg-bot/internal/storage/sqlite"
	"vpn-tg-bot/web/admin_panel"

	"github.com/joho/godotenv"
)

func main() {

	envPathPtr := flag.String("env", ".env", "path to .env file")
	storagePathPtr := flag.String("storage", "internal/storage/sqlite/test_data/db.db", "path to storage file")
	flag.Parse()

	envPath := *envPathPtr
	log.Printf("env path: %s", envPath)

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file: %v \n env path: %s", err, envPath)
	}

	addr := os.Getenv("ADMIN_PANEL_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	storagePath := *storagePathPtr
	log.Printf("storage path: %s", storagePath)
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		log.Fatalf("[ERR]: Storage not found by given path: %s", storagePath)
	} else if err != nil {
		log.Fatalf("[ERR]: Unexpected error os.Stat: %v", err)
	}

	storage, err := sqlite.New(storagePath)
	if err != nil {
		log.Fatalf("[ERR]: Can't create sqlite storage instance: %v", err)
	}

	sessionKey := os.Getenv("ADMIN_PANEL_SESSION_KEY")
	if sessionKey == "" && len(sessionKey) != 32 {
		log.Fatalf("ADMIN_PANEL_SESSION_KEY is incorrect, value: %s", sessionKey)
	}

	settings := admin_panel.Settings{
		Addr:       addr,
		Storage:    storage,
		SessionKey: sessionKey,
	}
	p := admin_panel.New(settings)

	if err := p.Run(); err != nil {
		log.Fatalf("[ERR]: Can't start admin panel routes: %v", err)
	}
}
