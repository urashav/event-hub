package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/urashav/event-hub/internal/models"
	"github.com/urashav/event-hub/internal/repository/postgres"
	"github.com/urashav/event-hub/pkg/auth"
	"github.com/urashav/event-hub/pkg/hasher"
)

type UsersService struct {
	repo         *postgres.UserRepository
	hasher       *hasher.Hasher
	tokenManager *auth.TokenManager
}

func NewUserService(
	repo *postgres.UserRepository,
	hasher *hasher.Hasher,
	tokenManager *auth.TokenManager,
) *UsersService {
	// TODO: Отвязать от конкретного типа репозитория
	return &UsersService{
		repo:         repo,
		hasher:       hasher,
		tokenManager: tokenManager,
	}
}

func (s *UsersService) AuthenticateUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("Invalid user Email or Password")
		}
		return "", err
	}
	if !s.hasher.Verify(user.Password, password) {
		return "", errors.New("Invalid credentials")

	}
	token, err := s.tokenManager.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UsersService) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	hashedPassword, err := s.hasher.Hash(user.Password)
	if err != nil {
		return nil, fmt.Errorf("Failed to hash password: %v", err)
	}
	user.Password = hashedPassword
	user.Role = models.RoleUser

	id, err := s.repo.Create(ctx, &user)
	if err != nil {
		return nil, err
	}
	user.ID = int(id)
	return &user, nil
}

func (s *UsersService) UpdateUserRole(ctx context.Context, adminID int, userID int, newRole models.Role) error {
	admin, err := s.repo.GetByID(ctx, adminID)
	if err != nil {
		return err
	}

	if admin.Role != models.RoleAdmin {
		return errors.New("insufficient permissions")
	}

	return s.repo.UpdateRole(ctx, userID, newRole)
}

func (s *UsersService) ListUsers(ctx context.Context) ([]*models.User, error) {
	return s.repo.List(ctx)
}
