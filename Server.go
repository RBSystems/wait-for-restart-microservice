package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

//If we want the hanlder to have access to the channel we have to build a wrapper around it.
func makeSubmissonHandler(submissionChannel chan<- request) func(web.C, http.ResponseWriter, *http.Request) {
	//This is our actual handler - submitRequest
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		bits, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not read request body: %s\n", err.Error())
			return
		}

		var req request

		err = json.Unmarshal(bits, &req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error with the request body: %s", err.Error())
			return
		}
		submissionChannel <- req //add the request body to the channel queue

		fmt.Fprintf(w, "Added to queue.")
	}
}

func importConfig(configPath string) configuration {
	fmt.Printf("Importing the configuration information from %v\n", configPath)

	f, err := ioutil.ReadFile(configPath)
	check(err)

	var configurationData configuration
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
	var ConfigFileLocation = flag.String("config", "./config.json", "The locaton of the config file.")

	flag.Parse()

	config := importConfig(*ConfigFileLocation)

	submissionChannel := make(chan request, 50)

	submitRequest := makeSubmissonHandler(submissionChannel)

	go runService(submissionChannel, config)

	goji.Post("/submit", submitRequest)
	goji.Post("/submit/", submitRequest)
	//TODO: add some way to check the status (to make sure it hasn't gotten lost).
	goji.Serve()
}
