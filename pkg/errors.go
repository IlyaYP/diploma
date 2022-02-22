package pkg

import "errors"

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrInvalidPassword = errors.New("invalid password")
	ErrAlreadyExists   = errors.New("object exists in the DB")
	ErrNotExists       = errors.New("object not exists in the DB")
)
