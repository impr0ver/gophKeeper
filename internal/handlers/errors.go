package handlers

import "errors"

// Handlers errors.
var (
	ErrEmptyField  = errors.New("field is empty")
	ErrWrongAESKey = errors.New("wrong AES key")
)
