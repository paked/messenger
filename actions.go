package messenger

type Action int

const (
	UnknownAction Action = iota - 1
	TextAction
	DeliveryAction
)
