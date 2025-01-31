package admin_panel

import (
	"log"
	"net/http"
	"vpn-tg-bot/internal/storage"
)

type AdminPanel struct {
	Addr string

	storage storage.Storage
	mux     http.Handler
}

// type AdminPanelState struct {
// 	Tabs []AdminPanelTab
// 	MainContent templ.Component
// }

func New(addr string, storage storage.Storage) *AdminPanel {
	return &AdminPanel{
		Addr:    addr,
		storage: storage,
	}
}

func (p *AdminPanel) Run() (err error) {

	p.mux = p.registerRoutes()

	log.Printf("Admin panel is running on %s", p.Addr)
	return http.ListenAndServe(p.Addr, p.mux)
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	return nil
}
