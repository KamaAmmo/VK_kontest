package storage 

import "errors"

var (
	ErrNoRecord error = errors.New("storage: no matching record found")
)