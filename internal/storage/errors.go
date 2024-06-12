package storage

import "errors"

// Storage errors.
var (
	ErrLoginExists      = errors.New("login already exists")
	ErrWrongCredentials = errors.New("wrong login or password")
	ErrUnauthenticated  = errors.New("user is unauthorized")
	ErrNotFound         = errors.New("not found record with id")
	ErrUnknown          = errors.New("internal server error")
)
