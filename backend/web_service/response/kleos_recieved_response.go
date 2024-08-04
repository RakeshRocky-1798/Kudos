package service

type KleosReceivedResponse struct {
	UserId        string `json:"userId"`
	UserEmail     string `json:"userEmail"`
	GivenCount    string `json:"givenCount"`
	ReceivedCount string `json:"receivedCount"`
}
