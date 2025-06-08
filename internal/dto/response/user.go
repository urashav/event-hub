package response

import "github.com/urashav/event-hub/internal/models"

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func FromModel(user *models.User) *UserResponse {
	return &UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
