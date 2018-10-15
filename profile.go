package messenger

// Profile is the public information of a Facebook user
type Profile struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	ProfilePicURL string  `json:"profile_pic"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Locale        string  `json:"locale"`
	Timezone      float64 `json:"timezone"`
	Gender        string  `json:"gender"`
}
