package usersDb

type UserRequest struct {
	SlackUserId   string `json:"slack_user_id,omitempty"`
	UserName      string `json:"user_name,omitempty"`
	Email         string `json:"email,omitempty"`
	SlackImageUrl string `json:"slack_image_url,omitempty"`
	RealName      string `json:"real_name,omitempty"`
	GivenCount    int    `json:"given_count,omitempty"`
	ReceivedCount int    `json:"received_count,omitempty"`
}
