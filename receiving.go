package messenger

type Receive struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID        int64         `json:"id"`
	Time      int64         `json:"time"`
	Messaging []MessageInfo `json:"messaging"`
}

type MessageInfo struct {
	Sender    Sender    `json:"sender"`
	Recipient Recipient `json:"recipient"`
	Timestamp int64     `json:"timestamp"`
	Message   *Message  `json:"message"`
	Delivery  *Delivery `json:"delivery"`
}

type Sender struct {
	ID int64 `json:"id"`
}

type Recipient struct {
	ID int64 `json:"id"`
}

type Attachment struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	URL string `json:"url"`
}
