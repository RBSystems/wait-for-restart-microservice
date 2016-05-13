package controllers

import (
	"net/http"

	"github.com/byuoitav/ftp-microservice/helpers"
	"github.com/labstack/echo"
)

// SubmitInfo returns information about the /submit endpoint
func SubmitInfo(c echo.Context) error {
	response := &helpers.Response{
		Message: "Send a POST request to the /submit endpoint with a body including at least MachineAddress and CallbackAddress tokens",
	}

	return c.JSON(http.StatusOK, *response)
}
