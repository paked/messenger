package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

// AttachmentType is attachment type.
type AttachmentType string
type MessagingType string
type TopElementStyle string
type ImageAspectRatio string

const (
	// SendMessageURL is API endpoint for sending messages.
	SendMessageURL = "https://graph.facebook.com/v2.11/me/messages"

	// ImageAttachment is image attachment type.
	ImageAttachment AttachmentType = "image"
	// AudioAttachment is audio attachment type.
	AudioAttachment AttachmentType = "audio"
	// VideoAttachment is video attachment type.
	VideoAttachment AttachmentType = "video"
	// FileAttachment is file attachment type.
	FileAttachment AttachmentType = "file"

	// ResponseType is response messaging type
	ResponseType MessagingType = "RESPONSE"
	// UpdateType is update messaging type
	UpdateType MessagingType = "UPDATE"
	// MessageTagType is message_tag messaging type
	MessageTagType MessagingType = "MESSAGE_TAG"
	// NonPromotionalSubscriptionType is NON_PROMOTIONAL_SUBSCRIPTION messaging type
	NonPromotionalSubscriptionType MessagingType = "NON_PROMOTIONAL_SUBSCRIPTION"

	// TopElementStyle is compact.
	CompactTopElementStyle TopElementStyle = "compact"
	// TopElementStyle is large.
	LargeTopElementStyle TopElementStyle = "large"

	// ImageAspectRatio is horizontal (1.91:1). Default.
	HorizontalImageAspectRatio ImageAspectRatio = "horizontal"
	// ImageAspectRatio is square.
	SquareImageAspectRatio ImageAspectRatio = "square"
)

// QueryResponse is the response sent back by Facebook when setting up things
// like greetings or call-to-actions
type QueryResponse struct {
	Error  *QueryError `json:"error,omitempty"`
	Result string      `json:"result,omitempty"`
}

// QueryError is representing an error sent back by Facebook
type QueryError struct {
	Message   string `json:"message"`
	Type      string `json:"type"`
	Code      int    `json:"code"`
	FBTraceID string `json:"fbtrace_id"`
}

func checkFacebookError(r io.Reader) error {
	var err error

	qr := QueryResponse{}
	err = json.NewDecoder(r).Decode(&qr)
	if qr.Error != nil {
		err = fmt.Errorf("Facebook error : %s", qr.Error.Message)
		return err
	}

	return nil
}

// Response is used for responding to events with messages.
type Response struct {
	token string
	to    Recipient
}

// SetToken is for using DispatchMessage from outside.
func (r *Response) SetToken(token string) {
	r.token = token
}

// Text sends a textual message.
func (r *Response) Text(message string, messagingType MessagingType, tags ...string) error {
	return r.TextWithReplies(message, nil, messagingType, tags...)
}

// TextWithReplies sends a textual message with some replies
// messagingType should be one of the following: "RESPONSE","UPDATE","MESSAGE_TAG","NON_PROMOTIONAL_SUBSCRIPTION"
// only supply tags when messagingType == "MESSAGE_TAG" (see https://developers.facebook.com/docs/messenger-platform/send-messages#messaging_types for more)
func (r *Response) TextWithReplies(message string, replies []QuickReply, messagingType MessagingType, tags ...string) error {
	var tag string
	if len(tags) > 0 {
		tag = tags[0]
	}

	m := SendMessage{
		MessagingType: messagingType,
		Recipient:     r.to,
		Message: MessageData{
			Text:         message,
			Attachment:   nil,
			QuickReplies: replies,
		},
		Tag: tag,
	}
	return r.DispatchMessage(&m)
}

// AttachmentWithReplies sends a attachment message with some replies
func (r *Response) AttachmentWithReplies(attachment *StructuredMessageAttachment, replies []QuickReply, messagingType MessagingType, tags ...string) error {
	var tag string
	if len(tags) > 0 {
		tag = tags[0]
	}

	m := SendMessage{
		MessagingType: messagingType,
		Recipient:     r.to,
		Message: MessageData{
			Attachment:   attachment,
			QuickReplies: replies,
		},
		Tag: tag,
	}
	return r.DispatchMessage(&m)
}

