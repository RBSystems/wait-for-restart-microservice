package handlers

import (
	"net/http"

	"github.com/jessemillar/jsonresp"
	"github.com/labstack/echo"
)

// SubmitInfo returns information about the /submit endpoint
func SubmitInfo(context echo.Context) error {
	jsonresp.New(context.Response(), http.StatusOK, "Send a POST request to the /submit endpoint with a body including at least Address and CallbackAddress tokens")
	return nil
}
