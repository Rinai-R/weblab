package dao

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrDuplicated = errors.New("duplicated record")
)
