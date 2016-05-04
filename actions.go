package messenger

// Action is used to determine what kind of message a webhook event is.
type Action int

const (
	// UnknownAction means that the event was not able to be classified.
	UnknownAction Action = iota - 1
	// TextAction means that the event was a text message (May contain attachments).
	TextAction
	// DeliveryAction means that the event was a previous recipient reading their respective
	// messages.
	DeliveryAction
	// PostBackAction represents post call back
	PostBackAction
)
