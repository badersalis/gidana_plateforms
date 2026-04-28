package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type expoPushPayload struct {
	To    string `json:"to"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Data  any    `json:"data,omitempty"`
}

// SendExpoPush sends a push notification via the Expo push API.
// It is fire-and-forget; errors are logged but not returned.
func SendExpoPush(token, title, body string, data any) {
	if token == "" {
		return
	}
	payload := expoPushPayload{To: token, Title: title, Body: body, Data: data}
	b, _ := json.Marshal(payload)
	resp, err := http.Post("https://exp.host/api/v2/push/send", "application/json", bytes.NewReader(b))
	if err != nil {
		log.Printf("push: failed to send notification: %v", err)
		return
	}
	resp.Body.Close()
}
