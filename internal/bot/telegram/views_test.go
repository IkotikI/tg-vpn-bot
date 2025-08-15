package telegram

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	xui_service "vpn-tg-bot/internal/service/subscription/xui"
	"vpn-tg-bot/internal/service/vpnserver"
	"vpn-tg-bot/internal/storage/sqlite"

	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/telebot.v4"
)

func TestLocalize(t *testing.T) {
	b := prepareBot()
	loc := b.localizers["ru"]
	fmt.Printf("localizers: %v\n", b.localizers)

	out, err := loc.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "msgSettings",
			Other: `*Settings*`,
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Localized message:", out)

}

func TestLocalizeMarkup(t *testing.T) {
	b := prepareBot()
	loc := b.localizers["en"]

	markup := b.view.markupMain()
	markupControl := b.view.markupMain()
	fmt.Printf("Unlocalized keyboard: \n%+v", printKeyboard(markup))
	// fmt.Printf("Unlocalized keyboard: \n%+v\n", markup.InlineKeyboard)
	b.view.LocalizeMarkup(loc, markup)
	fmt.Printf("Localized keyboard: \n%+v", printKeyboard(markup))
	// fmt.Printf("Localized keyboard: \n%+v\n", markup.InlineKeyboard)

	fails := []telebot.InlineButton{}

	for i, container := range markup.InlineKeyboard {
		for j, item := range container {
			if item.Text == markupControl.InlineKeyboard[i][j].Text {
				fails = append(fails, item)
			}
		}
	}

	if len(fails) > 0 {
		out := "Localization fails for:\n"
		for _, item := range fails {
			out += fmt.Sprintf("%s: %s\n", item.Unique, item.Text)
		}
		t.Fatal(out)
	}

}

func TestView(t *testing.T) {
	b := prepareBot()
	loc := b.localizers["en"]

	msg := `*Settings*`
	msg, err := b.view.localizeMessage(loc, "msgSettings", msg, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Localized message:", msg)

	markup := b.view.markupSettings()
	b.view.LocalizeMarkup(loc, markup)
	// fmt.Printf("%+v", markup.InlineKeyboard)
	fmt.Print(printKeyboard(markup))
}

func printKeyboard(m *telebot.ReplyMarkup) string {
	out := ""
	out += "Inline keyboard: \n"
	for _, container := range m.InlineKeyboard {
		for _, item := range container {
			out += fmt.Sprintf("	%s: %s\n", item.Unique, item.Text)
		}
	}
	out += "Reply keyboard: \n"
	for _, container := range m.ReplyKeyboard {
		for _, item := range container {
			out += fmt.Sprintf("	%s\n", item.Text)
		}
	}
	return out
}

var basePath = "../../../"

func prepareBot() *TelegramBot {
	envPathPtr := flag.String("env", basePath+".env", "path to .env file")
	storagePathPtr := flag.String("storage", basePath+"internal/storage/sqlite/test_data/db.db", "path to storage file")
	flag.Parse()

	envPath := *envPathPtr
	log.Printf("env path: %s", envPath)

	storagePath := *storagePathPtr
	log.Printf("storage path: %s", storagePath)

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
	subsService := xui_service.NewXUIService(token, storage, storage)

	settings := Settings{
		Settings: telebot.Settings{
			Token:   token,
			Poller:  &telebot.LongPoller{Timeout: 10},
			Verbose: true,
		},
	}

	bot, err := New(settings, storage, subsService, vpnSeverManager, basePath+"internal/bot/telegram/i18n")
	if err != nil {
		log.Fatalf("[ERR]: Can't create bot instance: %v", err)
	}

	return bot
}
