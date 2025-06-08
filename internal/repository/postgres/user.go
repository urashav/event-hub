package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/urashav/event-hub/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u UserRepository) Create(ctx context.Context, user *models.User) (int, error) {
	var id int
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"

	err := u.db.QueryRowContext(ctx, query, user.Email, user.Password).Scan(&id)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return 0, errors.New("Email already exists")
		}
		return 0, err
	}

	return id, nil
}

func (u UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := "SELECT id, email FROM users WHERE id = $1"
	row := u.db.QueryRowContext(ctx, query, id)
	user := &models.User{}
	if err := row.Scan(&user.ID, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err // Other error
	}
	return user, nil
}

func (u UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, email, password FROM users WHERE email = $1"
	row := u.db.QueryRowContext(ctx, query, email)
	user := &models.User{}
	if err := row.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
