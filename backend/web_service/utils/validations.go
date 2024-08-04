package utils

import (
	"errors"
	"strconv"
)

func ValidatePage(pageSize, pageNumber string) (int64, int64, error) {
	pageSizeInteger, err := strconv.ParseInt(pageSize, 0, 64)
	if err != nil {
		pageSizeInteger = 10
	}
	pageNumberInteger, err := strconv.ParseInt(pageNumber, 0, 64)
	if err != nil {
		pageNumberInteger = 0
	}

	if pageSizeInteger == 0 {
		return 0, 0, errors.New("page size cannot be zero")
	}

	return pageSizeInteger, pageNumberInteger, nil
}
