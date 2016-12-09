package messenger

// Defines the different sizes available when setting up a CallToActionsItem
// of type "web_url". These values can be used in the "WebviewHeightRatio"
// field.
const (
	// WebviewCompact opens the page in a web view that takes half the screen
	// and covers only part of the conversation.
	WebviewCompact = "compact"

	// WebviewTall opens the page in a web view that covers about 75% of the
	// conversation.
	WebviewTall = "tall"

	// WebviewFull opens the page in a web view that completely covers the
	// conversation, and has a "back" button instead of a "close" one.
	WebviewFull = "full"
)

// GreetingSetting is the setting for greeting message
type GreetingSetting struct {
	SettingType string       `json:"setting_type"`
	Greeting    GreetingInfo `json:"greeting"`
}

// GreetingInfo contains greeting message
type GreetingInfo struct {
	Text string `json:"text"`
}

// CallToActionsSetting is the settings for Get Started and Persist Menu
type CallToActionsSetting struct {
	SettingType   string              `json:"setting_type"`
	ThreadState   string              `json:"thread_state"`
	CallToActions []CallToActionsItem `json:"call_to_actions"`
}

// CallToActionsItem contains Get Started button or item of Persist Menu
type CallToActionsItem struct {
	Type               string `json:"type,omitempty"`
	Title              string `json:"title,omitempty"`
	Payload            string `json:"payload,omitempty"`
	URL                string `json:"url,omitempty"`
	WebviewHeightRatio string `json:"webview_height_ratio,omitempty"`
	MessengerExtension bool   `json:"messenger_extensions,omitempty"`
}
