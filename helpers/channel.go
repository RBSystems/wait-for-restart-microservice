package helpers

import (
	"net/http"

	"github.com/labstack/echo"
)

func MakeSubmissonHandler(submissionChannel chan<- Request) func(c echo.Context) error {
	// If we want the handler to have access to the channel we have to build a wrapper around it
	return func(context echo.Context) error {
		request := Request{}
		context.Bind(&request)

		if len(request.CallbackAddress) < 1 || len(request.Address) < 1 {
			return context.JSON(http.StatusBadRequest, "Request must include at least Address and CallbackAddress tokens")
		}

		if request.Port == 0 {
			request.Port = 23
		}

		if request.Timeout <= 10 {
			request.Timeout = 500
		}

		submissionChannel <- request // Add the request body to the channel queue

		return context.JSON(http.StatusOK, "Added to queue")
	}
}
