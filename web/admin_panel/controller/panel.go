package controller

import (
	"context"
	"log"
	"net/http"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/web/admin_panel/middleware"
	"vpn-tg-bot/web/admin_panel/service"
	"vpn-tg-bot/web/admin_panel/views"
	"vpn-tg-bot/web/admin_panel/views/templates"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

type PanelController struct {
	BaseController

	storage service.StorageService
	// service?
	// cookieStore *sessions.CookieStore
	// mux         http.Handler
}

type State struct {
	Request *http.Request
	Page    Page
	// Menu        Menu
	// MainContent templ.Component
}

type Page struct {
	Title string
	// Template func() templ.Component
}

func NewPanelController(r *mux.Router, storage storage.Storage) *PanelController {

	storage_service, err := service.NewStorageService(storage)
	if err != nil {
		log.Fatalf("[ERR} Can't start storage service: %v", err)
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

func (p *PanelController) registerRoutes(r *mux.Router) {
	// Admin Panel Router
	apiRouter := r.PathPrefix("/").Subrouter()
	apiRouter.Use(middleware.LoggingMiddleware)
	apiRouter.HandleFunc("/", p.defaultView).Methods("GET")
	// r.HandleFunc("GET /", p.defaultView)
}

func (p *PanelController) defaultView(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	users, err := p.storage.GetUsersWithSubscription(ctx, nil)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err.Error())
		log.Printf("[ERR] Can't get users: %v", err)
		return
	}

	UIviews := map[string]templ.Component{
		"main": templates.UsersTable(users),
		"menu": templates.MainNavigation(),
	}

	index := views.Index(views.UI(r, UIviews))
	index.Render(ctx, w)
}

// func (p *AdminPanel) appUI(state AdminPanelState) templ.Component {

// }
