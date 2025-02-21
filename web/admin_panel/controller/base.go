package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"vpn-tg-bot/pkg/e"
	"vpn-tg-bot/web/admin_panel/views"
)

type BaseController struct{}

func (c *BaseController) checkLogin(r *http.Request) bool {
	return true
}

type Response struct {
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Obj     interface{} `json:"obj"`
}

func writeError(w http.ResponseWriter, ctx context.Context, status int, msg string) {
	w.WriteHeader(status)
	index := views.Index(views.Error(status, msg))
	index.Render(ctx, w)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return e.Wrap("can't marshal json", err)
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return e.Wrap("can't write bytes", err)
	}
	return nil
}
