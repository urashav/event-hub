package repository

import (
	"context"
	"github.com/urashav/event-hub/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	UpdateRole(ctx context.Context, userID int, role models.Role) error
	List(ctx context.Context) ([]*models.User, error)
}
