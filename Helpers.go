package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ziutek/telnet"
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
			fmt.Printf("%s Adding the item to the channel: \n", req.IPAddressHostname)
			req.SubmissionTime = time.Now()
			requestList = append(requestList, req)
			continue //go back to get everything out of the channel that's there
		}
		select {
		case req := <-submissionChannel: //if there's something in the channel get it
			fmt.Printf("%s Adding item to channel.\n", req.IPAddressHostname)
			req.SubmissionTime = time.Now()
			requestList = append(requestList, req)
			continue //go back to get everything out of the channel that's there
		default: //otherwise just bypass
		}

		//we have to use a descending list otherwise our deletion gets in the way.
		for curIndex := len(requestList) - 1; curIndex >= 0; curIndex-- {
			curReq := requestList[curIndex]
			fmt.Printf("%s Pinging \n", curReq.IPAddressHostname)

			timeout := 85 * time.Millisecond

			//fmt.Printf("Timeout: %v", timeout)

			conn, err := net.DialTimeout("tcp", curReq.IPAddressHostname+":"+strconv.Itoa(curReq.Port), timeout)

			if err == nil { //successfully connected
				defer conn.Close()

				if !systemIsBusy(curReq) {
					sendResponse(curReq, "Success")
					fmt.Printf("%s Success!\n", curReq.IPAddressHostname)
					requestList = append(requestList[:curIndex], requestList[curIndex+1:]...)
					continue
				}
			}

			fmt.Printf("%s No response.\n", curReq.IPAddressHostname)
			//we didn't connect, check the timeout.
			fmt.Printf("%s Time since init: %v\n", curReq.IPAddressHostname, time.Since(curReq.SubmissionTime).Seconds())
			if int(time.Since(curReq.SubmissionTime).Seconds()) > curReq.Timeout { //We've timed out

				sendResponse(curReq, "Timeout")
				fmt.Printf("%s Failure, timeout %s.\n", curReq.IPAddressHostname, curReq.Timeout)

				requestList = append(requestList[:curIndex], requestList[curIndex+1:]...)
				continue
			}
		}

		if len(requestList) == 0 { //get back to wait for another request
			continue
		} else if len(requestList) < config.WaitThreshold {
			time.Sleep(time.Duration(config.IterativeTime) * time.Second)
		}
	}

}

//Check to make sure we're not getting a 'system is busy' error.

//We thought about using the telnet microservice to check if we get a 'system is busy'
//response, but that could wait for too long and chew up our process.
//TODO: make this not just work for folks that will respond to telnet over port 41795
func systemIsBusy(curReq request) bool {
	var conn *telnet.Conn

	conn, err := telnet.Dial("tcp", curReq.IPAddressHostname+":41795")

	if err != nil {
		return true
	}

	conn.SetUnixWriteMode(true) // Convert any '\n' (LF) to '\r\n' (CR LF) This is apparently very important
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	_, err = conn.Write([]byte("\n\n"))

	if err != nil {
		return true
	}

	//Dynamically get the prompt
	conn.SkipUntil(">")
	promptBytes, err := conn.ReadUntil(">")

	if err != nil {
		return true
	}
	regex := "\\S.*?>"

	re := regexp.MustCompile(regex)

	prompt := string(re.Find(promptBytes))

	_, err = conn.Write([]byte("hostname\n\n")) // Send a second newline so we get the prompt

	if err != nil {
		return true
	}

	err = conn.SkipUntil(prompt)

	if err != nil {
		return true
	}

	response, err := conn.ReadUntil(prompt) // Read until the second prompt delimiter (provided by sending two commands in sendCommand)

	if err != nil {
		return true
	}

	conn.Close()
	fmt.Printf("%s\n", string(response))

	if strings.Contains(string(response), "system is busy") {
		return true
	}

	return false
}
