package controller

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/web/admin_panel/middleware"
	"vpn-tg-bot/web/admin_panel/service"
	"vpn-tg-bot/web/admin_panel/views"
	"vpn-tg-bot/web/admin_panel/views/templates"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

const ErrMsgServerInternalError = "Server Internal Error."

// Page Timeout in Milliseconds
var PageTimeout = 1000 * time.Millisecond

type PanelController struct {
	BaseController

	storage service.StorageService
	// service?
	// cookieStore *sessions.CookieStore
	// mux         http.Handler
}

func NewPanelController(r *mux.Router, storage storage.Storage) *PanelController {

	storage_service, err := service.NewStorageService(storage)
	if err != nil {
		log.Fatalf("[ERR] Can't start storage service: %v", err)
	}

	p := &PanelController{
		storage: storage_service,
	}
	p.registerRoutes(r)
	return p
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	return nil
}

func (c *PanelController) registerRoutes(r *mux.Router) {
	// Admin Panel Router
	apiRouter := r.PathPrefix("/").Subrouter()
	apiRouter.Use(middleware.LoggingMiddleware)
	apiRouter.HandleFunc("/", c.defaultView)
	apiRouter.HandleFunc("/users", c.usersView)
	apiRouter.HandleFunc("/servers", c.serversView)
	apiRouter.HandleFunc("/tools", c.defaultView)
	// r.HandleFunc("GET /", p.defaultView)
}

func (c *PanelController) defaultView(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), PageTimeout)
	defer cancel()

	UIviews := map[string]templ.Component{
		"main":   templates.DefaultMain(),
		"header": templates.DefaultHeader(nil),
	}
	index := views.Index(views.UI(r, UIviews))
	index.Render(ctx, w)

}

func (c *PanelController) usersView(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), PageTimeout)
	defer cancel()

	var values url.Values

	switch r.Method {
	case "GET":
		values = r.URL.Query()
	case "POST":
		r.ParseForm()
		values = r.PostForm
	default:
		return
	}

	args := service.ParseSelectQueryArgs(values)

	users, err := c.storage.GetUsers(ctx, args)
	if err != nil {
		writeError(w, ctx, http.StatusInternalServerError, ErrMsgServerInternalError)
		log.Printf("[ERR] Can't get users: %v", err)
		return
	}

	switch r.Method {
	case "GET":
		UIviews := map[string]templ.Component{
			"main":   templates.UsersMain(users),
			"header": templates.DefaultHeader(nil),
		}

		index := views.Index(views.UI(r, UIviews))
		index.Render(ctx, w)
	case "POST":
		table := templates.UsersTable(users)
		table.Render(ctx, w)
	}

}

func (p *PanelController) serversView(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), PageTimeout)
	defer cancel()

	var values url.Values

	switch r.Method {
	case "GET":
		values = r.URL.Query()
	case "POST":
		r.ParseForm()
		values = r.PostForm
	}

	args := service.ParseSelectQueryArgs(values)

	servers, err := p.storage.GetServers(ctx, args)
	if err != nil {
		writeError(w, ctx, http.StatusInternalServerError, ErrMsgServerInternalError)
		log.Printf("[ERR] Can't get users: %v", err)
		return
	}

	switch r.Method {
	case "GET":
		UIviews := map[string]templ.Component{
			"main":   templates.ServersTable(servers),
			"header": templates.DefaultHeader(nil),
		}

		index := views.Index(views.UI(r, UIviews))
		index.Render(ctx, w)
	case "POST":
		table := templates.ServersTable(servers)
		table.Render(ctx, w)
	}

}

// func (p *AdminPanel) appUI(state AdminPanelState) templ.Component {

// }
