package service

import "strings"

func ErrorResponse(err error, code int, metadata interface{}) Response {
	return Response{
		Error: Error{
			Message:  err.Error(),
			Metadata: metadata,
		},
		Status: code,
	}
}

func SuccessResponse(data interface{}, code int) Response {
	return Response{
		Data:   data,
		Status: code,
	}
}

func SuccessPaginatedResponse(data interface{}, page Page, code int) PaginatedResponse {
	return PaginatedResponse{
		Data:   data,
		Page:   page,
		Status: code,
	}
}

func SplitUntilWord(input, stopWord string) (string, string) {

	lowercaseInput := strings.ToLower(input)
	lowercaseStopWord := strings.ToLower(stopWord)

	stopIndex := strings.Index(lowercaseInput, lowercaseStopWord)

	if stopIndex == -1 {
		// If stopWord is not found, return the entire input as the first part
		return input, ""
	}

	return strings.TrimSpace(input[:stopIndex]), strings.TrimSpace(input[stopIndex+len(stopWord):])
}
