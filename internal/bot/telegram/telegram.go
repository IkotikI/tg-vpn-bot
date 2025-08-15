package telegram

import (
	"fmt"
	"log"
	"path/filepath"
	"time"
	"vpn-tg-bot/internal/service/subscription"
	"vpn-tg-bot/internal/service/vpnserver"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/e"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml"
	"golang.org/x/text/language"
	telebot "gopkg.in/telebot.v4"
)

type TelegramBot struct {
	MaxResponseTime          time.Duration
	DemoSubscriptionDuration time.Duration
	LocalizationDir          string

	tg               *telebot.Bot
	storage          storage.Storage
	subscriptions    subscription.VPN_API
	vpnserverManager *vpnserver.VPNServerManager

	view       *ViewBuilder
	bundle     *i18n.Bundle
	localizers map[string]*i18n.Localizer
}

type Settings struct {
	telebot.Settings
}

func New(settings Settings, storage storage.Storage, subscriptions subscription.VPN_API, vpnserverManager *vpnserver.VPNServerManager, i18nDir string) (*TelegramBot, error) {

	tgbot := &TelegramBot{
		MaxResponseTime:          2000,
		DemoSubscriptionDuration: 2 * 24 * time.Hour,
		storage:                  storage,
		subscriptions:            subscriptions,
		vpnserverManager:         vpnserverManager,
		LocalizationDir:          i18nDir,
	}
	// TODO: think about better error handling function passing.
	// Is external user can manage views to override defautl logging?
	if settings.OnError == nil {
		settings.OnError = tgbot.DefaultErrorHandler
	}
	// Creating bot later to pass DefaultHandler as a function.
	bot, err := telebot.NewBot(settings.Settings)
	if err != nil {
		return nil, err
	}
	tgbot.tg = bot
	// Register handlers for bot. Field tgbot.tg should be set before!
	tgbot.registerHandlers()

	err = tgbot.loadMessageFiles([]string{"en", "ru"})
	if err != nil {
		return nil, err
	}
	tgbot.view = NewViewBulder(tgbot.LocalizationDir)
	return tgbot, err
}

func (b *TelegramBot) Run() {
	b.tg.Start()
}

func (b *TelegramBot) Stop() {
	b.tg.Stop()
}

func (b *TelegramBot) loadMessageFiles(langs []string) error {
	b.localizers = make(map[string]*i18n.Localizer)
	b.bundle = i18n.NewBundle(language.English)
	b.bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	for _, lang := range langs {
		path := lang + ".toml"
		path = filepath.Join(b.LocalizationDir, path)
		_, err := b.bundle.LoadMessageFile(path)
		if err != nil {
			return e.Wrap(fmt.Sprintf("can't load message file for language \"%s\" by path \"%s\"", lang, path), err)
		}
		b.localizers[lang] = i18n.NewLocalizer(b.bundle, lang)
		log.Printf("Loaded message file \"%s\"", path)
	}
	return nil
}
