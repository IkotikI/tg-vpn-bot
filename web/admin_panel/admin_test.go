package admin_panel

import (
	"log"
	"os"
	"testing"
	"vpn-tg-bot/internal/storage/sqlite"

	"github.com/joho/godotenv"
)

func TestRunUntilCancel(t *testing.T) {

	envPath := "../../.env"

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	addr := os.Getenv("ADMIN_PANEL_ADDR")
	if addr == "" {
		addr = ":8088"
	}

	storage, err := sqlite.New("../../internal/storage/sqlite/test_data/db.db")
	if err != nil {
		t.Fatalf("[ERR]: Can't create sqlite storage instance: %v", err)
	}

	p := New(addr, storage)

	if err := p.Run(); err != nil {
		t.Fatalf("[ERR]: Can't start admin panel routes: %v", err)
	}

}
