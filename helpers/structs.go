package helpers

import "time"

type Configuration struct {
	WaitThreshold     int // How many items to have in our list before we stop waiting between iterations
	IterativeTime     int // How long to wait between iterations if the threshold isn't met
	IndividualTimeout int // Time in seconds to wait before timing out
}

type Request struct {
	Address         string    `json:"Address"`         // Address to be tested for restart
	Port            int       `json:"port"`            // The port to be used when testing connection
	Timeout         int       `json:"timeout"`         // Time in seconds to wait
	CallbackAddress string    `json:"callbackAddress"` // Complete address to send the notification that the host is responding
	SubmissionTime  time.Time `json:"submissionTime"`  // Will be filled by the service to indicate when the process started pinging
	CompletionTime  time.Time `json:"completionTime"`  // Will be filled by the service to indicate when the machine responded or timed out
	Status          string    `json:"status"`          // Timeout or Success
	Identifier      string    `json:"identifier"`      // Optional value so the requester can identify the host when it's sent back
}
