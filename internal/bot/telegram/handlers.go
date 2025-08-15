package telegram

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/pkg/e"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/telebot.v4"
)

func (b *TelegramBot) registerHandlers() {
	b.tg.Handle("/start", b.handleStart)

	// Profile
	// b.tg.Handle(&btnProfile, b.handleProfile)

	// Reply Button
	b.tg.Handle(&btnMainMenu, b.handleMainMenu)
	b.tg.Handle(&replyBtnMainMenu, b.handleMainMenu)

	// Inline Buttons
	b.tg.Handle(&btnProfile, b.handleProfile)
	b.tg.Handle(&btnSettings, b.handleSettings)
	b.tg.Handle(&btnMySubs, b.handleMySubs)
	b.tg.Handle(&btnBuySub, b.handleBuySub)
	b.tg.Handle(&btnSettings, b.handleSettings)
	b.tg.Handle(&btnUpBalance, b.handleUpBalance)
	b.tg.Handle(&btnLanguages, b.handleLanguages)

	b.tg.Handle(&telebot.Btn{Unique: nameBtnBuySub}, b.handleBuySubBtn)
	b.tg.Handle(&telebot.Btn{Unique: nameBtnLanguage}, b.handleSelectLanguage)
	b.tg.Handle(&telebot.Btn{Unique: nameBtnMySub}, b.handleMySub)
	// b.tg.Handle(&telebot.Btn{Unique: nameBtnExtendSub}, b.handleMySub)

	b.tg.Handle(telebot.OnText, b.handleText)
}

func (b *TelegramBot) handleMainMenu(c telebot.Context) error {
	log.Printf("Received text from %s: %s", c.Sender().Username, c.Text())
	msg, args, err := b.view.viewMain(b.localizerByContext(c), nil)
	if err != nil {
		return err
	}
	// if c.Message() != nil || c.Message().Text != "" {
	// 	fmt.Printf("%+v", c.Message())
	err = c.Edit(msg, args...)
	if err == nil {
		return nil
	}
	if err == telebot.ErrBadContext {
		return c.Send(msg, args...)
	}
	return err
	// }

}

func (b *TelegramBot) handleStart(c telebot.Context) error {
	log.Printf("Received /start from %s", c.Sender().Username)
	ctx, cancel := context.WithTimeout(context.Background(), b.MaxResponseTime*time.Millisecond)
	defer cancel()

	loc := b.localizerByContext(c)

	tgID := storage.TelegramID(c.Sender().ID)
	_, err := b.storage.GetUserByTelegramID(ctx, tgID)
	if err == storage.ErrNoSuchUser {
		// If user is not registered, then register him and give demo subscription
		msg, _ := b.view.localizeMessage(loc, "msgRegister", "Registering..", nil)
		c.Send(msg, b.view.LocalizeMarkup(loc, b.view.markupMainMenuReplyBtn()), telebot.ModeMarkdown)
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
		link, err := b.subscriptions.SubscriptionLink(ctx, demoServer.ID, userID)
		if err != nil {
			log.Printf("[ERR] can't get subscription link %v", err)
			link = "Failed to get free subscription link("
		}
		msg, opts, err := b.view.viewNewcomer(loc, map[string]interface{}{"SubLink": link})
		if err != nil {
			return err
		}
		return c.Send(msg, opts...)
	} else if err != nil {
		return err
	} else {
		// TODO: update keyboard on /start
		// c.Send("Hello!\nChose, what do you want to do?", b.view.LocalizeMarkup(loc, b.view.markupMainMenu()), telebot.ModeMarkdown)
		return c.Send("Hello!\nChose, what do you want to do?", b.view.LocalizeMarkup(loc, b.view.markupMain()), telebot.ModeMarkdown)
	}

	return nil
}

func (b *TelegramBot) handleText(c telebot.Context) error {
	log.Printf("Received text from %s: %s", c.Sender().Username, c.Text())

	// TODO: this need to be cached! Othervies heavy calcualtions
	// MB think global text handler
	var inList bool = false
	for _, loc := range b.localizers {
		msg, err := b.view.localizeMessage(loc, b.view.LocalizationIDFromText(btnMainMenu.Text), "", nil)
		if err == nil && msg == c.Text() {
			inList = true
			break
		}
	}

	if inList {
		return b.handleMainMenu(c)
	}

	return c.Send("Can't process this(", b.view.LocalizeMarkup(b.localizerByContext(c), b.view.markupMain()), telebot.ModeMarkdown)
}

func (b *TelegramBot) handleProfile(c telebot.Context) error {
	log.Printf("Received text from %s: %s", c.Sender().Username, c.Text())

	ctx, cancel := context.WithTimeout(context.Background(), b.MaxResponseTime*time.Millisecond)
	defer cancel()

	args := map[string]interface{}{}

	user, err := b.storage.GetUserByTelegramID(ctx, storage.TelegramID(c.Sender().ID))
	if err != nil {
		log.Printf("[ERR]: Telegram Bot: handleProfile: can't get user %v", err)
		return err
	}

	queryArgs := &storage.QueryArgs{
		From: storage.TableSubscriptions,
		Where: []storage.Where{
			{
				Column:   "user_id",
				Operator: "=",
				Value:    user.ID,
			},
		},
	}

	args["ID"] = user.TelegramID
	n, err := b.storage.Count(ctx, queryArgs)
	if err == nil {
		args["SubsCount"] = n
	} else {
		log.Printf("[ERR]: Telegram Bot: handleProfile: can't count subscriptions %v", err)
		args["SubsCount"] = "err"
	}

	args["Balance"] = "\\<Hardcoded nill>\\"

	msg, opts, err := b.view.viewProfile(b.localizerByContext(c), args)
	if err != nil {
		return err
	}

	return c.Edit(msg, opts...)
}

