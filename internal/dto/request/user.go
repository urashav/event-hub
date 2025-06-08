package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/urashav/event-hub/internal/models"
)

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (s *SignUpRequest) Validate() error {
	return validate.Struct(s)
}

func (s *SignInRequest) Validate() error {
	return validate.Struct(s)
}

func (r *SignUpRequest) ToModel() *models.User {
	return &models.User{
		Email:    r.Email,
		Password: r.Password,
	}
}
