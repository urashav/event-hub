package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/urashav/event-hub/internal/dto/request"
	"github.com/urashav/event-hub/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserService struct {
	users map[string]*models.User
}

func newMockUserService() *mockUserService {
	return &mockUserService{
		users: make(map[string]*models.User),
	}
}

func (m *mockUserService) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	if _, exists := m.users[user.Email]; exists {
		return nil, errors.New("user already exists")
	}
	user.ID = len(m.users) + 1
	m.users[user.Email] = &user
	return &user, nil
}

func (m *mockUserService) AuthenticateUser(ctx context.Context, email, password string) (string, error) {
	_, exists := m.users[email]
	if !exists {
		return "", errors.New("user not found")
	}
	// В тестах просто возвращаем фиктивный токен
	return "test_token", nil
}

func (m *mockUserService) UpdateUserRole(ctx context.Context, adminID int, userID int, newRole models.Role) error {
	// Реализация для тестов
	return nil
}

func (m *mockUserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	users := make([]*models.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func TestUserHandler_SignUp(t *testing.T) {
	tests := []struct {
		name        string
		requestBody request.SignUpRequest
		wantStatus  int
	}{
		{
			name: "valid signup",
			requestBody: request.SignUpRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid email",
			requestBody: request.SignUpRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "short password",
			requestBody: request.SignUpRequest{
				Email:    "test@example.com",
				Password: "123",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newMockUserService()
			handler := NewUserHandler(service)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/user/signup", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.SignUp(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("SignUp() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}
