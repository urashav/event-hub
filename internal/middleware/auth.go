package middleware

import (
	"context"
	"github.com/urashav/event-hub/internal/models"
	"github.com/urashav/event-hub/pkg/auth"
	httputils "github.com/urashav/event-hub/pkg/httputilst"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	tokenManager *auth.TokenManager
}

func NewAuthMiddleware(tokenManager *auth.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: tokenManager,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httputils.ErrorResponse.Unauthorized(w, "Authorization header is required")
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			httputils.ErrorResponse.Unauthorized(w, "Invalid authorization header format")
			return
		}

		claims, err := m.tokenManager.ValidateToken(headerParts[1])
		if err != nil {
			httputils.ErrorResponse.Unauthorized(w, "Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "user_role", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole, ok := r.Context().Value("user_role").(string)
		if !ok {
			httputils.ErrorResponse.Forbidden(w, "Роль пользователя не определена")
			return
		}

		if userRole != string(models.RoleAdmin) {
			httputils.ErrorResponse.Forbidden(w, "Требуются права администратора")
			return
		}

		next.ServeHTTP(w, r)
	})
}
