package messenger

import "time"

// Message represents a Facebook messenge message.
type Message struct {
	// Sender is who the message was sent from.
	Sender Sender `json:"-"`
	// Recipient is who the message was sent to.
	Recipient Recipient `json:"-"`
	// Time is when the message was sent.
	Time time.Time `json:"-"`
	// Mid is the ID of the message.
	Mid string `json:"mid"`
	// Seq is order the message was sent in relation to other messages.
	Seq int `json:"seq"`
	// Text is the textual contents of the message.
	Text string `json:"text"`
	// Attachments is the information about the attachments which were sent
	// with the message.
	Attachments []Attachment `json:"attachments"`
}

// Delivery represents a the event fired when a recipient reads one of Messengers sent
// messages.
type Delivery struct {
	// Mids are the IDs of the messages which were read.
	Mids []string `json:"mids"`
	// RawWatermark is the timestamp contained in the message of when the read was.
	RawWatermark int64 `json:"watermark"`
	// Seq is the sequence the message was sent in.
	Seq int `json:"seq"`
}

// Watermark is the RawWatermark timestamp rendered as a time.Time.
func (d Delivery) Watermark() time.Time {
	return time.Unix(d.RawWatermark, 0)
}
