package messenger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	WebhookURL = "/webhook"
)

type MessengerOptions struct {
	Verify      bool
	VerifyToken string
}

type MessageHandler func(Message)

type Messenger struct {
	mux      *http.ServeMux
	handlers map[Action]MessageHandler
}

func New(mo MessengerOptions) *Messenger {
	m := &Messenger{
		mux:      http.NewServeMux(),
		handlers: make(map[Action]MessageHandler),
	}

	if mo.Verify {
		m.mux.HandleFunc(WebhookURL, newVerifyHandler(mo.VerifyToken))
	} else {
		m.mux.HandleFunc(WebhookURL, m.handle)
	}

	return m
}

func (m *Messenger) Handle(a Action, f MessageHandler) {
	m.handlers[a] = f
}

func (m *Messenger) Handler() http.Handler {
	return m.mux
}

func (m *Messenger) handle(w http.ResponseWriter, r *http.Request) {
	var rec Receive

	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		fmt.Println(err)

		fmt.Fprintln(w, `{status: 'not ok'}`)
		return
	}

	if rec.Object != "page" {
		fmt.Println("Object is not page, undefined behaviour. Got", rec.Object)
	}

	m.dispatch(rec)

	fmt.Fprintln(w, `{status: 'ok'}`)
}

func (m *Messenger) dispatch(r Receive) {
	for _, entry := range r.Entry {
		for _, info := range entry.Messaging {
			a := m.classify(info, entry)
			if a == UnknownAction {
				fmt.Println("Unknown action:", info)
				continue
			}

			if f := m.handlers[a]; f != nil {
				f(Message{
					Sender:    info.Sender.ID,
					Recipient: info.Recipient.ID,
					Time:      time.Unix(info.Timestamp, 0),
					Text:      info.Message.Text,
				})
			}
		}
	}
}

func (m *Messenger) classify(info MessageInfo, e Entry) Action {
	if info.Message != nil {
		return TextAction
	}

	return UnknownAction
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
