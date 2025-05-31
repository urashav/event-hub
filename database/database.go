package database

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type DB interface {
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
	Close() error
}
