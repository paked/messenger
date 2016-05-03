package messenger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	// ProfileURL is the API endpoint used for retrieving profiles.
	// Used in the form: https://graph.facebook.com/v2.6/<USER_ID>?fields=first_name,last_name,profile_pic&access_token=<PAGE_ACCESS_TOKEN>
	ProfileURL = "https://graph.facebook.com/v2.6/"
)

// Options are the settings used when creating a Messenger client.
type Options struct {
	// Verify sets whether or not to be in the "verify" mode. Used for
	// verifying webhooks on the Facebook Developer Portal.
	Verify bool
	// VerifyToken is the token to be used when verifying the webhook. Is set
	// when the webhook is created.
	VerifyToken string
	// Token is the access token of the Facebook page to send messages from.
	Token string
    // WebhookURL is where the Messenger client should listen for webhook events.
	WebhookURL string
}

// MessageHandler is a handler used for responding to a message containing text.
type MessageHandler func(Message, *Response)

// DeliveryHandler is a handler used for responding to a read receipt.
type DeliveryHandler func(Delivery, *Response)

// PostBackHandler is a handler used postback callbacks.
type PostBackHandler func(PostBack, *Response)

// Messenger is the client which manages communication with the Messenger Platform API.
type Messenger struct {
	mux              *http.ServeMux
	messageHandlers  []MessageHandler
	deliveryHandlers []DeliveryHandler
    postBackHandlers []PostBackHandler
	token            string
	verifyHandler    func(http.ResponseWriter, *http.Request)
}

// New creates a new Messenger. You pass in Options in order to affect settings.
func New(mo Options) *Messenger {
	m := &Messenger{
		mux:   http.NewServeMux(),
		token: mo.Token,
	}

	m.verifyHandler = newVerifyHandler(mo.VerifyToken)
	m.mux.HandleFunc(mo.WebhookURL, m.handle)

	return m
}

// HandleMessage adds a new MessageHandler to the Messenger which will be triggered
// when a message is received by the client.
func (m *Messenger) HandleMessage(f MessageHandler) {
	m.messageHandlers = append(m.messageHandlers, f)
}

// HandleDelivery adds a new DeliveryHandler to the Messenger which will be triggered
// when a previously sent message is read by the recipient.
func (m *Messenger) HandleDelivery(f DeliveryHandler) {
	m.deliveryHandlers = append(m.deliveryHandlers, f)
}

// HandlePostBack adds a new PostBackHandler to the Messenger
func (m *Messenger) HandlePostBack(f PostBackHandler) {
	m.postBackHandlers = append(m.postBackHandlers, f)
}

// Handler returns the Messenger in HTTP client form.
func (m *Messenger) Handler() http.Handler {
	return m.mux
}

// ProfileByID retrieves the Facebook user associated with that ID
func (m *Messenger) ProfileByID(id int64) (Profile, error) {
	p := Profile{}
	url := fmt.Sprintf("%v%v", ProfileURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return p, err
	}

	req.URL.RawQuery = "fields=first_name,last_name,profile_pic&access_token=" + m.token

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&p)

	return p, err
}

// handle is the internal HTTP handler for the webhooks.
func (m *Messenger) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		m.verifyHandler(w, r)
		return
	}

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

// dispatch triggers all of the relevant handlers when a webhook event is received.
func (m *Messenger) dispatch(r Receive) {
	for _, entry := range r.Entry {
		for _, info := range entry.Messaging {
			a := m.classify(info, entry)
			if a == UnknownAction {
				fmt.Println("Unknown action:", info)
				continue
			}

			resp := &Response{
				to:    Recipient{info.Sender.ID},
				token: m.token,
			}

			switch a {
                case TextAction:
                    for _, f := range m.messageHandlers {
                        message := *info.Message
                        message.Sender = info.Sender
                        message.Recipient = info.Recipient
                        message.Time = time.Unix(info.Timestamp, 0)
                        f(message, resp)
                    }
                case DeliveryAction:
                    for _, f := range m.deliveryHandlers {
                        f(*info.Delivery, resp)
                    }
                case PostBackAction:
                    for _, f := range m.postBackHandlers {
                        message := *info.PostBack
                        message.Sender = info.Sender
                        message.Recipient = info.Recipient
                        message.Time = time.Unix(info.Timestamp, 0)
                        f(message, resp)
                    }
			}
		}
	}
}

// classify determines what type of message a webhook event is.
func (m *Messenger) classify(info MessageInfo, e Entry) Action {
	if info.Message != nil {
		return TextAction
	} else if info.Delivery != nil {
		return DeliveryAction
	} else if info.PostBack != nil {
		return PostBackAction
	}
	return UnknownAction
}

// newVerifyHandler returns a function which can be used to handle webhook verification
func newVerifyHandler(token string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("hub.verify_token") == token {
			fmt.Fprintln(w, r.FormValue("hub.challenge"))
			return
		}
		fmt.Fprintln(w, "Incorrect verify token.")
	}
}
