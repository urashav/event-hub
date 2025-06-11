package service

import (
	"context"
	"github.com/urashav/event-hub/internal/models"
)

type UserService interface {
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	AuthenticateUser(ctx context.Context, email, password string) (string, error)
	UpdateUserRole(ctx context.Context, adminID int, userID int, newRole models.Role) error
	ListUsers(ctx context.Context) ([]*models.User, error)
}
