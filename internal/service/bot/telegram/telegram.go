package telegram

import (
	"time"
	"vpn-tg-bot/internal/service/subscription"
	"vpn-tg-bot/internal/service/vpnserver"
	"vpn-tg-bot/internal/storage"

	telebot "gopkg.in/telebot.v4"
)

type TelegramBot struct {
	MaxResponseTime          time.Duration
	DemoSubscriptionDuration time.Duration

	tg               *telebot.Bot
	storage          storage.Storage
	subscriptions    subscription.VPN_API
	vpnserverManager *vpnserver.VPNServerManager

	view *ViewBuilder
}

type Settings struct {
	telebot.Settings
}

func New(settings Settings, storage storage.Storage, subscriptions subscription.VPN_API, vpnserverManager *vpnserver.VPNServerManager) (*TelegramBot, error) {
	bot, err := telebot.NewBot(settings.Settings)
	if err != nil {
		return nil, err
	}
	tgbot := &TelegramBot{
		MaxResponseTime:          2000,
		DemoSubscriptionDuration: 2 * 24 * time.Hour,
		tg:                       bot,
		storage:                  storage,
		subscriptions:            subscriptions,
		vpnserverManager:         vpnserverManager,
	}
	tgbot.registerHandlers()
	tgbot.view = NewViewBulder()
	return tgbot, nil
}

func (b *TelegramBot) Run() {
	b.tg.Start()
}

func (b *TelegramBot) Stop() {
	b.tg.Stop()
}
