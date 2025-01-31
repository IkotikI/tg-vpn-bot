package admin_panel

import (
	"context"
	"log"
	"net/http"
	"vpn-tg-bot/web/admin_panel/templates"
)

func (p *AdminPanel) registerRoutes() http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("/", p.defaultView)

	return m
}

func (p *AdminPanel) defaultView(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	users, err := p.storage.GetAllUsers(ctx)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err.Error())
		log.Printf("[ERR] Can't get users: %v", err)
		return
	}

	component := index(templates.Users(users))
	component.Render(ctx, w)
}

// func (p *AdminPanel) appUI(state AdminPanelState) templ.Component {

// }
