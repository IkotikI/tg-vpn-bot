package controller

import "net/http"

type BaseController struct{}

func (c *BaseController) checkLogin(r *http.Request) bool {
	return true
}
