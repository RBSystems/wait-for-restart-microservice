package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/byuoitav/hateoas"
	"github.com/byuoitav/wait-for-restart-microservice/handlers"
	"github.com/byuoitav/wait-for-restart-microservice/helpers"
	"github.com/byuoitav/wso2jwt"
	"github.com/jessemillar/health"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	err := hateoas.Load("https://raw.githubusercontent.com/byuoitav/wait-for-restart-microservice/master/swagger.json")
	if err != nil {
		log.Fatalln("Could not load Swagger file. Error: " + err.Error())
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
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	// Use the `secure` routing group to require authentication
	secure := router.Group("", echo.WrapMiddleware(wso2jwt.ValidateJWT))

	router.GET("/", echo.WrapHandler(http.HandlerFunc(hateoas.RootResponse)))
	router.GET("/health", echo.WrapHandler(http.HandlerFunc(health.Check)))

	secure.GET("/submit", handlers.SubmitInfo)

	secure.POST("/submit", submitRequest)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
