package comms

type CommunicationResponse struct {
	ReferenceID         string `json:"referenceId"`
	ExternalReferenceID string `json:"externalReferenceId,omitempty"`
	Message             string `json:"message"`
	Status              string `json:"status"`
	Medium              string `json:"medium"`
	TemplateID          string `json:"templateId"`
	RecipientID         string `json:"recipientId,omitempty"`
	Tenant              string `json:"tenant"`
	SentAt              string `json:"sentAt"`
	StatusDescription   string `json:"statusDescription,omitempty"`
	ReplyTopic          string `json:"replyTopic,omitempty"`
	VendorName          string `json:"vendorName"`
	VendorStatus        string `json:"vendorStatus,omitempty"`
}
