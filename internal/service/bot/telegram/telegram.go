package telegram

import (
	"vpn-tg-bot/internal/storage"

	telebot "gopkg.in/telebot.v4"
)

type Processor struct {
	tg      *Bot
	storage storage.Storage
}

func New(bot *Bot, storage storage.Storage) *Processor {
	return &Processor{
		tg:      bot,
		storage: storage,
	}
}

func (p *Processor) Handle(cmd string, f func() error) {
	p.tg.Handle(cmd,
		func(ctx telebot.Context) error {
			return f()
		},
	)
}

type Settings struct {
	telebot.Settings
}

type Bot struct {
	*telebot.Bot
}

func NewBot(s Settings) (*Bot, error) {
	bot, err := telebot.NewBot(s.Settings)
	if err != nil {
		return nil, err
	}
	return &Bot{Bot: bot}, nil
}

func (p *Processor) Send(text string) {
	// p.tg.Bot.Send(text)
}
