package storage

import "errors"

var (
	ErrCarNotFound = errors.New("car not found")
	ErrCarExists = errors.New("car already exists")
	ErrOwnerExists = errors.New("owner already exists")
)