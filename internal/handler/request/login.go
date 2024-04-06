package request

import "github.com/go-playground/validator/v10"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
}

func (lr *LoginRequest) Validate(valid *validator.Validate) error {
	return valid.Struct(lr)
}
