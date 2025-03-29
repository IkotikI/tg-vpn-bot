package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/web/admin_panel/middleware"
	"vpn-tg-bot/web/admin_panel/views"
	"vpn-tg-bot/web/admin_panel/views/templates"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

// Page Timeout in Milliseconds
var PageTimeout = 1000 * time.Millisecond

type PanelController struct {
	BaseController

	storage storage.Storage
	// service?
	// cookieStore *sessions.CookieStore
	// mux         http.Handler
}

func NewPanelController(r *mux.Router, storage storage.Storage) *PanelController {
	p := &PanelController{
		storage: storage,
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
	apiRouter.HandleFunc("/subscriptions", c.subscriptionsView)
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

	queryArgs := storage.ParseQueryArgs(values)
	// queryArgs.ParseDefaultsFrom(storage.DefaultQueryArguments)

	args := queryArgs.ToQueryArgs()

	users, err := c.storage.GetUsers(ctx, args)
	if err != nil {
		writeError(w, ctx, http.StatusInternalServerError, ErrMsgServerInternalError)
		log.Printf("[ERR] Can't get users: %v", err)
		return
	}

	pagination, err := c.storage.MakePagination(ctx, c.storage, "users", queryArgs, args)
	if err != nil {
		log.Printf("[ERR] Error occurred while making pagination:", err)
	}
	paginationLinks := views.PaginationLinks("/users", pagination, 3)

	fmt.Printf("queryArgs: \n%+v\n", queryArgs)
	fmt.Printf("pagination: \n%+v\n", pagination)
	fmt.Printf("paginationLinks: \n%+v\n", paginationLinks)

	switch r.Method {
	case "GET":
		UIviews := map[string]templ.Component{
			"main":       templates.UsersMain(users),
			"header":     templates.DefaultHeader(nil),
			"pagination": templates.DefaultPagination(paginationLinks),
		}

		index := views.Index(views.UI(r, UIviews))
		index.Render(ctx, w)
	case "POST":
		table := templates.UsersTable(users)
		table.Render(ctx, w)
	}

}

func (c *PanelController) serversView(w http.ResponseWriter, r *http.Request) {
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

	queryArgs := storage.ParseQueryArgs(values)
	// queryArgs.ParseDefaultsFrom(storage.DefaultQueryArguments)

	args := queryArgs.ToQueryArgs()

	servers, err := c.storage.GetServersWithCountries(ctx, args)
	if err != nil {
		writeError(w, ctx, http.StatusInternalServerError, ErrMsgServerInternalError)
		log.Printf("[ERR] Can't get users: %v", err)
		return
	}

	pagination, err := c.storage.MakePagination(ctx, c.storage, "servers", queryArgs, args)
	if err != nil {
		log.Printf("[ERR] Error occurred while making pagination:", err)
	}
	paginationLinks := views.PaginationLinks("/servers", pagination, 3)

	fmt.Printf("queryArgs: \n%+v\n", queryArgs)
	fmt.Printf("pagination: \n%+v\n", pagination)
	fmt.Printf("paginationLinks: \n%+v\n", paginationLinks)

	switch r.Method {
	case "GET":
		UIviews := map[string]templ.Component{
			"main":       templates.ServersTable(servers),
			"header":     templates.DefaultHeader(nil),
			"pagination": templates.DefaultPagination(paginationLinks),
		}

		index := views.Index(views.UI(r, UIviews))
		index.Render(ctx, w)
	case "POST":
		table := templates.ServersTable(servers)
		table.Render(ctx, w)
	}

}

func (p *PanelController) subscriptionsView(w http.ResponseWriter, r *http.Request) {
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

	args := storage.ParseSelectQueryArgs(values)

	subs, err := p.storage.GetSubscriptionsWithUsersAndServers(ctx, args)
	if err != nil {
		writeError(w, ctx, http.StatusInternalServerError, ErrMsgServerInternalError)
		log.Printf("[ERR] Can't get subscriptions: %v", err)
		return
	}

	switch r.Method {
	case "GET":
		UIviews := map[string]templ.Component{
			"main":   templates.SubscriptionsTable(subs),
			"header": templates.DefaultHeader(nil),
		}

		index := views.Index(views.UI(r, UIviews))
		index.Render(ctx, w)
	case "POST":
		table := templates.SubscriptionsTable(subs)
		table.Render(ctx, w)
	}

}

// func (p *AdminPanel) appUI(state AdminPanelState) templ.Component {

// }
