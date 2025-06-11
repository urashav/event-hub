package middleware

import (
	"context"
	"github.com/urashav/event-hub/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/urashav/event-hub/pkg/auth"
)

func TestAuthMiddleware_RequireAuth(t *testing.T) {
	tokenManager := auth.NewTokenManager("test_key")
	middleware := NewAuthMiddleware(tokenManager)

	// Создаем тестовый токен
	token, _ := tokenManager.GenerateToken(1, "test@example.com", "user")

	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{
			name:       "valid token",
			token:      token,
			wantStatus: http.StatusOK,
		},
		{
			name:       "no token",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid token",
			token:      "invalid.token.here",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			w := httptest.NewRecorder()

			middleware.RequireAuth(nextHandler).ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("RequireAuth() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestAdminRequired(t *testing.T) {
	tests := []struct {
		name       string
		role       string
		wantStatus int
	}{
		{
			name:       "admin role",
			role:       string(models.RoleAdmin),
			wantStatus: http.StatusOK,
		},
		{
			name:       "user role",
			role:       string(models.RoleUser),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "no role",
			role:       "",
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.role != "" {
				ctx := context.WithValue(req.Context(), "user_role", tt.role)
				req = req.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			AdminRequired(nextHandler).ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("AdminRequired() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}
