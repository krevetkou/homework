package domain

import "errors"

var (
	ErrFieldsRequired = errors.New("all required fields must have values")
	ErrExists         = errors.New("already exists")
	ErrNotFound       = errors.New("not found")
)
