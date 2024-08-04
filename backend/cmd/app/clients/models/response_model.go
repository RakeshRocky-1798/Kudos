package clients

type MjolnirError struct {
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Params  interface{} `json:"params,omitempty"`
}

type MjolnirSessionResponse struct {
	SessionToken      string         `json:"sessionToken,omitempty"`
	ClientId          string         `json:"clientId,omitempty"`
	EmailId           string         `json:"emailId,omitempty"`
	AccountId         string         `json:"accountId,omitempty"`
	PhoneNumber       string         `json:"phoneNumber,omitempty"`
	PreferredUsername string         `json:"preferred_username,omitempty"`
	Roles             []string       `json:"roles,omitempty"`
	Groups            []string       `json:"groups,omitempty"`
	Permissions       []string       `json:"permissions,omitempty"`
	StatusCode        int            `json:"statusCode,omitempty"`
	Errors            []MjolnirError `json:"errors"`
}
