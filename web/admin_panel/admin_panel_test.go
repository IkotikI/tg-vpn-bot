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

	s := getTestSettings(t)
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
	s := New(Settings{Addr: ":8080", Storage: nil})
	d, err := os.Getwd()
	t.Log("\n current dir: ", d, err)
	err = debug.DisplayFileSystem(s.Assets())
	if err != nil {
		log.Fatal(err)
	}

	// // err = s.collectFileMetadata()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// t.Logf("\n public assets meta: %+v \n", publicAssetsMeta)

}

const basePath = "../../"

func getTestSettings(t *testing.T) Settings {
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

	return Settings{
		Addr:       addr,
		Scheme:     "http",
		Storage:    storage,
		SessionKey: sessionKey,
	}
}
