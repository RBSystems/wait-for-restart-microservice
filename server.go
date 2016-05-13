package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/byuoitav/hateoas"
	"github.com/byuoitav/listen-for-reboot-microservice/controllers"
	"github.com/byuoitav/listen-for-reboot-microservice/helpers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"
)

func checkRequest(req helpers.Request) error {
	if len(req.CallbackAddress) < 1 || req.Port == 0 || len(req.IPAddressHostname) < 1 {
		return errors.New("Invalid payload")
	}

	return nil
}

// If we want the handler to have access to the channel we have to build a wrapper around it.
func makeSubmissonHandler(submissionChannel chan<- helpers.Request) func(c echo.Context) error {
	// This is our actual handler - submitRequest
	return func(c echo.Context) error {
		request := helpers.Request{}
		c.Bind(request)

		err := checkRequest(request)
		if err != nil {
			return c.String(http.StatusBadRequest, `Invalid request. Request must be in form of:
			  "IPAddressHostname": "string",
				"Port": int,
				"Timeout": int,
				"CallbackAddress": "string",
				"Identifier": "Optional string"
			}`)
		}

		if request.Timeout <= 10 {
			request.Timeout = 500
		}

		submissionChannel <- request // Add the request body to the channel queue

		return c.String(http.StatusOK, "Added to queue")
	}
}

func importConfig(configPath string) helpers.Configuration {
	fmt.Printf("Importing configuration information from %v\n", configPath)

	f, err := ioutil.ReadFile(configPath)
	check(err)

	var configurationData helpers.Configuration
	json.Unmarshal(f, &configurationData)

	fmt.Printf("\n%s\n", f)

	return configurationData
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	err := hateoas.Load("https://raw.githubusercontent.com/byuoitav/listen-for-reboot-microservice/master/swagger.yml")
	if err != nil {
		fmt.Println("Could not load Swagger file")
		panic(err)
	}

	var configFile = flag.String("config", "./config.json", "The locaton of the config file")

	flag.Parse()

	config := importConfig(*configFile)

	submissionChannel := make(chan helpers.Request, 50)

	submitRequest := makeSubmissonHandler(submissionChannel)

	go helpers.RunService(submissionChannel, config)

	port := ":8003"
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	e.Get("/", controllers.Root)
	e.Get("/health", controllers.Health)

	e.Post("/submit", submitRequest)

	fmt.Printf("Listen for Reboot microservice is listening on %s\n", port)
	e.Run(fasthttp.New(port))
}
