package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

func sendResponse(req request, status string) {
	req.CompletionTime = time.Now()
	req.Status = status
	bits, err := json.Marshal(req) //Not sure what to do with the error here.

	if err != nil {
		fmt.Printf("Error marshalling: %s", err.Error())
		http.Post(req.CallbackAddress, "text/plain", bytes.NewBufferString("Error marshalling response."))
	}

	http.Post(req.CallbackAddress, "application/json", bytes.NewBuffer(bits))
}

func runService(submissionChannel <-chan request, config configuration) {

	var requestList []request

	for true { //start our loop, alive for the length of the service.

		//If we don't have anything in our list, we want to block and wait for
		//something to come in, so we're not wasting time.
		if len(requestList) < 1 {
			fmt.Printf("Waiting for entry into the channel.\n")
			req := <-submissionChannel
			fmt.Printf("Adding the item to the channel: %s\n", req.IPAddressHostname)
			req.SubmissionTime = time.Now()
			if req.Timeout == 0 {
				req.Timeout = 600
			}
			requestList = append(requestList, req)
		} else {
			select {
			case req := <-submissionChannel: //if there's something in the channel get it
				fmt.Printf("Adding item to channel %s.\n", req.IPAddressHostname)
				req.SubmissionTime = time.Now()
				requestList = append(requestList, req)
				if req.Timeout == 0 {
					req.Timeout = 600
				}

			default: //otherwise just bypass
			}
		}

		for curIndex := range requestList {
			curReq := requestList[curIndex]
			fmt.Printf("Pinging %s\n", curReq.IPAddressHostname)

			timeout := time.Duration(config.IndividualTimeout) * time.Second

			_, err := net.DialTimeout("tcp", curReq.IPAddressHostname+":"+strconv.Itoa(curReq.Port), timeout)

			if err == nil { //successfully connected
				sendResponse(curReq, "Success")
				fmt.Printf("Success!\n")
				requestList = append(requestList[:curIndex], requestList[curIndex+1:]...)
				continue
			}
			fmt.Printf("No response.\n")
			//we didn't connect, check the timeout.
			fmt.Printf("Time since init: %v\n", time.Since(curReq.SubmissionTime).Seconds())
			if int(time.Since(curReq.SubmissionTime).Seconds()) > curReq.Timeout { //We've timed out
				sendResponse(curReq, "Timeout")
				fmt.Printf("Failure, timeout.")
				requestList = append(requestList[:curIndex], requestList[curIndex+1:]...)
				continue
			}

		}
		if len(requestList) == 0 {
			continue
		} else if len(requestList) < config.WaitThreshold {
			time.Sleep(time.Duration(config.IterativeTime) * time.Second)
		}
	}
}
