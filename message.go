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
	// Message is mine
	IsEcho bool `json:"is_echo,omitempty"`
	// Mid is the ID of the message.
	Mid string `json:"mid"`
	// Seq is order the message was sent in relation to other messages.
	Seq int `json:"seq"`
	// Text is the textual contents of the message.
	Text string `json:"text"`
	// Attachments is the information about the attachments which were sent
	// with the message.
	Attachments []Attachment `json:"attachments"`
	// Selected quick reply
	QuickReply *QuickReply `json:"quick_reply,omitempty"`
}

// Delivery represents a the event fired when Facebook delivers a message to the
// recipient.
type Delivery struct {
	// Mids are the IDs of the messages which were read.
	Mids []string `json:"mids"`
	// RawWatermark is the timestamp of when the delivery was.
	RawWatermark int64 `json:"watermark"`
	// Seq is the sequence the message was sent in.
	Seq int `json:"seq"`
}

// Read represents a the event fired when a message is read by the
// recipient.
type Read struct {
	// RawWatermark is the timestamp before which all messages have been read
	// by the user
	RawWatermark int64 `json:"watermark"`
	// Seq is the sequence the message was sent in.
	Seq int `json:"seq"`
}

// PostBack represents postback callback
type PostBack struct {
	// Sender is who the message was sent from.
	Sender Sender `json:"-"`
	// Recipient is who the message was sent to.
	Recipient Recipient `json:"-"`
	// Time is when the message was sent.
	Time time.Time `json:"-"`
	// PostBack ID
	Payload string `json:"payload"`
	// Optional referral info
	Referral Referral `json:"referral"`
}

// Watermark is the RawWatermark timestamp rendered as a time.Time.
func (d Delivery) Watermark() time.Time {
	return time.Unix(d.RawWatermark/int64(time.Microsecond), 0)
}

// Watermark is the RawWatermark timestamp rendered as a time.Time.
func (r Read) Watermark() time.Time {
	return time.Unix(r.RawWatermark/int64(time.Microsecond), 0)
}
