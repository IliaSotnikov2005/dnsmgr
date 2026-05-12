package domain

import "errors"

var (
	ErrAlreadyExists = errors.New("dns entry already exists")
	ErrNotFound      = errors.New("dns entry not found")
	ErrInvalidIP     = errors.New("invalid IP address format")
)
