package helpers

import "time"

type Configuration struct {
	WaitThreshold     int //How many items to have in our list before we stop waiting between iterations
	IterativeTime     int //How long to wait between iterations if the threshold isn't met.
	IndividualTimeout int //time in seconds to wait before timing out indivi
}

type Request struct {
	IPAddressHostname string    //hostname to be pinged
	Port              int       //port to be used when testing connection
	Timeout           int       //Time in seconds to wait. Optional, will default to 300 seconds if not present or is 0.
	CallbackAddress   string    //complete address to send the notification that the host is responding
	SubmissionTime    time.Time //Will be filled by the server as the time the process started pinging
	CompletionTime    time.Time //Will be filled by the service as the time that a) Sucessfully responded or b) timed out
	Status            string    //Timeout or Success
	Identifier        string    //Optional value to be passed in so the requester can identify the host when it's sent back.
}
