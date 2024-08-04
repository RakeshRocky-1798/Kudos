package request

type KleosRequest struct {
	SenderId    string   `json:"from"`
	ReceiverId  []string `json:"to"`
	Achievement string   `json:"achievement"`
	Message     string   `json:"message"`
	NeedSlack   bool     `json:"needSlack"`
}
