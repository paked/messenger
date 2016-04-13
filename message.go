package messenger

import "time"

type Message struct {
	Sender    int64
	Recipient int64
	Time      time.Time
	Text      string
}
