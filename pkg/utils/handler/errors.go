package handler

import "errors"

var (
	ErrNoHeaderProvided      = errors.New("no header provided")
	ErrInvalidHeaderProvided = errors.New("invalid header provided")

	ErrNoQueryParamProvided = errors.New("no query param provided")
)
