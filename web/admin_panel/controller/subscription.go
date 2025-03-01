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

var ErrMsgSubscriptionNotFound = "Subscription not found."

type SubscriptionController struct {
	BaseController

	storage service.StorageService
	// service?
	// cookieStore *sessions.CookieStore
	// mux         http.Handler
}

func NewSubscriptionController(r *mux.Router, storage storage.Storage) *SubscriptionController {

	storage_service, err := service.NewStorageService(storage)
	if err != nil {
		log.Fatalf("[ERR} Can't start storage service: %v", err)
	}

	c := &SubscriptionController{
		storage: storage_service,
	}
	c.registerRoutes(r)
	return c
}

func (c *SubscriptionController) registerRoutes(r *mux.Router) {

	subscriptionRouter := r.PathPrefix("/").Subrouter()
	subscriptionRouter.HandleFunc("/subscription/{ServerID:[0-9]+}/{UserID:[0-9]+}", c.subscriptionView).Methods("GET")
	subscriptionRouter.HandleFunc("/subscription/{ServerID:[0-9]+}/{UserID:[0-9]+}", c.subscriptionUpdate).Methods("PUT", "PATCH")
}

func (c *SubscriptionController) subscriptionView(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), PageTimeout)
	defer cancel()

	vars := mux.Vars(r)

	ServerID, err := getInt64FromVars[storage.ServerID](vars, "ServerID")
	if err != nil {
		c.writeErrorNotFound(w, ctx)
		return
	}
	UserID, err := getInt64FromVars[storage.UserID](vars, "UserID")
	if err != nil {
		c.writeErrorNotFound(w, ctx)
		return
	}

	log.Printf("userID: %d, serverID: %d", UserID, ServerID)

	sub, err := c.storage.GetSubscriptionWithUserAndServerByIDs(ctx, UserID, ServerID)
	if err == storage.ErrNoSuchSubscription {
		c.writeErrorNotFound(w, ctx)
		return
	} else if err != nil {
		log.Printf("[ERR] SubscriptionController: subscriptionView: %v", err)
		c.writeErrorServerInternal(w, ctx)
		return
	}
	// var statusMsgs []string
	// sub := &entity.SubscriptionWithUserAndServer{}

	// subscription, err := c.storage.GetSubscriptionByIDs(ctx, UserID, ServerID)
	// if err != nil {
	// 	writeError(w, ctx, http.StatusNotFound, ErrMsgUserNotFound)
	// 	return
	// }

	// user, err := c.storage.GetUserByID(ctx, UserID)
	// if err != nil {
	// 	sub.User = entity.User{}
	// 	statusMsgs = append(statusMsgs, "Error get User")
	// 	log.Printf("[ERR] SubscriptionController: subscriptionView: %v", err)
	// }

	// server, err := c.storage.GetEntityServerByID(ctx, ServerID)
	// if err != nil {
	// 	sub.Server = entity.Server{}
	// 	statusMsgs = append(statusMsgs, "Error get Server")
	// 	log.Printf("[ERR] SubscriptionController: subscriptionView: %v", err)
	// }

	// log.Printf("got subscription: \n%+v\n user: \n%+v\n server: \n%+v\n", subscription, user, server)

	// sub.Subscription.ParseDefaultsFrom(subscription)
	// sub.User.ParseDefaultsFrom(user)
	// sub.Server.ParseDefaultsFrom(server)

	log.Printf("got subscription before render: \n%+v\n", sub)

	UI := map[string]templ.Component{
		"main": templates.SubscriptionMain(sub),
		// "status": templ.Raw(strings.Join(statusMsgs, "<br>")),
		"header": templates.DefaultHeader(nil),
	}

	index := views.Index(views.UI(r, UI))
	index.Render(ctx, w)
}

func (c *SubscriptionController) subscriptionUpdate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), PageTimeout)
	defer cancel()

	vars := mux.Vars(r)

	ServerID, err := getInt64FromVars[storage.ServerID](vars, "ServerID")
	if err != nil {
		c.writeErrorNotFound(w, ctx)
		return
	}
	UserID, err := getInt64FromVars[storage.UserID](vars, "UserID")
	if err != nil {
		c.writeErrorNotFound(w, ctx)
		return
	}

	sub := &storage.Subscription{}

	switch r.Header.Get("Content-Type") {
	case "application/json":
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(sub); err != nil {
			writeJSON(w, http.StatusBadRequest, "Incorrect JSON format.")
			return
		}

	case "application/x-www-form-urlencoded":
		r.ParseForm()
		sub.ParseURLValues(r.Form, entity.TimeLayout)
		fmt.Printf("form values \n%+v\n", r.Form)
	}

	if sub.UserID != UserID {
		if sub.UserID != 0 {
			writeJSON(w, http.StatusNotFound, ErrMsgUserIDMismatch)
			return
		}
		sub.UserID = UserID
	}

	if sub.ServerID != ServerID {
		if sub.ServerID != 0 {
			writeJSON(w, http.StatusNotFound, ErrMsgServerIDMismatch)
			return
		}
		sub.ServerID = ServerID
	}

	fmt.Printf("updateUser get data:\n%+v\n", sub)

	// If method PATCH, parse default value from old user.
	if r.Method == "PATCH" {
		oldSub, err := c.storage.GetSubscriptionByIDs(ctx, UserID, ServerID)
		if err != nil {
			writeJSON(w, http.StatusNotFound, "Can't find user.")
			return
		}

		sub.ParseDefaultsFrom(oldSub)

		fmt.Printf("updateSubscription parsed default data:\n%+v\n", sub)
	}

	err = c.storage.SaveSubscription(ctx, sub)
	if err != nil {
		log.Printf("[ERR] SubscriptionController: subUpdate: %v", err)
		writeJSON(w, http.StatusInternalServerError, ErrMsgServerInternalError)
		return
	}

	resp := &Response{
		Success: true,
		Msg:     "Subscription updated.",
		Obj:     sub,
	}
	writeJSON(w, 200, resp)
}

func (c *SubscriptionController) writeErrorNotFound(w http.ResponseWriter, ctx context.Context) {
	writeError(w, ctx, http.StatusNotFound, ErrMsgSubscriptionNotFound)
}

func (c *SubscriptionController) writeErrorServerInternal(w http.ResponseWriter, ctx context.Context) {
	writeError(w, ctx, http.StatusInternalServerError, ErrMsgServerInternalError)
}
