package models

import "github.com/go-playground/validator/v10"

type User struct {
	ID       int    `json:"id,omitempty"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"` // Password should be hashed before storing
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)

}
