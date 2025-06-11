package models

import (
	"testing"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "valid user",
			user: User{
				Email:    "test@example.com",
				Password: "password123",
				Role:     RoleUser,
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			user: User{
				Email:    "invalid-email",
				Password: "password123",
				Role:     RoleUser,
			},
			wantErr: true,
		},
		{
			name: "short password",
			user: User{
				Email:    "test@example.com",
				Password: "123",
				Role:     RoleUser,
			},
			wantErr: true,
		},
		{
			name: "invalid role",
			user: User{
				Email:    "test@example.com",
				Password: "password123",
				Role:     "invalid_role",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.user.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
