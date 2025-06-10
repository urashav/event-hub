package models

import "github.com/go-playground/validator/v10"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID       int    `json:"id,omitempty"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"` // Password should be hashed before storing
	Role     Role   `json:"role" validate:"required,oneof=user admin"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)

}
