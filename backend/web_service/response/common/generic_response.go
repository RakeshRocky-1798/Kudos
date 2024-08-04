package service

type Error struct {
	Message  string      `json:"message,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
}

type Response struct {
	Data   interface{} `json:"data"`
	Error  Error       `json:"error,omitempty"`
	Status int         `json:"status,omitempty"`
}

type Page struct {
	TotalElements int   `json:"totalElements"`
	TotalPages    int64 `json:"totalPages"`
	PageSize      int64 `json:"pageSize"`
	PageNumber    int64 `json:"pageNumber"`
	HasData       bool  `json:"hasData"`
}

type PaginatedResponse struct {
	Data   interface{} `json:"data"`
	Page   Page        `json:"pages,omitempty"`
	Status int         `json:"status,omitempty"`
}
