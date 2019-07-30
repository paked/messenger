package messenger

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/xerrors"
)

const (
	// ProfileURL is the API endpoint used for retrieving profiles.
	// Used in the form: https://graph.facebook.com/v2.6/<USER_ID>?fields=<PROFILE_FIELDS>&access_token=<PAGE_ACCESS_TOKEN>
	ProfileURL = "https://graph.facebook.com/v2.6/"
	// SendSettingsURL is API endpoint for saving settings.
	SendSettingsURL = "https://graph.facebook.com/v2.6/me/thread_settings"

	// MessengerProfileURL is the API endpoint where you set properties that define various aspects of the following Messenger Platform features.
	// Used in the form https://graph.facebook.com/v2.6/me/messenger_profile?access_token=<PAGE_ACCESS_TOKEN>
	// https://developers.facebook.com/docs/messenger-platform/reference/messenger-profile-api/
	MessengerProfileURL = "https://graph.facebook.com/v2.6/me/messenger_profile"
)

// Options are the settings used when creating a Messenger client.
type Options struct {
	// Verify sets whether or not to be in the "verify" mode. Used for
	// verifying webhooks on the Facebook Developer Portal.
	Verify bool
	// AppSecret is the app secret from the Facebook Developer Portal. Used when
	// in the "verify" mode.
	AppSecret string
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

// OptInHandler is a handler used to handle opt-ins.
type OptInHandler func(OptIn, *Response)

// ReferralHandler is a handler used postback callbacks.
type ReferralHandler func(ReferralMessage, *Response)

// AccountLinkingHandler is a handler used to react to an account
// being linked or unlinked.
type AccountLinkingHandler func(AccountLinking, *Response)

// Messenger is the client which manages communication with the Messenger Platform API.
type Messenger struct {
	mux                    *http.ServeMux
	messageHandlers        []MessageHandler
	deliveryHandlers       []DeliveryHandler
	readHandlers           []ReadHandler
	postBackHandlers       []PostBackHandler
	optInHandlers          []OptInHandler
	referralHandlers       []ReferralHandler
	accountLinkingHandlers []AccountLinkingHandler
	token                  string
	verifyHandler          func(http.ResponseWriter, *http.Request)
	verify                 bool
	appSecret              string
}

// New creates a new Messenger. You pass in Options in order to affect settings.
func New(mo Options) *Messenger {
	if mo.Mux == nil {
		mo.Mux = http.NewServeMux()
	}

	m := &Messenger{
		mux:       mo.Mux,
		token:     mo.Token,
		verify:    mo.Verify,
		appSecret: mo.AppSecret,
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

// HandleOptIn adds a new OptInHandler to the Messenger which will be triggered
// once a user opts in to communicate with the bot.
func (m *Messenger) HandleOptIn(f OptInHandler) {
	m.optInHandlers = append(m.optInHandlers, f)
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

// HandleReferral adds a new ReferralHandler to the Messenger
func (m *Messenger) HandleReferral(f ReferralHandler) {
	m.referralHandlers = append(m.referralHandlers, f)
}

// HandleAccountLinking adds a new AccountLinkingHandler to the Messenger
func (m *Messenger) HandleAccountLinking(f AccountLinkingHandler) {
	m.accountLinkingHandlers = append(m.accountLinkingHandlers, f)
}

// Handler returns the Messenger in HTTP client form.
func (m *Messenger) Handler() http.Handler {
	return m.mux
}

// ProfileByID retrieves the Facebook user profile associated with that ID.
// According to the messenger docs: https://developers.facebook.com/docs/messenger-platform/identity/user-profile,
// Developers must ask for access except for some fields that are accessible without permissions.
//
// At the time of writing (2019-01-04), these fields are
// - Name
// - First Name
// - Last Name
// - Profile Picture
func (m *Messenger) ProfileByID(id int64, profileFields []string) (Profile, error) {
	p := Profile{}
	url := fmt.Sprintf("%v%v", ProfileURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return p, err
	}

	fields := strings.Join(profileFields, ",")

	req.URL.RawQuery = "fields=" + fields + "&access_token=" + m.token

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
			return p, xerrors.Errorf("facebook error: %w", qr.Error)
		}
	}

	return p, nil
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

// CallToActionsSetting sends settings for Get Started or Persistent Menu
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

	// consume a *copy* of the request body
	body, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	err := json.Unmarshal(body, &rec)
	if err != nil {
		err = xerrors.Errorf("could not decode response: %w", err)
		fmt.Println(err)
		fmt.Fprintln(w, `{status: 'not ok'}`)
		return
	}

	if rec.Object != "page" {
		fmt.Println("object is not page, undefined behaviour. Got", rec.Object)
	}

	if m.verify {
		if err := m.checkIntegrity(r); err != nil {
			fmt.Println("could not verify request:", err)
			fmt.Fprintln(w, `{status: 'not ok'}`)
			return
		}
	}

	m.dispatch(rec)

	fmt.Fprintln(w, `{status: 'ok'}`)
}

// checkIntegrity checks the integrity of the requests received
func (m *Messenger) checkIntegrity(r *http.Request) error {
	if m.appSecret == "" {
		return xerrors.New("missing app secret")
	}

	sigHeader := "X-Hub-Signature"
	sig := strings.SplitN(r.Header.Get(sigHeader), "=", 2)
	if len(sig) == 1 {
		if sig[0] == "" {
			return xerrors.Errorf("missing %s header", sigHeader)
		}
		return xerrors.Errorf("malformed %s header: %v", sigHeader, strings.Join(sig, "="))
	}

	checkSHA1 := func(body []byte, hash string) error {
		mac := hmac.New(sha1.New, []byte(m.appSecret))
		if mac.Write(body); fmt.Sprintf("%x", mac.Sum(nil)) != hash {
			return xerrors.Errorf("invalid signature: %s", hash)
		}
		return nil
	}

	body, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	sigEnc := strings.ToLower(sig[0])
	sigHash := strings.ToLower(sig[1])
	switch sigEnc {
	case "sha1":
		return checkSHA1(body, sigHash)
	default:
		return xerrors.Errorf("unknown %s header encoding, expected sha1: %s", sigHeader, sig[0])
	}
}

// dispatch triggers all of the relevant handlers when a webhook event is received.
func (m *Messenger) dispatch(r Receive) {
	for _, entry := range r.Entry {
		for _, info := range entry.Messaging {
			a := m.classify(info)
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
			case OptInAction:
				for _, f := range m.optInHandlers {
					message := *info.OptIn
					message.Sender = info.Sender
					message.Recipient = info.Recipient
					message.Time = time.Unix(info.Timestamp/int64(time.Microsecond), 0)
					f(message, resp)
				}
			case ReferralAction:
				for _, f := range m.referralHandlers {
					message := *info.ReferralMessage
					message.Sender = info.Sender
					message.Recipient = info.Recipient
					message.Time = time.Unix(info.Timestamp/int64(time.Microsecond), 0)
					f(message, resp)
				}
			case AccountLinkingAction:
				for _, f := range m.accountLinkingHandlers {
					message := *info.AccountLinking
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
func (m *Messenger) Send(to Recipient, message string, messagingType MessagingType, tags ...string) error {
	return m.SendWithReplies(to, message, nil, messagingType, tags...)
}

// SendGeneralMessage will send the GenericTemplate message
func (m *Messenger) SendGeneralMessage(to Recipient, elements *[]StructuredMessageElement, messagingType MessagingType, tags ...string) error {
	r := &Response{
		token: m.token,
		to:    to,
	}
	return r.GenericTemplate(elements, messagingType, tags...)
}

// SendWithReplies sends a textual message to a user, but gives them the option of numerous quick response options.
func (m *Messenger) SendWithReplies(to Recipient, message string, replies []QuickReply, messagingType MessagingType, tags ...string) error {
	response := &Response{
		token: m.token,
		to:    to,
	}

	return response.TextWithReplies(message, replies, messagingType, tags...)
}

// Attachment sends an image, sound, video or a regular file to a given recipient.
func (m *Messenger) Attachment(to Recipient, dataType AttachmentType, url string, messagingType MessagingType, tags ...string) error {
	response := &Response{
		token: m.token,
		to:    to,
	}

	return response.Attachment(dataType, url, messagingType, tags...)
}

// EnableChatExtension set the homepage url required for a chat extension.
func (m *Messenger) EnableChatExtension(homeURL HomeURL) error {
	wrap := map[string]interface{}{
		"home_url": homeURL,
	}
	data, err := json.Marshal(wrap)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", MessengerProfileURL, bytes.NewBuffer(data))
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

// classify determines what type of message a webhook event is.
func (m *Messenger) classify(info MessageInfo) Action {
	if info.Message != nil {
		return TextAction
	} else if info.Delivery != nil {
		return DeliveryAction
	} else if info.Read != nil {
		return ReadAction
	} else if info.PostBack != nil {
		return PostBackAction
	} else if info.OptIn != nil {
		return OptInAction
	} else if info.ReferralMessage != nil {
		return ReferralAction
	} else if info.AccountLinking != nil {
		return AccountLinkingAction
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
