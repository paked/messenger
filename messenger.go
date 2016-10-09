package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// ProfileURL is the API endpoint used for retrieving profiles.
	// Used in the form: https://graph.facebook.com/v2.6/<USER_ID>?fields=<PROFILE_FIELDS>&access_token=<PAGE_ACCESS_TOKEN>
	ProfileURL = "https://graph.facebook.com/v2.6/"
	// ProfileFields is a list of JSON field names which will be populated by the profile query.
	ProfileFields = "first_name,last_name,profile_pic,locale,timezone,gender"
	// SendSettingsURL is API endpoint for saving settings.
	SendSettingsURL = "https://graph.facebook.com/v2.6/me/thread_settings"
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
	// WebhookURL is where the Messenger client should listen for webhook events. Leaving the string blank implies a path of "/".
	WebhookURL string
	// Mux is shared mux between several Messenger objects
	Mux *http.ServeMux
}

// MessageHandler is a handler used for responding to a message containing text.
type MessageHandler func(Message, *Response)

// DeliveryHandler is a handler used for responding to a delivery receipt.
type DeliveryHandler func(Delivery, *Response)

// ReadHandler is a handler used for responding to a read receipt.
type ReadHandler func(Read, *Response)

// PostBackHandler is a handler used postback callbacks.
type PostBackHandler func(PostBack, *Response)

// Messenger is the client which manages communication with the Messenger Platform API.
type Messenger struct {
	mux              *http.ServeMux
	messageHandlers  []MessageHandler
	deliveryHandlers []DeliveryHandler
	readHandlers     []ReadHandler
	postBackHandlers []PostBackHandler
	token            string
	verifyHandler    func(http.ResponseWriter, *http.Request)
}

// New creates a new Messenger. You pass in Options in order to affect settings.
func New(mo Options) *Messenger {
	if mo.Mux == nil {
		mo.Mux = http.NewServeMux()
	}

	m := &Messenger{
		mux:   mo.Mux,
		token: mo.Token,
	}

	if mo.WebhookURL == "" {
		mo.WebhookURL = "/"
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
// when a previously sent message is delivered to the recipient.
func (m *Messenger) HandleDelivery(f DeliveryHandler) {
	m.deliveryHandlers = append(m.deliveryHandlers, f)
}

// HandleRead adds a new DeliveryHandler to the Messenger which will be triggered
// when a previously sent message is read by the recipient.
func (m *Messenger) HandleRead(f ReadHandler) {
	m.readHandlers = append(m.readHandlers, f)
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

	req.URL.RawQuery = "fields=" + ProfileFields + "&access_token=" + m.token

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return p, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return p, err
	}

	err = json.Unmarshal(content, &p)
	if err != nil {
		return p, err
	}

	if p == *new(Profile) {
		qr := QueryResponse{}
		err = json.Unmarshal(content, &qr)
		if qr.Error != nil {
			err = fmt.Errorf("Facebook error : %s", qr.Error.Message)
		}
	}

	return p, err
}

// GreetingSetting sends settings for greeting
func (m *Messenger) GreetingSetting(text string) error {
	d := GreetingSetting{
		SettingType: "greeting",
		Greeting: GreetingInfo{
			Text: text,
		},
	}

	data, err := json.Marshal(d)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", SendSettingsURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.URL.RawQuery = "access_token=" + m.token

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkFacebookError(resp.Body)
}

// CallToActionsSetting sends settings for Get Started or Persist Menu
func (m *Messenger) CallToActionsSetting(state string, actions []CallToActionsItem) error {
	d := CallToActionsSetting{
		SettingType:   "call_to_actions",
		ThreadState:   state,
		CallToActions: actions,
	}

	data, err := json.Marshal(d)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", SendSettingsURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.URL.RawQuery = "access_token=" + m.token

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkFacebookError(resp.Body)
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
		fmt.Println("could not decode response:", err)
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
					message.Time = time.Unix(info.Timestamp/int64(time.Microsecond), 0)
					f(message, resp)
				}
			case DeliveryAction:
				for _, f := range m.deliveryHandlers {
					f(*info.Delivery, resp)
				}
			case ReadAction:
				for _, f := range m.readHandlers {
					f(*info.Read, resp)
				}
			case PostBackAction:
				for _, f := range m.postBackHandlers {
					message := *info.PostBack
					message.Sender = info.Sender
					message.Recipient = info.Recipient
					message.Time = time.Unix(info.Timestamp/int64(time.Microsecond), 0)
					f(message, resp)
				}
			}
		}
	}
}

// Response returns new Response object
func (m *Messenger) Response(to int64) *Response {
	return &Response{
		to:    Recipient{to},
		token: m.token,
	}
}

// Send will send a textual message to a user. This user must have previously initiated a conversation with the bot.
func (m *Messenger) Send(to Recipient, message string) error {
	return m.SendWithReplies(to, message, nil)
}

// SendWithReplies sends a textual message to a user, but gives them the option of numerous quick response options.
func (m *Messenger) SendWithReplies(to Recipient, message string, replies []QuickReply) error {
	response := &Response{
		token: m.token,
		to:    to,
	}

	return response.TextWithReplies(message, replies)
}

// Attachment sends an image, sound, video or a regular file to a given recipient.
func (m *Messenger) Attachment(to Recipient, dataType AttachmentType, url string) error {
	response := &Response{
		token: m.token,
		to:    to,
	}

	return response.Attachment(dataType, url)
}

// classify determines what type of message a webhook event is.
func (m *Messenger) classify(info MessageInfo, e Entry) Action {
	if info.Message != nil {
		return TextAction
	} else if info.Delivery != nil {
		return DeliveryAction
	} else if info.Read != nil {
		return ReadAction
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
