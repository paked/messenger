package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
)

// AttachmentType is attachment type.
type AttachmentType string

const (
	// SendMessageURL is API endpoint for sending messages.
	SendMessageURL = "https://graph.facebook.com/v2.6/me/messages"

	// ImageAttachment is image attachment type.
	ImageAttachment AttachmentType = "image"
	// AudioAttachment is audio attachment type.
	AudioAttachment AttachmentType = "audio"
	// VideoAttachment is video attachment type.
	VideoAttachment AttachmentType = "video"
	// FileAttachment is file attachment type.
	FileAttachment AttachmentType = "file"
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
	}

	return err
}

// Response is used for responding to events with messages.
type Response struct {
	token string
	to    Recipient
}

// Text sends a textual message.
func (r *Response) Text(message string) error {
	return r.TextWithReplies(message, nil)
}

// TextWithReplies sends a textual message with some replies
func (r *Response) TextWithReplies(message string, replies []QuickReply) error {
	m := SendMessage{
		Recipient: r.to,
		Message: MessageData{
			Text:         message,
			QuickReplies: replies,
		},
	}

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

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	return err
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
func (r *Response) Attachment(dataType AttachmentType, url string) error {
	m := SendStructuredMessage{
		Recipient: r.to,
		Message: StructuredMessageData{
			Attachment: StructuredMessageAttachment{
				Type: dataType,
				Payload: StructuredMessagePayload{
					Url: url,
				},
			},
		},
	}

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

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	return err
}

// AttachmentData sends an image, sound, video or a regular file to a chat via an io.Reader.
func (r *Response) AttachmentData(
	dataType AttachmentType,
	filename string,
	filedata io.Reader,
) error {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	data, err := w.CreateFormFile("filedata", filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(data, filedata)
	if err != nil {
		return err
	}

	w.WriteField("recipient", fmt.Sprintf(`{"id":"%v"}`, r.to.ID))
	w.WriteField("message", fmt.Sprintf(`{"attachment":{"type":"%v", "payload":{}}}`, dataType))

	req, err := http.NewRequest("POST", SendMessageURL, &b)
	if err != nil {
		return err
	}

	req.URL.RawQuery = "access_token=" + r.token

	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	var res bytes.Buffer
	res.ReadFrom(resp.Body)
	fmt.Println(res.String(), "DONE!")
	return nil
}

// ButtonTemplate sends a message with the main contents being button elements
func (r *Response) ButtonTemplate(text string, buttons *[]StructuredMessageButton) error {
	m := SendStructuredMessage{
		Recipient: r.to,
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
	}

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

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkFacebookError(resp.Body)
}

// GenericTemplate is a message which allows for structural elements to be sent
func (r *Response) GenericTemplate(elements *[]StructuredMessageElement) error {
	m := SendStructuredMessage{
		Recipient: r.to,
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
	}

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

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkFacebookError(resp.Body)
}

// SenderAction sends a info about sender action
func (r *Response) SenderAction(action string) error {
	m := SendSenderAction{
		Recipient:    r.to,
		SenderAction: action,
	}

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

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkFacebookError(resp.Body)
}

// SendMessage is the information sent in an API request to Facebook.
type SendMessage struct {
	Recipient Recipient   `json:"recipient"`
	Message   MessageData `json:"message"`
}

// MessageData is a text message with optional replies to be sent.
type MessageData struct {
	Text         string       `json:"text,omitempty"`
	QuickReplies []QuickReply `json:"quick_replies,omitempty"`
}

// SendStructuredMessage is a structured message template.
type SendStructuredMessage struct {
	Recipient Recipient             `json:"recipient"`
	Message   StructuredMessageData `json:"message"`
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
	TemplateType string                      `json:"template_type,omitempty"`
	Text         string                      `json:"text,omitempty"`
	Elements     *[]StructuredMessageElement `json:"elements,omitempty"`
	Buttons      *[]StructuredMessageButton  `json:"buttons,omitempty"`
	Url          string                      `json:"url,omitempty"`
}

// StructuredMessageElement is a response containing structural elements
type StructuredMessageElement struct {
	Title    string                    `json:"title"`
	ImageURL string                    `json:"image_url"`
	ItemURL  string                    `json:"item_url"`
	Subtitle string                    `json:"subtitle"`
	Buttons  []StructuredMessageButton `json:"buttons"`
}

// StructuredMessageButton is a response containing buttons
type StructuredMessageButton struct {
	Type    string `json:"type"`
	URL     string `json:"url,omitempty"`
	Title   string `json:"title"`
	Payload string `json:"payload,omitempty"`
}

// SendSenderAction is the information about sender action
type SendSenderAction struct {
	Recipient    Recipient `json:"recipient"`
	SenderAction string    `json:"sender_action"`
}
