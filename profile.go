package messenger

// Profile is the public information of a Facebook user
type Profile struct {
	Name          string  `json:"name"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	ProfilePicURL string  `json:"profile_pic"`
	Locale        string  `json:"locale"`
	Timezone      float64 `json:"timezone"`
	Gender        string  `json:"gender"`

	// instagram user profile
	Username  string `json:"username,omitempty"`
	IsPrivate bool   `json:"is_private,omitempty"`
	//FollowCount          int32  `json:"follow_count,omitempty"`
	FollowedByCount      int32 `json:"follower_count,omitempty"` // by the documentation followed_by_count
	IsVerifiedUser       bool  `json:"is_verified_user"`
	IsUserFollowBusiness bool  `json:"is_user_follow_business"`
	IsBusinessFollowUser bool  `json:"is_business_follow_user"`
}
