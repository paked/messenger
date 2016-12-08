package messenger

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
