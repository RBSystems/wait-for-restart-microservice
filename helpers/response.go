package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func CallCallback(request Request, status string) {
	request.CompletionTime = time.Now()
	request.Status = status

	response, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("Error marshalling: %s", err.Error())
		http.Post(request.CallbackAddress, "text/plain", bytes.NewBufferString("Error marshalling response"))
	}

	http.Post(request.CallbackAddress, "application/json", bytes.NewBuffer(response))
}
