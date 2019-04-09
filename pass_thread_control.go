package messenger

import "encoding/json"

func unmarshalPassThreadControl(data []byte) (passThreadControl, error) {
	var r passThreadControl
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *passThreadControl) marshal() ([]byte, error) {
	return json.Marshal(r)
}

type passThreadControl struct {
	Recipient   Recipient `json:"recipient"`
	TargetAppID int64     `json:"target_app_id"`
	Metadata    string    `json:"metadata"`
}
