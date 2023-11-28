package domain

import "errors"

var (
	ErrFieldsRequired = errors.New("all required fields must have values")
	ErrActorExists    = errors.New("actor already exists")
	ErrNotFound       = errors.New("actor not found")
)
