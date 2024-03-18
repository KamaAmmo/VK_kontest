package storage 

import "errors"

var (
	ErrNoRecord error = errors.New("Storage: no matching record found")

	ErrInvalidData error = errors.New("Film: invalid data")
)