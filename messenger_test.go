package messenger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessenger_Classify(t *testing.T) {
	m := New(Options{})

	for name, test := range map[string]struct {
		msgInfo  MessageInfo
		expected Action
	}{
		"unknown": {
			msgInfo:  MessageInfo{},
			expected: UnknownAction,
		},
		"message": {
			msgInfo: MessageInfo{
				Message: &Message{},
			},
			expected: TextAction,
		},
		"delivery": {
			msgInfo: MessageInfo{
				Delivery: &Delivery{},
			},
			expected: DeliveryAction,
		},
		"read": {
			msgInfo: MessageInfo{
				Read: &Read{},
			},
			expected: ReadAction,
		},
		"postback": {
			msgInfo: MessageInfo{
				PostBack: &PostBack{},
			},
			expected: PostBackAction,
		},
		"optin": {
			msgInfo: MessageInfo{
				OptIn: &OptIn{},
			},
			expected: OptInAction,
		},
		"referral": {
			msgInfo: MessageInfo{
				ReferralMessage: &ReferralMessage{},
			},
			expected: ReferralAction,
		},
	} {
		t.Run("action "+name, func(t *testing.T) {
			action := m.classify(test.msgInfo)
			assert.Exactly(t, action, test.expected)
		})
	}
}

func TestMessenger_Dispatch(t *testing.T) {
	type handlersCalls struct {
		message  int
		delivery int
		optin    int
		read     int
		postback int
		referral int
	}

	assertHandlersCalls := func(t *testing.T, actual *handlersCalls, expected handlersCalls) {
		assert.Equal(t, actual.message, expected.message)
		assert.Equal(t, actual.delivery, expected.delivery)
		assert.Equal(t, actual.optin, expected.optin)
		assert.Equal(t, actual.read, expected.read)
		assert.Equal(t, actual.postback, expected.postback)
		assert.Equal(t, actual.referral, expected.referral)
	}

	newReceive := func(msgInfo []MessageInfo) Receive {
		return Receive{
			Entry: []Entry{
				{
					Messaging: msgInfo,
				},
			},
		}
	}

	t.Run("message handlers", func(t *testing.T) {
		m := &Messenger{}
		h := &handlersCalls{}

		handler := func(msg Message, r *Response) {
			h.message++
			assert.NotNil(t, r)
			assert.EqualValues(t, 111, msg.Sender.ID)
			assert.EqualValues(t, 222, msg.Recipient.ID)
			assert.Equal(t, time.Unix(1543095111, 0), msg.Time)
		}

		messages := []MessageInfo{
			{
				Sender:    Sender{111},
				Recipient: Recipient{222},
				// 2018-11-24 21:31:51 UTC + 999ms
				Timestamp: 1543095111999,
				Message:   &Message{},
			},
		}

		// First handler
		m.HandleMessage(handler)

		m.dispatch(newReceive(messages))
		assertHandlersCalls(t, h, handlersCalls{message: 1})

		// Another handler
		m.HandleMessage(handler)

		m.dispatch(newReceive(messages))
		assertHandlersCalls(t, h, handlersCalls{message: 3})
	})

	t.Run("delivery handlers", func(t *testing.T) {
		m := &Messenger{}
		h := &handlersCalls{}

		handler := func(_ Delivery, r *Response) {
			h.delivery++
			assert.NotNil(t, r)
		}

		messages := []MessageInfo{
			{
				Sender:    Sender{111},
				Recipient: Recipient{222},
				// 2018-11-24 21:31:51 UTC + 999ms
				Timestamp: 1543095111999,
				Delivery:  &Delivery{},
			},
		}

		// First handler
		m.HandleDelivery(handler)

		m.dispatch(newReceive(messages))
		assertHandlersCalls(t, h, handlersCalls{delivery: 1})

		// Another handler
		m.HandleDelivery(handler)

		m.dispatch(newReceive(messages))
		assertHandlersCalls(t, h, handlersCalls{delivery: 3})
	})

	t.Run("read handlers", func(t *testing.T) {
		m := &Messenger{}
		h := &handlersCalls{}

		handler := func(_ Read, r *Response) {
			h.read++
			assert.NotNil(t, r)
		}

		messages := []MessageInfo{
			{
				Sender:    Sender{111},
				Recipient: Recipient{222},
				// 2018-11-24 21:31:51 UTC + 999ms
				Timestamp: 1543095111999,
				Read:      &Read{},
			},
		}

		// First handler
		m.HandleRead(handler)

		m.dispatch(newReceive(messages))
		assertHandlersCalls(t, h, handlersCalls{read: 1})

		// Another handler
		m.HandleRead(handler)

		m.dispatch(newReceive(messages))
		assertHandlersCalls(t, h, handlersCalls{read: 3})
	})

	t.Run("postback handlers", func(t *testing.T) {
		m := &Messenger{}
		h := &handlersCalls{}

		handler := func(msg PostBack, r *Response) {
			h.postback++
			assert.NotNil(t, r)
			assert.EqualValues(t, 111, msg.Sender.ID)
			assert.EqualValues(t, 222, msg.Recipient.ID)
			assert.Equal(t, time.Unix(1543095111, 0), msg.Time)
		}

		messages := []MessageInfo{
			{
				Sender:    Sender{111},
				Recipient: Recipient{222},
				// 2018-11-24 21:31:51 UTC + 999ms
				Timestamp: 1543095111999,
				PostBack:  &PostBack{},
			},
		}

		// First handler
		m.HandlePostBack(handler)

		m.dispatch(newReceive(messages))
		assertHandlersCalls(t, h, handlersCalls{postback: 1})

		// Another handler
		m.HandlePostBack(handler)

		m.dispatch(newReceive(messages))
		assertHandlersCalls(t, h, handlersCalls{postback: 3})
	})

	t.Run("optin handlers", func(t *testing.T) {
		m := &Messenger{}
		h := &handlersCalls{}

		handler := func(msg OptIn, r *Response) {
			h.optin++
			assert.NotNil(t, r)
			assert.EqualValues(t, 111, msg.Sender.ID)
			assert.EqualValues(t, 222, msg.Recipient.ID)
			assert.Equal(t, time.Unix(1543095111, 0), msg.Time)
		}

		messages := []MessageInfo{
			{
				Sender:    Sender{111},
				Recipient: Recipient{222},
				// 2018-11-24 21:31:51 UTC + 999ms
				Timestamp: 1543095111999,
				OptIn:     &OptIn{},
			},
		}

		// First handler
		m.HandleOptIn(handler)

		m.dispatch(newReceive(messages))
		assertHandlersCalls(t, h, handlersCalls{optin: 1})

		// Another handler
		m.HandleOptIn(handler)

		m.dispatch(newReceive(messages))
		assertHandlersCalls(t, h, handlersCalls{optin: 3})
	})

	t.Run("referral handlers", func(t *testing.T) {
		m := &Messenger{}
		h := &handlersCalls{}

		handler := func(msg ReferralMessage, r *Response) {
			h.referral++
			assert.NotNil(t, r)
			assert.EqualValues(t, 111, msg.Sender.ID)
			assert.EqualValues(t, 222, msg.Recipient.ID)
			assert.Equal(t, time.Unix(1543095111, 0), msg.Time)
		}

		messages := []MessageInfo{
			{
				Sender:    Sender{111},
				Recipient: Recipient{222},
				// 2018-11-24 21:31:51 UTC + 999ms
				Timestamp:       1543095111999,
				ReferralMessage: &ReferralMessage{},
			},
		}

		// First handler
		m.HandleReferral(handler)

		m.dispatch(newReceive(messages))
		assertHandlersCalls(t, h, handlersCalls{referral: 1})

		// Another handler
		m.HandleReferral(handler)

		m.dispatch(newReceive(messages))
		assertHandlersCalls(t, h, handlersCalls{referral: 3})
	})
}
