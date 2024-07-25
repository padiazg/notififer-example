package webhook

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	m "github.com/padiazg/notifier/model"
)

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming JSON payload
	var (
		notification m.Notification
		err          error
	)

	err = json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, "Failed to decode JSON payload", http.StatusBadRequest)
		return
	}

	formated, err := json.MarshalIndent(notification, "", "  ")
	if err != nil {
		log.Printf("Failed to format payload: %v", err)
		return
	}

	// Process the received notification
	log.Printf("%v Received notification: %+v\n", r.Method, string(formated))

	// Print the request headers
	var headers []byte
	headers, _ = json.MarshalIndent(r.Header, "", "  ")
	log.Printf("%v Headers: %s\n", r.Method, string(headers))

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Notification received successfully!")
}
