package controllers

import (
	"net/http"

	"github.com/jessemillar/jsonresp"
	"github.com/labstack/echo"
)

// SubmitInfo returns information about the /submit endpoint
func SubmitInfo(c echo.Context) error {
	return jsonresp.Create(c, http.StatusOK, "Send a POST request to the /submit endpoint with a body including at least Address and CallbackAddress tokens")
}
