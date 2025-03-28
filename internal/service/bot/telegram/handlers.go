package telegram

import (
	"context"
	"fmt"
	"log"
	"time"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/e"

	"gopkg.in/telebot.v4"
)

func (b *TelegramBot) registerHandlers() {
	b.tg.Handle("/start", b.handleStart)
	b.tg.Handle(telebot.OnText, b.handleText)
	// Profile
	// b.tg.Handle(&btnProfile, b.handleProfile)

	b.tg.Handle(&btnProfile, b.handleProfile)
	// b.tg.Handle(&btnMySubs, b.handleMySubs)
	// b.tg.Handle(&btnBuySub, b.handleBuySub)
	// b.tg.Handle(&btnSettings, b.handleSettings)

}

func (b *TelegramBot) handleStart(c telebot.Context) error {
	log.Printf("Received /start from %s", c.Sender().Username)
	ctx, cancel := context.WithTimeout(context.Background(), b.MaxResponseTime*time.Millisecond)
	defer cancel()

	tgID := storage.TelegramID(c.Sender().ID)
	_, err := b.storage.GetUserByTelegramID(ctx, tgID)
	if err == storage.ErrNoSuchUser {
		// If user is not registered, then register him and give demo subscription
		c.Send("Registering...")
		fmt.Printf("Registring new user %s\n", c.Sender().Username)

		demoServer, err := b.vpnserverManager.GetDemoServer(ctx)
		if err != nil {
			return e.Wrap("can't get demo server", err)
		}
		user := &storage.User{
			TelegramID:   tgID,
			TelegramName: c.Sender().Username,
		}
		userID, err := b.storage.SaveUser(ctx, user)
		if err != nil {
			return e.Wrap("can't save user", err)
		}
		sub := &storage.Subscription{
			UserID:                userID,
			ServerID:              demoServer.ID,
			SubscriptionStatus:    storage.SubscriptionStatusActive,
			SubscriptionExpiredAt: time.Now().Add(b.DemoSubscriptionDuration),
		}
		err = b.subscriptions.UpdateSubscription(ctx, sub)
		if err != nil {
			return e.Wrap("can't update subscription", err)
		}
		c.Edit("Registered!", b.view.MustKeyboard("main"), telebot.ModeMarkdown)
	} else if err != nil {
		return err
	} else {
		return c.Send("Hello!\nChose, what do you want to do?", b.view.MustKeyboard("main"), telebot.ModeMarkdown)
	}

	return nil
}

func (b *TelegramBot) handleText(c telebot.Context) error {
	log.Printf("Received text from %s: %s", c.Sender().Username, c.Text())
	return c.Send("Wtf u write?!", b.view.MustKeyboard("reply"), telebot.ModeMarkdown)
}

func (b *TelegramBot) handleProfile(c telebot.Context) error {
	log.Printf("Received text from %s: %s", c.Sender().Username, c.Text())

	ctx, cancel := context.WithTimeout(context.Background(), b.MaxResponseTime*time.Millisecond)
	defer cancel()

	args := map[string]interface{}{}

	queryArgs := &storage.QueryArgs{
		From: storage.TableSubscriptions,
		Where: []storage.Where{
			{
				Column:   "user_id",
				Operator: "=",
				Value:    c.Sender().ID,
			},
		},
	}
	n, err := b.storage.Count(ctx, queryArgs)
	if err == nil {
		args["SubscriptionCount"] = n
	} else {
		log.Printf("[ERR]: Telegram Bot: handleProfile: %v", err)
		args["SubscriptionCount"] = "err"
	}

	args["Balance"] = "<Hardcoded nill>"

	msg, opts, err := b.view.viewProfile(c, args)
	if err != nil {
		return err
	}

	return c.Send(msg, telebot.ModeMarkdown, opts)
}

// func (b *TelegramBot) handleMySubs(c telebot.Context) error {
// 	log.Printf("Received text from %s: %s", c.Sender().Username, c.Text())
// 	ctx, cancel := context.WithTimeout(context.Background(), b.MaxResponseTime*time.Millisecond)
// 	defer cancel()

// 	user, err := b.storage.GetUserByTelegramID(ctx, storage.TelegramID(c.Sender().ID))
// 	if err != nil {
// 		return err
// 	}
// 	subs, err := b.storage.GetSubscriptionsByUserID(ctx, user.ID)
// 	if err != nil {
// 		return err
// 	}
// 	servers, err := b.storage.GetServers(ctx, &storage.QueryArgs{Where: []storage.Where{{Column: "id", Operator: "in", Value: subs}}})
// 	if err != nil {
// 		return err
// 	}
// 	msg, opts, err := b.view.viewMySubs(c, map[string]interface{}{"subscriptions": &subs, "servers": &servers})
// 	if err != nil {
// 		return err
// 	}
// 	return c.Edit(msg, opts)
// }
