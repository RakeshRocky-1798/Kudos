package kleosDb

type CreateKleosRequest struct {
	SenderID    string `json:"from_user,omitempty"`
	Message     string `json:"title,omitempty"`
	Achievement string `json:"level_give_kleos,omitempty"`
	ReceiverID  string `json:"member_slack_select,omitempty"`
	Year        string `json:"year,omitempty"`
	Month       string `json:"month,omitempty"`
	Week        string `json:"week,omitempty"`
	Day         string `json:"day,omitempty"`
}
