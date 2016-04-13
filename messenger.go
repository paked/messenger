package messenger

import (
	"fmt"
	"net/http"
)

const (
	WebhookURL = "/webhook"
)

type MessengerOptions struct {
	Verify      bool
	VerifyToken string
}

type Messenger struct {
	mux *http.ServeMux
}

func New(mo MessengerOptions) *Messenger {
	m := &Messenger{
		mux: http.NewServeMux(),
	}

	if mo.Verify {
		m.mux.HandleFunc(WebhookURL, newVerifyHandler(mo.VerifyToken))
	}

	return m
}

func (m *Messenger) Handler() http.Handler {
	return m.mux
}

func newVerifyHandler(token string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("hub.verify_token") == token {
			fmt.Fprintln(w, r.FormValue("hub.challenge"))
			return
		}

		fmt.Fprintln(w, "Incorrect verify token.")
	}
}
