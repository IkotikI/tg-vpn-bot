package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"vpn-tg-bot/internal/service/subscription"
	"vpn-tg-bot/internal/storage"
	"vpn-tg-bot/web/admin_panel/entity"
	"vpn-tg-bot/web/admin_panel/views"
	"vpn-tg-bot/web/admin_panel/views/templates"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

const (
	ErrMsgUserNotFound   = "User not found."
	ErrMsgUserIDMismatch = "UserID mismatch."
)

type UserController struct {
	BaseController

	storage      storage.Storage
	subscription subscription.VPN_API
	// service?
	// cookieStore *sessions.CookieStore
	// mux         http.Handler
}

func NewUserController(r *mux.Router, storage storage.Storage, subscription subscription.VPN_API) *UserController {
	c := &UserController{
		storage:      storage,
		subscription: subscription,
	}
	c.registerRoutes(r)
	return c
}

func (c *UserController) registerRoutes(r *mux.Router) {

	userRouter := r.PathPrefix("/").Subrouter()
	userRouter.HandleFunc("/user/{id:[0-9]+}", c.userView).Methods("GET")
	userRouter.HandleFunc("/user/{id:[0-9]+}", c.userUpdate).Methods("PUT", "PATCH")
	userRouter.HandleFunc("/user/{id:[0-9]+}", c.userDelete).Methods("DELETE")

}

func (c *UserController) userView(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), PageTimeout)
	defer cancel()

	vars := mux.Vars(r)

	id, err := getInt64FromVars[storage.UserID](vars, "id")
	if err != nil {
		c.writeErrorNotFound(w, ctx)
		return
	}

	user, err := c.storage.GetUserByID(ctx, id)
	if err == storage.ErrNoSuchUser {
		c.writeErrorNotFound(w, ctx)
		return
	} else if err != nil {
		log.Printf("[ERR] UserController: userView: %v", err)
		c.writeErrorServerInternal(w, ctx)
		return
	}

	subs, err := c.storage.GetSubscriptionsWithServersByUserID(ctx, id, nil)
	if err != nil {
		log.Printf("[ERR] UserController: userView: %v", err)
		subs = &[]storage.SubscriptionWithServer{}
	}

	UI := map[string]templ.Component{
		"main":   templates.UserMain(user, subs),
		"header": templates.DefaultHeader(nil),
	}
	index := views.Index(views.UI(r, UI))
	index.Render(ctx, w)
}

func (c *UserController) userUpdate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), PageTimeout)
	defer cancel()

	vars := mux.Vars(r)

	id, err := getInt64FromVars[storage.UserID](vars, "id")
	if err != nil {
		c.writeJSONNotFound(w)
		return
	}

	user := &storage.User{}

	switch r.Header.Get("Content-Type") {
	case "application/json":
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(user); err != nil {
			writeJSON(w, http.StatusBadRequest, "Incorrect JSON format.")
			return
		}

	case "application/x-www-form-urlencoded":
		r.ParseForm()
		user.ParseURLValues(r.Form, entity.TimeLayout)
		fmt.Printf("form values \n%+v\n", r.Form)
	}

	if user.ID != id {
		if user.ID != 0 {
			writeJSON(w, http.StatusNotFound, ErrMsgUserIDMismatch)
			return
		}
		user.ID = id
	}

	fmt.Printf("updateUser get data:\n%+v\n", user)

	// If method PATCH, parse default value from old user.
	if r.Method == "PATCH" {
		oldUser, err := c.storage.GetUserByID(ctx, id)
		if err != nil {
			writeJSON(w, http.StatusNotFound, "Can't find user.")
			return
		}

		user.ParseDefaultsFrom(oldUser)

		fmt.Printf("updateUser parsed default data:\n%+v\n", user)
	}

	id, err = c.storage.SaveUser(ctx, user)
	if err != nil {
		log.Printf("[ERR] UserController: userUpdate: %v", err)
		c.writeJSONServerInternal(w)
		return
	}
	if id == 0 {
		log.Printf("[ERR] UserController: userUpdate: storage return 0 updated id")
		c.writeJSONServerInternal(w)
		return
	}

	resp := &Response{
		Success: true,
		Msg:     "User updated.",
		Obj:     id,
	}
	writeJSON(w, 200, resp)
}

func (c *UserController) userDelete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), PageTimeout)
	defer cancel()

	vars := mux.Vars(r)

	id, err := getInt64FromVars[storage.UserID](vars, "id")
	if err != nil {
		c.writeJSONNotFound(w)
		return
	}

	subs, err := c.storage.GetSubscriptionsByUserID(ctx, id)
	if err != nil && err != storage.ErrNoSuchSubscription {
		c.writeJSONServerInternal(w)
	}

	for _, sub := range *subs {
		err = c.subscription.DeleteUserSubscription(ctx, sub.ServerID, sub.UserID)
		if err != nil {
			log.Printf("[ERR] UserController: DeleteUser: %v", err)
			c.writeJSONServerInternal(w)
		}
		log.Printf("[INFO] UserController: DeleteUser: Subsciption deleted: userID: %d, serverID: %d", sub.UserID, sub.ServerID)
	}

	err = c.storage.RemoveUserByID(ctx, id)
	if err != nil {
		log.Printf("[ERR] UserController: DeleteUser: %v", err)
		c.writeJSONServerInternal(w)
		return
	}

	resp := &Response{
		Success: true,
		Msg:     "User deleted.",
		Obj:     id,
	}

	w.Header().Set("HX-Redirect", "/users")
	writeJSON(w, 200, resp)
}

func (c *UserController) writeErrorNotFound(w http.ResponseWriter, ctx context.Context) {
	writeError(w, ctx, http.StatusNotFound, ErrMsgUserNotFound)
}

func (c *UserController) writeJSONNotFound(w http.ResponseWriter) {
	writeJSON(w, http.StatusNotFound, ErrMsgUserNotFound)
}
