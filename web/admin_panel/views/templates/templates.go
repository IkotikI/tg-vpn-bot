package templates

import (
	"vpn-tg-bot/internal/storage"

	"github.com/a-h/templ"
)

func UserLink(id storage.UserID) string {
	return "/user/" + id.String()
}

func ServerLink(id storage.ServerID) string {
	return "/server/" + id.String()
}

// Provide link in templ.SafeURL format.
func UserLinkT(id storage.UserID) templ.SafeURL {
	return templ.URL(UserLink(id))
}

// Provide link in templ.SafeURL format.
func ServerLinkT(id storage.ServerID) templ.SafeURL {
	return templ.URL(ServerLink(id))
}
