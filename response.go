package messenger

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	// SendMessageURL is API endpoint for sending messages.
	SendMessageURL = "https://graph.facebook.com/v2.6/me/messages"
)

// Response is used for responding to events with messages.
type Response struct {
	token string
	to    Recipient
}

// Text sends a textual message.
func (r *Response) Text(message string) error {
	m := SendMessage{
		Recipient: r.to,
		Message: MessageData{
			Text: message,
		},
	}

	data, err := json.Marshal(m)
	if err != nil {
		return nil
	}

	req, err := http.NewRequest("POST", SendMessageURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.URL.RawQuery = "access_token=" + r.token

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	return err
}

// SendMessage is the information sent in an API request to Facebook.
type SendMessage struct {
	Recipient Recipient   `json:"recipient"`
	Message   MessageData `json:"message"`
}

// MessageData is a text message to be sent.
type MessageData struct {
	Text string `json:"text,omitempty"`
}
