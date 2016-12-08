package messenger

// Action is used to determine what kind of message a webhook event is.
type Action int

const (
	// UnknownAction means that the event was not able to be classified.
	UnknownAction Action = iota - 1
	// TextAction means that the event was a text message (May contain attachments).
	TextAction
	// DeliveryAction means that the event was advising of a successful delivery to a
	// previous recipient.
	DeliveryAction
	// ReadAction means that the event was a previous recipient reading their respective
	// messages.
	ReadAction
	// PostBackAction represents post call back
	PostBackAction
	// OptInAction represents opting in through the Send to Messenger button
	OptInAction
)
