package domain

import "errors"

var (
	ErrFieldsRequired = errors.New("all required fields must have values")
	ErrMovieExists    = errors.New("movie already exists")
	ErrNotFound       = errors.New("movie not found")
)
