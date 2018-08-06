package messenger

// Profile is the public information of a Facebook user
type Profile struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ProfilePicURL string `json:"profile_pic"`
}
