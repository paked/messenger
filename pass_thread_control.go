package messenger

type passThreadControl struct {
	Recipient   Recipient `json:"recipient"`
	TargetAppID int64     `json:"target_app_id"`
	Metadata    string    `json:"metadata"`
}
