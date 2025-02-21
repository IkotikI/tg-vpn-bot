package test

import (
	"log"
	"os"
	"testing"
	"vpn-tg-bot/internal/storage/sqlite"
	"vpn-tg-bot/web/admin_panel"

	"github.com/joho/godotenv"
)

const basePath = "../../../"

func getTestSettings(t *testing.T) *admin_panel.Settings {
	envPath := basePath + ".env"

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	addr := os.Getenv("ADMIN_PANEL_ADDR")
	if addr == "" {
		addr = ":8088"
	}

	sessionKey := os.Getenv("ADMIN_PANEL_SESSION_KEY")
	if sessionKey == "" && len(sessionKey) != 32 {
		t.Fatalf("ADMIN_PANEL_SESSION_KEY is incorrect, value: %s", sessionKey)
	}

	storage, err := sqlite.New(basePath + "internal/storage/sqlite/test_data/db.db")
	if err != nil {
		t.Fatalf("[ERR]: Can't create sqlite storage instance: %v", err)
	}

	return &admin_panel.Settings{
		Addr:       addr,
		Scheme:     "http",
		Storage:    storage,
		SessionKey: sessionKey,
	}
}
