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
	query := "INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id"

	err := u.db.QueryRowContext(ctx, query, user.Email, user.Password, user.Role).Scan(&id)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return 0, errors.New("Email already exists")
		}
		return 0, err
	}

	return id, nil
}

func (u UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := "SELECT id, email, role FROM users WHERE id = $1"
	row := u.db.QueryRowContext(ctx, query, id)
	user := &models.User{}
	if err := row.Scan(&user.ID, &user.Email, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err // Other error
	}
	return user, nil
}

func (u UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, email, role FROM users WHERE email = $1"
	row := u.db.QueryRowContext(ctx, query, email)
	user := &models.User{}
	if err := row.Scan(&user.ID, &user.Email, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateRole(ctx context.Context, userID int, role models.Role) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET role = $1 WHERE id = $2",
		role, userID)
	return err
}

func (r *UserRepository) List(ctx context.Context) ([]*models.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, email, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
