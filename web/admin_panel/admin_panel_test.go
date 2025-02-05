package admin_panel

import (
	"fmt"
	"log"
	"os"
	"testing"
	"vpn-tg-bot/internal/storage/sqlite"
	"vpn-tg-bot/pkg/debug"

	"github.com/gorilla/securecookie"

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

	sessionKey := os.Getenv("ADMIN_PANEL_SESSION_KEY")
	if sessionKey == "" && len(sessionKey) != 32 {
		log.Fatalf("ADMIN_PANEL_SESSION_KEY is incorrect, value: %s", sessionKey)
	}

	storage, err := sqlite.New("../../internal/storage/sqlite/test_data/db.db")
	if err != nil {
		t.Fatalf("[ERR]: Can't create sqlite storage instance: %v", err)
	}

	s := Settings{
		Addr:       addr,
		Storage:    storage,
		SessionKey: sessionKey,
	}
	p := New(s)

	if err := p.Run(); err != nil {
		t.Fatalf("[ERR]: Can't start admin panel routes: %v", err)
	}

}

func TestGenerateSessionKey(t *testing.T) {
	key := securecookie.GenerateRandomKey(32)
	fmt.Printf("New cookies key generated \n bytes: %v,\n string: %x\n", key, key)
}

func TestAssets(t *testing.T) {
	p := New(Settings{Addr: ":8080", Storage: nil})
	d, err := os.Getwd()
	fmt.Println("current dir: ", d, err)
	err = debug.DisplayFileSystem(p.Assets())
	if err != nil {
		log.Fatal(err)
	}
}
