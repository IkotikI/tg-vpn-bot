package telegram

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/template"
	"vpn-tg-bot/internal/storage"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/telebot.v4"
)

var ErrKeyboardNotFound = errors.New("keyboard not found")

// Telegram supports 2 types of keyboard:
// 1. Reply keyboard: postioned under the chat text input
// 2. Inline keyboard: attached to a message
var (
	replyBtnMainMenu = m.Text("Main Menu")
	btnMainMenu      = m.Data("Main Menu", "MainMenu")
	btnProfile       = m.Data("Profile", "Profile")
	btnSettings      = m.Data("Settings", "Settings")
	btnBuySub        = m.Data("New VPN", "NewSub")
	btnMySubs        = m.Data("My VPNs", "MySubs")
	btnUpBalance     = m.Data("Up balance", "UpBalance")
	btnLanguages     = m.Data("Languages", "Languages")

	nameBtnLanguage  = "language"
	nameBtnBuySub    = "buySub"
	nameBtnMySub     = "mySub"
	nameBtnExtendSub = "ExtendSub"

	msgServerError = m.Data("Sorry, something went wrong. Please try this feauture later or contact support.", "msgServerError")
	// Default menu markup
	m = &telebot.ReplyMarkup{}
)

// TODO: class is redundand
type ViewBuilder struct {
	translationsPath string
}

func NewViewBulder(translationsPath string) *ViewBuilder {
	b := &ViewBuilder{translationsPath: translationsPath}
	return b
}

func (b *ViewBuilder) LocalizeMarkup(l *i18n.Localizer, m *telebot.ReplyMarkup) *telebot.ReplyMarkup {
	for i, row := range m.InlineKeyboard {
		for j, btn := range row {
			if btn.Unique == "" {
				btn.Unique = b.LocalizationIDFromText(btn.Text)
			}
			newText, err := l.Localize(&i18n.LocalizeConfig{
				MessageID: btn.Unique,
			})
			// fmt.Printf("Localized button: %s: from %s to %s \n", btn.Unique, btn.Text, newText)
			if err != nil {
				fmt.Printf("[ERR] Telegram Bot: can't localize inline button: %s\n", err)
			} else {
				m.InlineKeyboard[i][j].Text = newText
			}

		}
	}

	// Reply button hasn't unique field, thus we convers text from
	// "Main Menu" to "btnMainMenu"
	for i, row := range m.ReplyKeyboard {
		for j, btn := range row {
			btn.Text = b.LocalizationIDFromText(btn.Text)
			newText, err := l.Localize(&i18n.LocalizeConfig{
				MessageID: btn.Text,
			})
			// fmt.Printf("Localized button: %s: from %s to %s \n", btn.Unique, btn.Text, newText)
			if err != nil {
				fmt.Printf("[ERR] Telegram Bot: can't localize reply button: %s\n", err)
			} else {
				m.ReplyKeyboard[i][j].Text = newText
			}

		}
	}

	// TODO: Test Reply keyboard
	return m
}

func (b *ViewBuilder) LocalizationIDFromText(text string) string {

	text = strings.ReplaceAll(text, " ", "")
	// text = "btn" + text
	return text
}

// This using Reply makup and should be used in separate requiest!
func (b *ViewBuilder) markupMainMenuReplyBtn() *telebot.ReplyMarkup {
	m = &telebot.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: false}
	m.Reply(
		m.Row(replyBtnMainMenu),
	)
	return m
}

// This using Reply makup and should be used in separate requiest!
func (b *ViewBuilder) markupMainMenuInlineBtn() *telebot.ReplyMarkup {
	m = &telebot.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: false}
	m.Inline(
		m.Row(btnMainMenu),
	)
	return m
}

// Up balance keyboard
func (b *ViewBuilder) markupUpBalance() *telebot.ReplyMarkup {
	m = &telebot.ReplyMarkup{}
	m.Inline(
		m.Row(btnProfile, btnSettings),
	)
	return m
}

var priceMap map[string]string = map[string]string{
	"1m": "$5",
	"6m": "$27 (-10%)",
	"1y": "$48.7 (-20%)",
	"3y": "$109.5 (-40%)",
}

func periodFromPriceMapKey(key string) (count int, units string) {
	units = key[len(key)-1:]
	count, _ = strconv.Atoi(key[:len(key)-1])
	return count, units
}

func unitsToMsgID(units string) string {
	switch units {
	case "m":
		return "nMonths"
	case "y":
		return "nYears"
	}
	return "undefinedUnits"
}

func (b *ViewBuilder) priceTable(priceMap map[string]string) string {
	msg := `💰 Наши цены после истечения пробной версии:
	`
	for k, v := range priceMap {
		msg += fmt.Sprintf("├ %s: %s\n", k, v)
	}
	return msg
}

