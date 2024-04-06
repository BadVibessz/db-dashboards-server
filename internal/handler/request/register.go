package request

import "github.com/go-playground/validator/v10"

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

func (rr *RegisterRequest) Validate(valid *validator.Validate) error {
	return valid.Struct(rr)
}
