package messenger

// Receive is the format in which webhook events are sent.
type Receive struct {
	// Object should always be `page`. (I don't quite understand why)
	Object string `json:"object"`
	// Entry is all of the different messenger types which were
	// sent in this event.
	Entry []Entry `json:"entry"`
}

// Entry is a batch of events which were sent in this webhook trigger.
type Entry struct {
	// ID is the ID of the batch.
	ID int64 `json:"id,string"`
	// Time is when the batch was sent.
	Time int64 `json:"time"`
	// Messaging is the events that were sent in this Entry
	Messaging []MessageInfo `json:"messaging"`
}

// MessageInfo is an event that is fired by the webhook.
type MessageInfo struct {
	// Sender is who the event was sent from.
	Sender Sender `json:"sender"`
	// Recipient is who the event was sent to.
	Recipient Recipient `json:"recipient"`
	// Timestamp is the true time the event was triggered.
	Timestamp int64 `json:"timestamp"`
	// Message is the contents of a message if it is a MessageAction.
	// Nil if it is not a MessageAction.
	Message *Message `json:"message"`
	// Delivery is the contents of a message if it is a DeliveryAction.
	// Nil if it is not a DeliveryAction.
	Delivery *Delivery `json:"delivery"`

	PostBack *PostBack `json:"postback"`
}

// Sender is who the message was sent from.
type Sender struct {
	ID int64 `json:"id,string"`
}

// Recipient is who the message was sent to.
type Recipient struct {
	ID int64 `json:"id,string"`
}

// Attachment is a file which used in a message.
type Attachment struct {
	// Type is what type the message is. (image, video or audio)
	Type string `json:"type"`
	// Payload is the information for the file which was sent in the attachment.
	Payload Payload `json:"payload"`
}

// QuickReplie is a file which used in a message.
type QuickReplie struct {
	// ContentType is the type of replie
	ContentType string `json:"content_type"`
	// Title is the replie title
	Title string `json:"title"`
	// Payload is the  replie information
	Payload string `json:"payload"`
}

// Payload is the information on where an attachment is.
type Payload struct {
	// URL is where the attachment resides on the internet.
	URL string `json:"url,omitempty"`
}
