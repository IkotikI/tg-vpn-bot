package telegram

import (
	"bytes"
	"errors"
	"text/template"

	"gopkg.in/telebot.v4"
)

var ErrKeyboardNotFound = errors.New("keyboard not found")

// Telegram supports 2 types of keyboard:
// 1. Reply keyboard: postioned under the chat text input
// 2. Inline keyboard: attached to a message
var (
	btnProfile   = m.Text("Profile")
	btnSettings  = m.Text("Settings")
	btnBuySub    = m.Text("New VPN")
	btnMySubs    = m.Text("My VPNs")
	btnUpBalance = m.Text("Up balance")
	btnSettings2 = m.Text("Settings")

	// Default menu markup
	m = &telebot.ReplyMarkup{}
)

type ViewBuilder struct {
	markups map[string]*telebot.ReplyMarkup
}

func NewViewBulder() *ViewBuilder {
	b := &ViewBuilder{}
	b.createMarkups()
	return b
}

func (b *ViewBuilder) MustKeyboard(name string) *telebot.ReplyMarkup {
	k, ok := b.markups[name]
	if !ok {
		panic("keyboard not found")
	}
	return k
}

func (b *ViewBuilder) Keyboard(name string) (*telebot.ReplyMarkup, error) {
	k, ok := b.markups[name]
	if !ok {
		return nil, ErrKeyboardNotFound
	}
	return k, nil
}

func ref[T any](v T) *T {
	return &v
}

func (b *ViewBuilder) createMarkups() {
	var m *telebot.ReplyMarkup

	// Main menu keyboard
	m = &telebot.ReplyMarkup{}
	m.Inline(
		m.Row(btnBuySub, btnMySubs),
		m.Row(btnUpBalance, btnProfile),
		m.Row(btnSettings),
	)
	b.markups["main"] = m

	// Reply keyboard
	m = &telebot.ReplyMarkup{ResizeKeyboard: true}
	m.Reply(
		m.Row(btnProfile, btnSettings),
	)
	b.markups["reply"] = m

	// Newcomer keyboard
	m = &telebot.ReplyMarkup{}
	m.Reply(
		m.Row(btnProfile, btnSettings),
	)
	b.markups["newcomer"] = m

	// Profile keyboard
	m = &telebot.ReplyMarkup{ResizeKeyboard: true}
	m.Reply(
		m.Row(btnProfile, btnSettings),
	)
	b.markups["profile"] = m

	// Settings keyboard
	m = &telebot.ReplyMarkup{ResizeKeyboard: true}
	m.Reply(
		m.Row(btnProfile, btnSettings),
	)
	b.markups["settings"] = m

	// Up balance keyboard
	m = &telebot.ReplyMarkup{ResizeKeyboard: true}
	m.Reply(
		m.Row(btnProfile, btnSettings),
	)
	b.markups["up_balance"] = m

}

func (c *ViewBuilder) priceTable() string {
	return `💰 Наши цены после истечения пробной версии:
├ 1 мес: $5
├ 6 мес: $27 (-10%)
├ 1 год: $48.7 (-20%)
├ 3 года: $109.5 (-40%)`
}

func (c *ViewBuilder) PrepareMessage(msg string, args interface{}) (string, error) {
	temp, err := template.New("msg").Parse(msg)
	if err != nil {
		return "", err
	}

	b := &bytes.Buffer{}
	err = temp.Execute(b, args)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func (b *ViewBuilder) newcomerView(c telebot.Context, args map[string]interface{}) (msg string, opts []interface{}, err error) {
	msg = `Добро пожаловать в @VA_VPN_TG_Dev_bot

Ваш VPN уже готов к работе и будет доступен БЕСПЛАТНО три дня!

Установите приложение для вашей OS 👇

🍏 iOS (https://apps.apple.com/ru/app/v2raytun/id6476628951)

🤖 Android (https://github.com/hiddify/HiddifyNG/releases/download/v6.0.4/HiddifyNG.apk)

🖥️ Windows (https://github.com/hiddify/hiddify-next/releases/download/v1.5.2/Hiddify-Windows-Setup-x64.exe)

🍏 MacOS (https://apps.apple.com/us/app/v2box-v2ray-client/id6446814690)

Подключите VPN ключ в приложение 👇

(Нажмите на текст ниже, чтобы скопировать):
`
	msg += "`{{ .SubLink }}`\n\n"
	msg += "-----------------------------"
	msg += b.priceTable()

	m := b.MustKeyboard("newcomer")

	return msg, []interface{}{m, telebot.ModeMarkdownV2}, nil
}

func (b *ViewBuilder) mainView(c telebot.Context, args map[string]interface{}) (msg string, opts []interface{}, err error) {
	msg = `🔥 Наши серверы не имеют ограничений по скорости и трафику, VPN работает на всех устройствах, YouTube в 4К – без задержек!

🔥 Максимальная анонимность и безопасность, которую не даст ни один VPN сервис в мире.

✅ Наш канал: @VA_VPN_TG_Dev
	`

	m := b.MustKeyboard("main")

	return msg, []interface{}{m, telebot.ModeMarkdownV2}, nil
}

func (b *ViewBuilder) viewProfile(c telebot.Context, args map[string]interface{}) (msg string, opts []interface{}, err error) {
	t := `*Profile*

ID: {{ .ID }}
Balance: {{ .Balance }}
You have {{ .SubscriptionCount }} active subscriptions
`
	msg, err = b.PrepareMessage(t, args)
	if err != nil {
		return "", nil, err
	}

	m := b.MustKeyboard("profile")
	// m = defaultInlineKeyboard()

	// m.Reply(
	// 	m.Row(btnProfile, btnSettings),
	// )

	return msg, []interface{}{m, telebot.ModeMarkdownV2}, nil
}

func (b *ViewBuilder) viewMySubs(c telebot.Context, args map[string]interface{}) (msg string, opts []interface{}, err error) {
	// subs, ok := args["subscriptions"].(*[]storage.Subscription)
	// if !ok {
	// 	return "", nil, fmt.Errorf("can't exctract `subscriptions: *[]storage.Subscription` from `args map[string]interface{}`")
	// }
	// msg = `*My Subscriptions*`
	// m := telebot.ReplyMarkup{}
	// rows := telebot.Row()
	// for _, s := range *subs {
	// 	rows = append(rows, m.Data(s.Server))
	// }
	return msg, []interface{}{m, telebot.ModeMarkdownV2}, nil
}
