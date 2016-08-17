package messenger

// Profile is the public information of a Facebook user
type Profile struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	ProfilePicURL string `json:"profile_pic"`
	Locale        string `json:"locale"`
	Timezone      int    `json:"timezone"`
	Gender        string `json:"gender"`
}
