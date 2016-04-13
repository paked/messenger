package messenger

import "time"

type Message struct {
	Sender    Sender
	Recipient Recipient
	Time      time.Time
	Text      string
	Seq       int
}

type Delivery struct {
	Mids      []string
	Watermark time.Time
	Seq       int
}
