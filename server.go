package main

import (
	"flag"
	"fmt"

	"github.com/byuoitav/hateoas"
	"github.com/byuoitav/listen-for-reboot-microservice/controllers"
	"github.com/byuoitav/listen-for-reboot-microservice/helpers"
	"github.com/jessemillar/health"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"
)

func main() {
	err := hateoas.Load("https://raw.githubusercontent.com/byuoitav/listen-for-reboot-microservice/master/swagger.yml")
	if err != nil {
		fmt.Println("Could not load Swagger file")
		panic(err)
	}

	var configFile = flag.String("config", "./config.json", "The location of the config file")

	flag.Parse()

	config, err := helpers.ImportConfig(*configFile)
	if err != nil {
		fmt.Println("Could not load config file")
		panic(err)
	}

	submissionChannel := make(chan helpers.Request, 50)

	submitRequest := helpers.MakeSubmissonHandler(submissionChannel)

	go helpers.RunService(submissionChannel, config)

	port := ":8003"
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	e.Get("/", controllers.Root)
	e.Get("/health", health.Check)
	e.Get("/submit", controllers.SubmitInfo)

	e.Post("/submit", submitRequest)

	fmt.Printf("Wait for Reboot microservice is listening on %s\n", port)
	e.Run(fasthttp.New(port))
}
