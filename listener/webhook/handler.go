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
		log.Printf("webhook: Failed to format payload: %v", err)
		return
	}

	log.Println("===== webhook =====")
	// Process the received notification
	log.Printf("payload: %+v\n", string(formated))

	// Print the request headers
	if len(r.Header) > 0 {
		var headers []byte
		headers, _ = json.MarshalIndent(r.Header, "", "  ")
		log.Printf("headers: %s\n", string(headers))
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Notification received successfully!")
}
