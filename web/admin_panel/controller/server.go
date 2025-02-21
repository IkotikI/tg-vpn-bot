package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/web/admin_panel/entity"
	"vpn-tg-bot/web/admin_panel/service"
	"vpn-tg-bot/web/admin_panel/views"
	"vpn-tg-bot/web/admin_panel/views/templates"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

const ErrMsgServerNotFound = "Server not found."

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
	idStr, ok := vars["id"]
	if !ok {
		writeError(w, ctx, 404, ErrMsgServerNotFound)
		return
	}

	idInt64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, ctx, 404, ErrMsgServerNotFound)
		return
	}

	id := storage.ServerID(idInt64)

	server, err := c.storage.GetEntityServerByID(ctx, id)
	if err == storage.ErrNoSuchServer {
		writeError(w, ctx, 404, ErrMsgServerNotFound)
		return
	} else if err != nil {
		log.Printf("[ERR] ServerController: ServerView: %v", err)
		writeError(w, ctx, 500, ErrMsgServerInternalError)
		return
	}

	subs, err := c.storage.GetSubscriptionsWithUsersByServerID(ctx, id, nil)
	if err != nil {
		log.Printf("[ERR] ServerController: ServerView: %v", err)
		writeError(w, ctx, 500, ErrMsgServerInternalError)
		return
	}

	UI := map[string]templ.Component{
		"main":   templates.ServerMain(server, subs),
		"header": templates.DefaultHeader(nil),
	}
	index := views.Index(views.UI(r, UI))
	index.Render(ctx, w)

}

func (c *ServerController) serverUpdate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), PageTimeout)
	defer cancel()

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		writeJSON(w, http.StatusNotFound, ErrMsgServerNotFound)
		return
	}

	idInt64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ErrMsgServerNotFound)
		return
	}

	id := storage.ServerID(idInt64)

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
			writeJSON(w, http.StatusNotFound, "Server ID mismatch.")
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
		writeJSON(w, http.StatusInternalServerError, ErrMsgServerInternalError)
		return
	}
	if id == 0 {
		log.Printf("[ERR] ServerController: serverUpdate: storage return 0 updated id")
		writeJSON(w, http.StatusInternalServerError, ErrMsgServerInternalError)
		return
	}

	resp := &Response{
		Success: true,
		Msg:     "Server updated.",
		Obj:     id,
	}
	writeJSON(w, 200, resp)
}
