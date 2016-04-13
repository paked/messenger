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
	Sender    Sender           `json:"sender"`
	Recipient Recipient        `json:"recipient"`
	Timestamp int64            `json:"timestamp"`
	Message   *MessageCallback `json:"message"`
}

type Sender struct {
	ID int64 `json:"id"`
}

type Recipient struct {
	ID int64 `json:"id"`
}

type MessageCallback struct {
	Mid  string `json:"mid"`
	Seq  int    `json:"seq"`
	Text string `json:"text"`
}
