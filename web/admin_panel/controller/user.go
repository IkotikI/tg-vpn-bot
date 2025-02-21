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

const ErrMsgUserNotFound = "User not found."

type UserController struct {
	BaseController

	storage service.StorageService
	// service?
	// cookieStore *sessions.CookieStore
	// mux         http.Handler
}

func NewUserController(r *mux.Router, storage storage.Storage) *UserController {

	storage_service, err := service.NewStorageService(storage)
	if err != nil {
		log.Fatalf("[ERR} Can't start storage service: %v", err)
	}

	c := &UserController{
		storage: storage_service,
	}
	c.registerRoutes(r)
	return c
}

func (c *UserController) registerRoutes(r *mux.Router) {

	userRouter := r.PathPrefix("/").Subrouter()
	userRouter.HandleFunc("/user/{id:[0-9]+}", c.userView).Methods("GET")
	userRouter.HandleFunc("/user/{id:[0-9]+}", c.userUpdate).Methods("PUT", "PATCH")
}

func (c *UserController) userView(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), PageTimeout)
	defer cancel()

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		writeError(w, ctx, http.StatusNotFound, ErrMsgUserNotFound)
		return
	}

	idInt64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, ctx, http.StatusNotFound, ErrMsgUserNotFound)
		return
	}

	id := storage.UserID(idInt64)

	user, err := c.storage.GetUserByID(ctx, id)
	if err == storage.ErrNoSuchUser {
		writeError(w, ctx, http.StatusNotFound, ErrMsgUserNotFound)
		return
	} else if err != nil {
		log.Printf("[ERR] UserController: userView: %v", err)
		writeError(w, ctx, 500, ErrMsgServerInternalError)
		return
	}

	subs, err := c.storage.GetSubscriptionsWithServersByUserID(ctx, id, nil)
	if err != nil {
		log.Printf("[ERR] UserController: userView: %v", err)
		subs = &[]entity.SubscriptionWithServer{}
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
	idStr, ok := vars["id"]
	if !ok {
		writeJSON(w, http.StatusNotFound, ErrMsgUserNotFound)
		return
	}

	idInt64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ErrMsgUserNotFound)
		return
	}

	id := storage.UserID(idInt64)

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
			writeJSON(w, http.StatusNotFound, "User ID mismatch.")
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
		writeJSON(w, http.StatusInternalServerError, ErrMsgServerInternalError)
		return
	}
	if id == 0 {
		log.Printf("[ERR] UserController: userUpdate: storage return 0 updated id")
		writeJSON(w, http.StatusInternalServerError, ErrMsgServerInternalError)
		return
	}

	resp := &Response{
		Success: true,
		Msg:     "User updated.",
		Obj:     id,
	}
	writeJSON(w, 200, resp)
}
