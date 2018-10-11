package messenger

import (
	"context"
)

type handlers struct {
	messageHandlers        []MessageHandler
	deliveryHandlers       []DeliveryHandler
	readHandlers           []ReadHandler
	postBackHandlers       []PostBackHandler
	optInHandlers          []OptInHandler
	referralHandlers       []ReferralHandler
	accountLinkingHandlers []AccountLinkingHandler
}

// MessageHandler is a handler used for responding to a message containing text.
type MessageHandler func(context.Context, Message, *Response)

// DeliveryHandler is a handler used for responding to a delivery receipt.
type DeliveryHandler func(context.Context, Delivery, *Response)

// ReadHandler is a handler used for responding to a read receipt.
type ReadHandler func(context.Context, Read, *Response)

// PostBackHandler is a handler used postback callbacks.
type PostBackHandler func(context.Context, PostBack, *Response)

// OptInHandler is a handler used to handle opt-ins.
type OptInHandler func(context.Context, OptIn, *Response)

// ReferralHandler is a handler used postback callbacks.
type ReferralHandler func(context.Context, ReferralMessage, *Response)

// AccountLinkingHandler is a handler used to react to an account
// being linked or unlinked.
type AccountLinkingHandler func(context.Context, AccountLinking, *Response)

// Handlers contain a registry of handlers matched off intent
var Handlers = handlers{}

// HandleMessage adds a new MessageHandler to the Messenger which will be triggered
// when a message is received by the client.
func (h *handlers) HandleMessage(f MessageHandler) {
	h.messageHandlers = append(h.messageHandlers, f)
}

// HandleDelivery adds a new DeliveryHandler to the Messenger which will be triggered
// when a previously sent message is delivered to the recipient.
func (h *handlers) HandleDelivery(f DeliveryHandler) {
	h.deliveryHandlers = append(h.deliveryHandlers, f)
}

// HandleOptIn adds a new OptInHandler to the Messenger which will be triggered
// once a user opts in to communicate with the bot.
func (h *handlers) HandleOptIn(f OptInHandler) {
	h.optInHandlers = append(h.optInHandlers, f)
}

// HandleRead adds a new DeliveryHandler to the Messenger which will be triggered
// when a previously sent message is read by the recipient.
func (h *handlers) HandleRead(f ReadHandler) {
	h.readHandlers = append(h.readHandlers, f)
}

// HandlePostBack adds a new PostBackHandler to the Messenger
func (h *handlers) HandlePostBack(f PostBackHandler) {
	h.postBackHandlers = append(h.postBackHandlers, f)
}

// HandleReferral adds a new ReferralHandler to the Messenger
func (h *handlers) HandleReferral(f ReferralHandler) {
	h.referralHandlers = append(h.referralHandlers, f)
}

// HandleAccountLinking adds a new AccountLinkingHandler to the Messenger
func (h *handlers) HandleAccountLinking(f AccountLinkingHandler) {
	h.accountLinkingHandlers = append(h.accountLinkingHandlers, f)
}
