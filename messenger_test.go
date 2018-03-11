package messenger

import (
	"testing"

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
			action := m.classify(test.msgInfo, Entry{})
			assert.Exactly(t, action, test.expected)
		})
	}
}
