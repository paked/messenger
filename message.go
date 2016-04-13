package messenger

import "time"

type Message struct {
	Sender    Sender
	Recipient Recipient
	Time      time.Time
	Text      string
	Seq       int
	Delivery  *Delivery
}

type Delivery struct {
	Mids      []string
	Watermark time.Time
	Seq       int
}
