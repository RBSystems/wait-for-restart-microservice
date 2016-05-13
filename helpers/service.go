package helpers

import (
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
			req.SubmissionTime = time.Now()
			requestList = append(requestList, req)

			continue // Go back to get everything out of the channel
		}

		select {
		case req := <-submissionChannel: // If there's something in the channel get it
			req.SubmissionTime = time.Now()
			requestList = append(requestList, req)
			continue // Go back to get everything out of the channel that's there
		default: // Otherwise just bypass
		}

		// We have to use a descending list otherwise our deletion gets in the way
		for curIndex := len(requestList) - 1; curIndex >= 0; curIndex-- {
			curReq := requestList[curIndex]
			timeout := time.Duration(config.IndividualTimeout) * time.Millisecond

			conn, err := net.DialTimeout("tcp", curReq.MachineAddress+":"+strconv.Itoa(curReq.Port), timeout)
			if err == nil { // Successfully connected
				defer conn.Close()

				if !IsSystemBusy(curReq) {
					SendResponse(curReq, "Success")
					requestList = append(requestList[:curIndex], requestList[curIndex+1:]...)
					continue
				}
			}

			if int(time.Since(curReq.SubmissionTime).Seconds()) > curReq.Timeout { // We've timed out
				SendResponse(curReq, "Timeout")

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
