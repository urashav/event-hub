package postgres

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL driver
	"time"
)

const timeout = 10 * time.Second

func NewPostgresDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

//
//func (p *PostgresDB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
//	return p.db.QueryContext(ctx, query, args...)
//}
//
//func (p *PostgresDB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
//	return p.db.QueryRowContext(ctx, query, args...)
//}
//
//func (p *PostgresDB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
//	return p.db.ExecContext(ctx, query, args...)
//}
