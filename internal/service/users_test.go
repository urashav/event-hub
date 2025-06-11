package service

import (
	"context"
	"errors"
	"github.com/urashav/event-hub/internal/models"
	"github.com/urashav/event-hub/pkg/auth"
	"github.com/urashav/event-hub/pkg/hasher"
	"testing"
)

// MockUserRepository реализует интерфейс UserRepository для тестов
type mockUserRepository struct {
	users map[string]*models.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) (int, error) {
	if _, exists := m.users[user.Email]; exists {
		return 0, errors.New("user already exists")
	}
	user.ID = len(m.users) + 1
	m.users[user.Email] = user
	return user.ID, nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockUserRepository) UpdateRole(ctx context.Context, userID int, role models.Role) error {
	for _, user := range m.users {
		if user.ID == userID {
			user.Role = role
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *mockUserRepository) List(ctx context.Context) ([]*models.User, error) {
	users := make([]*models.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func TestUsersService_CreateUser(t *testing.T) {
	ctx := context.Background()
	repo := newMockUserRepository()
	hasher := hasher.NewHasher()
	tokenManager := auth.NewTokenManager("test_key")
	service := NewUserService(repo, hasher, tokenManager)

	tests := []struct {
		name    string
		user    models.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: models.User{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "duplicate user",
			user: models.User{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdUser, err := service.CreateUser(ctx, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && createdUser.Role != models.RoleUser {
				t.Errorf("CreateUser() role = %v, want %v", createdUser.Role, models.RoleUser)
			}
		})
	}
}

func TestUsersService_AuthenticateUser(t *testing.T) {
	ctx := context.Background()
	repo := newMockUserRepository()
	hasher := hasher.NewHasher()
	tokenManager := auth.NewTokenManager("test_key")
	service := NewUserService(repo, hasher, tokenManager)

	// Создаем тестового пользователя
	password := "password123"
	hashedPassword, _ := hasher.Hash(password)
	testUser := &models.User{
		Email:    "test@example.com",
		Password: hashedPassword,
		Role:     models.RoleUser,
	}
	repo.users[testUser.Email] = testUser

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "valid credentials",
			email:    "test@example.com",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "invalid password",
			email:    "test@example.com",
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "user not found",
			email:    "nonexistent@example.com",
			password: "password123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.AuthenticateUser(ctx, tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthenticateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Error("AuthenticateUser() token is empty")
			}
		})
	}
}
