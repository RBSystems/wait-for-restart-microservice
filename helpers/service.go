package helpers

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func RunService(submissionChannel <-chan Request, config Configuration) {
	var requestList []Request

	for true { // Start the loop that will stay alive for the length of the service
		// If we don't have anything in our list, block and wait for something to come in
		if len(requestList) < 1 {
			req := <-submissionChannel
			fmt.Printf("%s Adding item to channel: \n", req.MachineAddress)
			req.SubmissionTime = time.Now()
			requestList = append(requestList, req)

			continue // Go back to get everything out of the channel
		}

		select {
		case req := <-submissionChannel: // If there's something in the channel get it
			fmt.Printf("%s Adding item to channel\n", req.MachineAddress)
			req.SubmissionTime = time.Now()
			requestList = append(requestList, req)
			continue // Go back to get everything out of the channel that's there
		default: // Otherwise just bypass
		}

		// We have to use a descending list otherwise our deletion gets in the way
		for curIndex := len(requestList) - 1; curIndex >= 0; curIndex-- {
			curReq := requestList[curIndex]
			fmt.Printf("%s Pinging \n", curReq.MachineAddress)

			timeout := time.Duration(config.IndividualTimeout) * time.Millisecond

			conn, err := net.DialTimeout("tcp", curReq.MachineAddress+":"+strconv.Itoa(curReq.Port), timeout)
			if err == nil { // Successfully connected
				defer conn.Close()

				if !IsSystemBusy(curReq) {
					SendResponse(curReq, "Success")
					fmt.Printf("%s Success!\n", curReq.MachineAddress)
					requestList = append(requestList[:curIndex], requestList[curIndex+1:]...)
					continue
				}
			}

			fmt.Printf("%s No response\n", curReq.MachineAddress)
			// We didn't connect, check the timeout
			fmt.Printf("%s Time since init: %v\n", curReq.MachineAddress, time.Since(curReq.SubmissionTime).Seconds())

			if int(time.Since(curReq.SubmissionTime).Seconds()) > curReq.Timeout { // We've timed out
				SendResponse(curReq, "Timeout")
				fmt.Printf("%s Failure, timeout %v\n", curReq.MachineAddress, curReq.Timeout)

				requestList = append(requestList[:curIndex], requestList[curIndex+1:]...)
				continue
			}
		}

		if len(requestList) == 0 { // Get back to wait for another request
			continue
		} else if len(requestList) < config.WaitThreshold {
			time.Sleep(time.Duration(config.IterativeTime) * time.Second)
		}
	}
}
