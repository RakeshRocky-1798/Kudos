package service

type CreateUserBulkResponse struct {
	SuccessCount int
	FailureCount int
	FailedRows   []map[string]string
}
