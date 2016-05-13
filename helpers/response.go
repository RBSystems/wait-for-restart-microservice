package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func SendResponse(req Request, status string) {
	req.CompletionTime = time.Now()
	req.Status = status
	bits, err := json.Marshal(req) // Not sure what to do with the error here
	if err != nil {
		fmt.Printf("Error marshalling: %s", err.Error())
		http.Post(req.CallbackAddress, "text/plain", bytes.NewBufferString("Error marshalling response"))
	}

	http.Post(req.CallbackAddress, "application/json", bytes.NewBuffer(bits))
}
