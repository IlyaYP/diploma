package pkg

import "errors"

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrInvalidLogin    = errors.New("invalid login")
	ErrInvalidPassword = errors.New("wrong password")
	ErrAlreadyExists   = errors.New("object exists in the DB")
	ErrNotExists       = errors.New("object not exists in the DB")
)
