package main

import (
	"log"
	"os"

	"vpn-tg-bot/internal/storage/sqlite"
	"vpn-tg-bot/web/admin_panel"

	"github.com/joho/godotenv"
)

func main() {
	envPath := ".env"

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	addr := os.Getenv("ADMIN_PANEL_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	storagePath := "internal/storage/sqlite/test_data/db.db"
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		log.Fatalf("[ERR]: Storage not found by given path: %s", storagePath)
	} else if err != nil {
		log.Fatalf("[ERR]: Unexpected error os.Stat: %v", err)
	}
	storage, err := sqlite.New(storagePath)
	if err != nil {
		log.Fatalf("[ERR]: Can't create sqlite storage instance: %v", err)
	}

	p := admin_panel.New(addr, storage)

	if err := p.Run(); err != nil {
		log.Fatalf("[ERR]: Can't start admin panel routes: %v", err)
	}
}