func (b *ViewBuilder) PrepareMessage(msg string, args interface{}) (string, error) {
	temp, err := template.New("msg").Parse(msg)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = temp.Execute(buf, args)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Newcomer view is preseted, when a new user write /start
// args: map[SubLink]<link_content>
func (b *ViewBuilder) viewNewcomer(loc *i18n.Localizer, args map[string]interface{}) (msg string, opts []interface{}, err error) {

	msg, _ = b.localizeMessage(loc, "msgNewcomer", msg, args)

	markup := b.markupNewcomer()
	b.LocalizeMarkup(loc, markup)

	return msg, []interface{}{markup, telebot.ModeHTML}, nil
}

// Newcomer keyboard
func (b *ViewBuilder) markupNewcomer() *telebot.ReplyMarkup {
	m = &telebot.ReplyMarkup{}
	m.Inline(
		m.Row(btnLanguages, btnLanguages),
		m.Row(btnProfile, btnSettings),
	)
	return m
}

func (b *ViewBuilder) viewMain(loc *i18n.Localizer, args map[string]interface{}) (msg string, opts []interface{}, err error) {
	msg = `Our servers have no speed and traffic limits, VPN works on all devices, YouTube in 4K - without delays!`
	msg, _ = b.localizeMessage(loc, "msgMain", msg, args)
	// 	msg = `🔥 Наши серверы не имеют ограничений по скорости и трафику, VPN работает на всех устройствах, YouTube в 4К – без задержек!

	// 🔥 Максимальная анонимность и безопасность, которую не даст ни один VPN сервис в мире.

	// ✅ Наш канал: @VA_VPN_TG_Dev
	// 	`

	markup := b.markupMain()
	b.LocalizeMarkup(loc, markup)

	return msg, []interface{}{markup, telebot.ModeHTML}, nil
}

// Main menu keyboard
func (b *ViewBuilder) markupMain() *telebot.ReplyMarkup {
	m = &telebot.ReplyMarkup{}

	m.Inline(
		m.Row(btnBuySub, btnMySubs),
		m.Row(btnUpBalance, btnProfile),
		m.Row(btnSettings),
	)
	return m
}

// User's profile.
// Require fields:
// - ID - user's Telegram ID
// - Balance - user's balance
// - SubsCount - number of user's active subscriptions
func (b *ViewBuilder) viewProfile(loc *i18n.Localizer, args map[string]interface{}) (msg string, opts []interface{}, err error) {
	msg = `**Profile**

ID: {{ .ID }}
Balance: {{ .Balance }}
Active subscriptions: {{ .SubsCount }} 
`
	msg, _ = b.localizeMessage(loc, "msgProfile", msg, args)

	// msg, err = b.PrepareMessage(t, args)
	// if err != nil {
	// 	return "", nil, err
	// }

	markup := b.markupProfile()
	b.LocalizeMarkup(loc, markup)

	return msg, []interface{}{markup, telebot.ModeHTML}, nil
}

// Profile keyboard
func (b *ViewBuilder) markupProfile() *telebot.ReplyMarkup {
	m = &telebot.ReplyMarkup{}
	m.Inline(
		m.Row(btnSettings, btnUpBalance),
	)
	return m
}

// My Subsciption shows list of user's active VPN subscriptions
func (b *ViewBuilder) viewMySubs(loc *i18n.Localizer, subs *[]storage.SubscriptionWithServer) (msg string, opts []interface{}, err error) {
	msg = `*My Subscriptions*`
	msg, _ = b.localizeMessage(loc, "msgMySubs", msg, nil)

	var markup *telebot.ReplyMarkup
	l := len(*subs)
	if l == 0 {
		msg += "\nYou have no active subscriptions\\."
		markup = &telebot.ReplyMarkup{}
		markup.Inline(markup.Row(btnBuySub))
	} else {
		msg += fmt.Sprintf("\nYou have %d active subscriptions\\.", l)
		markup = b.markupMySubscriptions(subs)
	}
	b.LocalizeMarkup(loc, markup)

	return msg, []interface{}{markup, telebot.ModeHTML}, nil
}

func (b *ViewBuilder) markupMySubscriptions(subs *[]storage.SubscriptionWithServer) *telebot.ReplyMarkup {
	m := telebot.ReplyMarkup{}
	rows := []telebot.Row{}
	for _, s := range *subs {
		rows = append(rows, m.Row(m.Data(s.Name, nameBtnMySub, s.ID.String())))
	}
	m.Inline(rows...)
	return &m
}

// Buy Subscription view
// args: map[ Period in format number+unit(string) ] Price (string)
func (b *ViewBuilder) viewBuySub(loc *i18n.Localizer, args map[string]string) (msg string, opts []interface{}, err error) {
	msg = "*Buy VPN*"
	msg, _ = b.localizeMessage(loc, "msgBuySub", msg, args)
	msg += "\n"

	var msgID string
	for k, v := range args {
		count, units := periodFromPriceMapKey(k)
		msgID = unitsToMsgID(units)
		// log.Printf("viewBuySub: localization args for new vpn period %v, %v, %v\n", unique, k, struct{ Count int }{Count: count})
		line, err := b.localizeMessage(loc, msgID, k, struct{ Count int }{Count: count})
		if err != nil {
			return "", nil, err
		}
		msg += fmt.Sprintf("├ %s: %s\n", line, v)
	}

	markup := b.markupBuySubscription(loc, args)
	// b.LocalizeMarkup(loc, markup)

	return msg, []interface{}{markup, telebot.ModeHTML}, nil
}

func (b *ViewBuilder) markupBuySubscription(loc *i18n.Localizer, priceMap map[string]string) *telebot.ReplyMarkup {
	m := telebot.ReplyMarkup{}
	rows := []telebot.Row{}
	for k := range priceMap {
		count, units := periodFromPriceMapKey(k)
		msgID := unitsToMsgID(units)
		// log.Printf("viewBuySub: localization args for new vpn period %v, %v, %v\n", unique, k, struct{ Count int }{Count: count})
		line, err := b.localizeMessage(loc, msgID, k, struct{ Count int }{Count: count})
		if err != nil {
			log.Printf("[ERR] Telegram Bot: can't localize inline button: %s\n", err)
		}

		rows = append(rows, m.Row(m.Data(line, nameBtnBuySub, k)))
	}
	m.Inline(rows...)
	return &m
}

// Settings view
func (b *ViewBuilder) viewSettings(loc *i18n.Localizer, args map[string]interface{}) (msg string, opts []interface{}, err error) {
	msg = `*Settings*`
	msg, _ = b.localizeMessage(loc, "msgSettings", msg, args)

	markup := b.markupSettings()
	b.LocalizeMarkup(loc, markup)

	return msg, []interface{}{markup, telebot.ModeHTML}, nil
}

// Settings keyboard
func (b *ViewBuilder) markupSettings() *telebot.ReplyMarkup {
	m = &telebot.ReplyMarkup{}
	m.Inline(
		m.Row(btnProfile, btnLanguages),
	)
	return m
}

func (b *ViewBuilder) viewLanguages(loc *i18n.Localizer, args map[string]interface{}) (msg string, opts []interface{}, err error) {
	msg, _ = b.localizeMessage(loc, "msgLanguages", msg, args)

	markup := b.markupLanguages(args)
	return msg, []interface{}{markup, telebot.ModeHTML}, nil
}

func (b *ViewBuilder) markupLanguages(args map[string]interface{}) *telebot.ReplyMarkup {
	m = &telebot.ReplyMarkup{}

	l := len(args)
	var i int
	var prev string
	var rows []telebot.Row
	for cur := range args {
		i++
		if i%2 == 0 {
			rows = append(rows, m.Row(
				m.Data(cur, nameBtnLanguage, cur),
				m.Data(prev, nameBtnLanguage, prev),
			))
		} else {
			if i == l {
				rows = append(rows, m.Row(
					m.Data(cur, nameBtnLanguage, cur),
				))
				break
			}
			prev = cur
		}
	}
	m.Inline(rows...)
	return m
}

func (b *ViewBuilder) viewSub(loc *i18n.Localizer, sub *storage.SubscriptionWithUserAndServer, link string) (msg string, opts []interface{}, err error) {
	msg = `*VPN Subscription*`

	var args = struct {
		Server  string
		Expired string
		Link    string
	}{
		Server:  sub.VPNServerWithCountry.CountryCode + " " + sub.VPNServer.Name,
		Expired: sub.SubscriptionExpiredAt.Format("02.01.2006 15:04"),
		Link:    link,
	}
	msg, _ = b.localizeMessage(loc, "msgSub", msg, args)

	markup := b.markupSub(sub)
	b.LocalizeMarkup(loc, markup)

	return msg, []interface{}{markup, telebot.ModeHTML}, nil
}

func (b *ViewBuilder) markupSub(sub *storage.SubscriptionWithUserAndServer) *telebot.ReplyMarkup {
	m = &telebot.ReplyMarkup{}
	m.Inline(
		m.Row(m.Data("Extend", nameBtnExtendSub, sub.ServerID.String())),
	)
	return m
}

func (b *ViewBuilder) localizeMessage(loc *i18n.Localizer, msgID string, msg string, args interface{}) (outMsg string, err error) {

	cfg := &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    msgID,
			Other: msg,
		},
	}
	if args != nil {
		cfg.TemplateData = args
	}
	msg, err = loc.Localize(cfg)

	if err != nil {
		log.Println(err)
	}

	return msg, err
}

func EscapeMarkdownV1(s string) string {
	specials := "_*`[\\"
	for _, ch := range specials {
		s = strings.ReplaceAll(s, string(ch), "\\"+string(ch))
	}
	return s
}

func EscapeMarkdownV2(s string) string {
	specials := "_*[]()~`>#+-=|{}.!"
	for _, ch := range specials {
		s = strings.ReplaceAll(s, string(ch), "\\"+string(ch))
	}
	return s
}
