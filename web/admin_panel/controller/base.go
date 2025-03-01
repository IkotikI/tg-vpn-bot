package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"vpn-tg-bot/pkg/e"
	"vpn-tg-bot/web/admin_panel/views"
)

const (
	ErrMsgNotFound            = "Not found."
	ErrMsgServerInternalError = "Server Internal Error."
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

func getInt64FromVars[T ~int64](vars map[string]string, key string) (T, error) {
	str, ok := vars[key]
	if !ok {
		return 0, fmt.Errorf("can't find %s in vars", key)
	}

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}

	return T(i), err
}

func (c *BaseController) writeErrorNotFound(w http.ResponseWriter, ctx context.Context) {
	writeError(w, ctx, http.StatusNotFound, ErrMsgNotFound)
}

func (c *BaseController) writeErrorServerInternal(w http.ResponseWriter, ctx context.Context) {
	writeError(w, ctx, http.StatusInternalServerError, ErrMsgServerInternalError)
}

func (c *BaseController) writeJSONNotFound(w http.ResponseWriter) {
	writeJSON(w, http.StatusNotFound, ErrMsgNotFound)
}

func (c *BaseController) writeJSONServerInternal(w http.ResponseWriter) {
	writeJSON(w, http.StatusInternalServerError, ErrMsgServerInternalError)
}
