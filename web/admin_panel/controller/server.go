package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/web/admin_panel/entity"
	"vpn-tg-bot/web/admin_panel/service"
	"vpn-tg-bot/web/admin_panel/views"
	"vpn-tg-bot/web/admin_panel/views/templates"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

const (
	ErrMsgServerNotFound   = "Server not found."
	ErrMsgServerIDMismatch = "ServerID mismatch."
)

type ServerController struct {
	BaseController

	storage service.StorageService
	// service?
	// cookieStore *sessions.CookieStore
	// mux         http.Handler
}

func NewServerController(r *mux.Router, storage storage.Storage) *ServerController {

	storage_service, err := service.NewStorageService(storage)
	if err != nil {
		log.Fatalf("[ERR} Can't start storage service: %v", err)
	}

	c := &ServerController{
		storage: storage_service,
	}
	c.registerRoutes(r)
	return c
}

func (c *ServerController) registerRoutes(r *mux.Router) {

	serverRouter := r.PathPrefix("/").Subrouter()
	serverRouter.HandleFunc("/server/{id:[0-9]+}", c.serverView).Methods("GET")
	serverRouter.HandleFunc("/server/{id:[0-9]+}", c.serverUpdate).Methods("PUT", "PATCH")
}

func (c *ServerController) serverView(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), PageTimeout)
	defer cancel()

	vars := mux.Vars(r)
	id, err := getInt64FromVars[storage.ServerID](vars, "id")
	if err != nil {
		c.writeErrorNotFound(w, ctx)
		return
	}

	server, err := c.storage.GetEntityServerByID(ctx, id)
	if err == storage.ErrNoSuchServer {
		c.writeErrorNotFound(w, ctx)
		return
	} else if err != nil {
		log.Printf("[ERR] ServerController: ServerView: GetServer %v", err)
		c.writeErrorServerInternal(w, ctx)
		return
	}

	subs, err := c.storage.GetSubscriptionsWithUsersByServerID(ctx, id, nil)
	if err != nil {
		log.Printf("[ERR] ServerController: ServerView: GetSubscriptions %v", err)
		c.writeErrorServerInternal(w, ctx)
		return
	}

	countries, err := c.storage.GetCountries(ctx, nil)
	if err != nil {
		log.Printf("[ERR] ServerController: ServerView: GetCountries %v", err)
		c.writeErrorServerInternal(w, ctx)
		return
	}

	UI := map[string]templ.Component{
		"main":   templates.ServerMain(server, subs, countries),
		"header": templates.DefaultHeader(nil),
	}
	index := views.Index(views.UI(r, UI))
	index.Render(ctx, w)

}

func (c *ServerController) serverUpdate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), PageTimeout)
	defer cancel()

	vars := mux.Vars(r)

	id, err := getInt64FromVars[storage.ServerID](vars, "id")
	if err != nil {
		c.writeJSONNotFound(w)
		return
	}

	server := &storage.VPNServer{}

	switch r.Header.Get("Content-Type") {
	case "application/json":
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(server); err != nil {
			writeJSON(w, http.StatusBadRequest, "Incorrect JSON format.")
			return
		}

	case "application/x-www-form-urlencoded":
		r.ParseForm()
		server.ParseURLValues(r.Form, entity.TimeLayout)
		fmt.Printf("form values \n%+v\n", r.Form)
	}

	if server.ID != id {
		if server.ID != 0 {
			writeJSON(w, http.StatusNotFound, ErrMsgServerIDMismatch)
			return
		}
		server.ID = id
	}

	fmt.Printf("updateServer get data:\n%+v\n", server)

	// If method PATCH, parse default value from old server.
	if r.Method == "PATCH" {
		oldServer, err := c.storage.GetServerByID(ctx, id)
		if err != nil {
			writeJSON(w, http.StatusNotFound, "Can't find server.")
			return
		}

		server.ParseDefaultsFrom(oldServer)

		fmt.Printf("updateServer parsed default data:\n%+v\n", server)
	}

	id, err = c.storage.SaveServer(ctx, server)
	if err != nil {
		log.Printf("[ERR] ServerController: serverUpdate: %v", err)
		c.writeJSONServerInternal(w)
		return
	}
	if id == 0 {
		log.Printf("[ERR] ServerController: serverUpdate: storage return 0 updated id")
		c.writeJSONServerInternal(w)
		return
	}

	resp := &Response{
		Success: true,
		Msg:     "Server updated.",
		Obj:     id,
	}
	writeJSON(w, 200, resp)
}

func (c *ServerController) writeErrorNotFound(w http.ResponseWriter, ctx context.Context) {
	writeError(w, ctx, http.StatusNotFound, ErrMsgServerNotFound)
}

func (c *ServerController) writeJSONNotFound(w http.ResponseWriter) {
	writeJSON(w, http.StatusNotFound, ErrMsgServerNotFound)
}
