package helpers

import "time"

type Configuration struct {
	WaitThreshold     int // How many items to have in our list before we stop waiting between iterations
	IterativeTime     int // How long to wait between iterations if the threshold isn't met
	IndividualTimeout int // Time in seconds to wait before timing out
}

type Request struct {
	MachineAddress  string    // Address to be tested for reboot
	Port            int       // The port to be used when testing connection
	Timeout         int       // Time in seconds to wait
	CallbackAddress string    // Complete address to send the notification that the host is responding
	SubmissionTime  time.Time // Will be filled by the service to indicate when the process started pinging
	CompletionTime  time.Time // Will be filled by the service to indicate when the machine responded or timed out
	Status          string    // Timeout or Success
	Identifier      string    // Optional value so the requester can identify the host when it's sent back
}