func (b *TelegramBot) handleMySubs(c telebot.Context) error {
	log.Printf("Received text from %s: %s", c.Sender().Username, c.Text())
	ctx, cancel := context.WithTimeout(context.Background(), b.MaxResponseTime*time.Millisecond)
	defer cancel()

	user, err := b.storage.GetUserByTelegramID(ctx, storage.TelegramID(c.Sender().ID))
	if err != nil {
		return err
	}
	subs, err := b.storage.GetSubscriptionsWithServersByUserID(ctx, user.ID, nil)
	if err != nil {
		return err
	}
	log.Printf("[INFO]: Telegram Bot: handleMySubs: %v", *subs)
	msg, opts, err := b.view.viewMySubs(b.localizerByContext(c), subs)
	if err != nil {
		return err
	}
	return c.Edit(msg, opts...)
}

func (b *TelegramBot) handleBuySub(c telebot.Context) error {
	log.Printf("Received text from %s: %s", c.Sender().Username, c.Text())

	msg, opts, err := b.view.viewBuySub(b.localizerByContext(c), priceMap)
	if err != nil {
		return err
	}
	return c.Edit(msg, opts...)
}

func (b *TelegramBot) handleSettings(c telebot.Context) error {
	log.Printf("Received text from %s: %s", c.Sender().Username, c.Text())

	msg, opts, err := b.view.viewSettings(b.localizerByContext(c), nil)
	if err != nil {
		return err
	}
	return c.Edit(msg, opts...)
}

func (b *TelegramBot) handleUpBalance(c telebot.Context) error {
	log.Printf("Received text from %s: %s", c.Sender().Username, c.Text())
	return c.Edit("Balance not implemented yet", b.view.LocalizeMarkup(b.localizerByContext(c), b.view.markupMain()), telebot.ModeMarkdown)
}

func (b *TelegramBot) handleLanguages(c telebot.Context) error {
	log.Printf("Received text from %s: %s", c.Sender().Username, c.Text())
	args := map[string]interface{}{}
	for lang := range b.localizers {
		args[lang] = nil
	}
	msg, opts, err := b.view.viewLanguages(b.localizerByContext(c), args)
	if err != nil {
		return err
	}
	return c.Edit(msg, opts...)
}

func (b *TelegramBot) handleSelectLanguage(c telebot.Context) error {
	log.Printf("Received text from %s: %s", c.Sender().Username, c.Text())
	return c.Edit("Select Language not implemented yet", b.view.LocalizeMarkup(b.localizerByContext(c), b.view.markupMain()), telebot.ModeMarkdown)
}

func (b *TelegramBot) handleBuySubBtn(c telebot.Context) error {
	log.Printf("Received text from %s: %s", c.Sender().Username, c.Text())
	return c.Edit("This button should provide you to payment page.\nBuy Subscription page not implemented yet", b.view.LocalizeMarkup(b.localizerByContext(c), b.view.markupMain()), telebot.ModeMarkdown)
}

func (b *TelegramBot) handleMySub(c telebot.Context) error {
	log.Printf("Received text from %s: %s", c.Sender().Username, c.Text())
	ctx, cancel := context.WithTimeout(context.Background(), b.MaxResponseTime*time.Millisecond)
	defer cancel()

	data := c.Data()
	if data == "" {
		return fmt.Errorf("empty data for button")
	}
	serverID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return e.Wrap(fmt.Sprintf("can't parse serverID, value: %s", data), err)
	}

	user, err := b.storage.GetUserByTelegramID(ctx, storage.TelegramID(c.Sender().ID))
	if err != nil {
		return err
	}

	sub, err := b.storage.GetSubscriptionWithUserAndServerByIDs(ctx, user.ID, storage.ServerID(serverID))
	if err != nil {
		return err
	}

	msg, opts, err := b.view.viewSub(b.localizerByContext(c), sub, "")
	if err != nil {
		return err
	}
	return c.Edit(msg, opts...)
}

func (b *TelegramBot) DefaultErrorHandler(err error, c telebot.Context) {
	log.Printf("[ERR]: Telegram Bot_: %v", err)

	loc := b.localizerByContext(c)
	msg, _ := b.view.localizeMessage(loc, "msgServerError", "Sorry, something went wrong.", nil)
	err = c.Edit(msg, b.view.LocalizeMarkup(loc, b.view.markupMainMenuInlineBtn()), telebot.ModeMarkdown)
	if err == nil {
		return
	}
	log.Printf("[INFO]: Telegram Bot: DefaultErrorHandler: %v", err)
	err = c.Send(msg, b.view.LocalizeMarkup(loc, b.view.markupMainMenuInlineBtn()), telebot.ModeMarkdown)
	if err != nil {
		log.Printf("[ERR]: Telegram Bot: DefaultErrorHandler: %v", err)
	}
}

func (b *TelegramBot) localizerByContext(c telebot.Context) *i18n.Localizer {

	// TOOD: Actually this logic must include steps to extract User's language setup.

	langCode := c.Sender().LanguageCode
	// DEBUG
	// langCode = "ru"
	// DEBUG
	loc, ok := b.localizers[langCode]
	if !ok {
		return b.localizers["en"]
	}
	return loc
}
