package userCountDb

type CreateUserCountRequest struct {
	UserId        string `json:"from_user,omitempty"`
	GivenCount    int    `json:"title,omitempty"`
	ReceivedCount int    `json:"level_give_kleos,omitempty"`
	Month         string `json:"month,omitempty"`
	Week          string `json:"week,omitempty"`
}
