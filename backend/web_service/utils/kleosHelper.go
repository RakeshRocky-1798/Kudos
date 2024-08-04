package utils

import (
	"strings"
)

func ExtractUsername(email string) string {

	parts := strings.Split(email, "@")

	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}

func ValidateHeaders(headers []string, mandatoryColumns []string) bool {
	if len(headers) < len(mandatoryColumns) {
		return false
	}

	lowerHeaders := make([]string, len(headers))

	for i, header := range headers {
		lowerHeaders[i] = strings.ToLower(header)
	}

	for i, col := range mandatoryColumns {
		col = strings.Trim(col, " ")
		if i >= len(lowerHeaders) || lowerHeaders[i] != strings.ToLower(col) {
			return false
		}
	}
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func FindIndex(headers []string, column string) int {
	for i, header := range headers {
		if strings.EqualFold(header, column) {
			return i
		}
	}
	return -1
}
