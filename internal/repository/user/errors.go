package user

import "errors"

var (
	ErrEmailExists = errors.New("user with this email already exists")
)
