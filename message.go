package messenger

import "time"

type Message struct {
	Sender      Sender       `json:"-"`
	Recipient   Recipient    `json:"-"`
	Time        time.Time    `json:"-"`
	Mid         string       `json:"mid"`
	Text        string       `json:"text"`
	Seq         int          `json:"seq"`
	Attachments []Attachment `json:"attachments"`
}

type Delivery struct {
	Mids         []string `json:"mids"`
	RawWatermark int64    `json:"watermark"`
	Seq          int      `json:"seq"`
}

func (d Delivery) Watermark() time.Time {
	return time.Unix(d.RawWatermark, 0)
}
