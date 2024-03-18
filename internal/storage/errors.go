package storage

import (
	"errors"
	"github.com/lib/pq"
)

var (
	ErrNoRecord error = errors.New("storage: no matching record found")

	ErrInvalidData error = errors.New("storage: invalid data")

	ErrDuplicateUserName error = errors.New("storage: duplicate username")

	UniqueViolationErr = pq.ErrorCode("23505")
)

func IsErrorCode(err error, errcode pq.ErrorCode) bool {
	if pgerr, ok := err.(*pq.Error); ok {
		return pgerr.Code == errcode
	}
	return false
}