// Image sends an image.
func (r *Response) Image(im image.Image) error {
	imageBytes := new(bytes.Buffer)
	err := jpeg.Encode(imageBytes, im, nil)
	if err != nil {
		return err
	}

	return r.AttachmentData(ImageAttachment, "meme.jpg", imageBytes)
}

// Attachment sends an image, sound, video or a regular file to a chat.
func (r *Response) Attachment(dataType AttachmentType, url string, messagingType MessagingType, tags ...string) error {
	var tag string
	if len(tags) > 0 {
		tag = tags[0]
	}

	m := SendStructuredMessage{
		MessagingType: messagingType,
		Recipient:     r.to,
		Message: StructuredMessageData{
			Attachment: StructuredMessageAttachment{
				Type: dataType,
				Payload: StructuredMessagePayload{
					Url: url,
				},
			},
		},
		Tag: tag,
	}
	return r.DispatchMessage(&m)
}

// copied from multipart package
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

// copied from multipart package
func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// copied from multipart package with slight changes due to fixed content-type there
func createFormFile(filename string, w *multipart.Writer, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="filedata"; filename="%s"`,
			escapeQuotes(filename)))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}

// AttachmentData sends an image, sound, video or a regular file to a chat via an io.Reader.
func (r *Response) AttachmentData(dataType AttachmentType, filename string, filedata io.Reader) error {

	filedataBytes, err := ioutil.ReadAll(filedata)
	if err != nil {
		return err
	}
	contentType := http.DetectContentType(filedataBytes[:512])
	fmt.Println("Content-type detected:", contentType)

	var body bytes.Buffer
	multipartWriter := multipart.NewWriter(&body)
	data, err := createFormFile(filename, multipartWriter, contentType)
	if err != nil {
		return err
	}

	_, err = bytes.NewBuffer(filedataBytes).WriteTo(data)
	if err != nil {
		return err
	}

	multipartWriter.WriteField("recipient", fmt.Sprintf(`{"id":"%v"}`, r.to.ID))
	multipartWriter.WriteField("message", fmt.Sprintf(`{"attachment":{"type":"%v", "payload":{}}}`, dataType))

	req, err := http.NewRequest("POST", SendMessageURL, &body)
	if err != nil {
		return err
	}

	req.URL.RawQuery = "access_token=" + r.token

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	return checkFacebookError(resp.Body)
}

// ButtonTemplate sends a message with the main contents being button elements
func (r *Response) ButtonTemplate(text string, buttons *[]StructuredMessageButton, messagingType MessagingType, tags ...string) error {
	var tag string
	if len(tags) > 0 {
		tag = tags[0]
	}

	m := SendStructuredMessage{
		MessagingType: messagingType,
		Recipient:     r.to,
		Message: StructuredMessageData{
			Attachment: StructuredMessageAttachment{
				Type: "template",
				Payload: StructuredMessagePayload{
					TemplateType: "button",
					Text:         text,
					Buttons:      buttons,
					Elements:     nil,
				},
			},
		},
		Tag: tag,
	}

	return r.DispatchMessage(&m)
}

// GenericTemplate is a message which allows for structural elements to be sent
func (r *Response) GenericTemplate(elements *[]StructuredMessageElement, messagingType MessagingType, tags ...string) error {
	var tag string
	if len(tags) > 0 {
		tag = tags[0]
	}

	m := SendStructuredMessage{
		MessagingType: messagingType,
		Recipient:     r.to,
		Message: StructuredMessageData{
			Attachment: StructuredMessageAttachment{
				Type: "template",
				Payload: StructuredMessagePayload{
					TemplateType: "generic",
					Buttons:      nil,
					Elements:     elements,
				},
			},
		},
		Tag: tag,
	}
	return r.DispatchMessage(&m)
}

// SenderAction sends a info about sender action
func (r *Response) SenderAction(action string) error {
	m := SendSenderAction{
		Recipient:    r.to,
		SenderAction: action,
	}
	return r.DispatchMessage(&m)
}

// DispatchMessage posts the message to messenger, return the error if there's any
func (r *Response) DispatchMessage(m interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", SendMessageURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.URL.RawQuery = "access_token=" + r.token

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return nil
	}
	return checkFacebookError(resp.Body)
}

// SendMessage is the information sent in an API request to Facebook.
type SendMessage struct {
	MessagingType MessagingType `json:"messaging_type"`
	Recipient     Recipient     `json:"recipient"`
	Message       MessageData   `json:"message"`
	Tag           string        `json:"tag,omitempty"`
}

// MessageData is a message consisting of text or an attachment, with an additional selection of optional quick replies.
type MessageData struct {
	Text         string                       `json:"text,omitempty"`
	Attachment   *StructuredMessageAttachment `json:"attachment,omitempty"`
	QuickReplies []QuickReply                 `json:"quick_replies,omitempty"`
}

// SendStructuredMessage is a structured message template.
type SendStructuredMessage struct {
	MessagingType MessagingType         `json:"messaging_type"`
	Recipient     Recipient             `json:"recipient"`
	Message       StructuredMessageData `json:"message"`
	Tag           string                `json:"tag,omitempty"`
}

// StructuredMessageData is an attachment sent with a structured message.
type StructuredMessageData struct {
	Attachment StructuredMessageAttachment `json:"attachment"`
}

// StructuredMessageAttachment is the attachment of a structured message.
type StructuredMessageAttachment struct {
	// Type must be template
	Type AttachmentType `json:"type"`
	// Payload is the information for the file which was sent in the attachment.
	Payload StructuredMessagePayload `json:"payload"`
}

// StructuredMessagePayload is the actual payload of an attachment
type StructuredMessagePayload struct {
	// TemplateType must be button, generic or receipt
	TemplateType     string                      `json:"template_type,omitempty"`
	TopElementStyle  TopElementStyle             `json:"top_element_style,omitempty"`
	Text             string                      `json:"text,omitempty"`
	ImageAspectRatio ImageAspectRatio            `json:"image_aspect_ratio,omitempty"`
	Sharable         bool                        `json:"sharable,omitempty"`
	Elements         *[]StructuredMessageElement `json:"elements,omitempty"`
	Buttons          *[]StructuredMessageButton  `json:"buttons,omitempty"`
	Url              string                      `json:"url,omitempty"`
	AttachmentID     string                      `json:"attachment_id,omitempty"`
}

// StructuredMessageElement is a response containing structural elements
type StructuredMessageElement struct {
	Title         string                    `json:"title"`
	ImageURL      string                    `json:"image_url"`
	ItemURL       string                    `json:"item_url,omitempty"`
	Subtitle      string                    `json:"subtitle"`
	DefaultAction *DefaultAction            `json:"default_action,omitempty"`
	Buttons       []StructuredMessageButton `json:"buttons"`
}

// DefaultAction is a response containing default action properties
type DefaultAction struct {
	Type                string `json:"type"`
	URL                 string `json:"url,omitempty"`
	WebviewHeightRatio  string `json:"webview_height_ratio,omitempty"`
	MessengerExtensions bool   `json:"messenger_extensions,omitempty"`
	FallbackURL         string `json:"fallback_url,omitempty"`
	WebviewShareButton  string `json:"webview_share_button,omitempty"`
}

// StructuredMessageButton is a response containing buttons
type StructuredMessageButton struct {
	Type                string `json:"type"`
	URL                 string `json:"url,omitempty"`
	Title               string `json:"title,omitempty"`
	Payload             string `json:"payload,omitempty"`
	WebviewHeightRatio  string `json:"webview_height_ratio,omitempty"`
	MessengerExtensions bool   `json:"messenger_extensions,omitempty"`
	FallbackURL         string `json:"fallback_url,omitempty"`
	WebviewShareButton  string `json:"webview_share_button,omitempty"`
}

// SendSenderAction is the information about sender action
type SendSenderAction struct {
	Recipient    Recipient `json:"recipient"`
	SenderAction string    `json:"sender_action"`
}
