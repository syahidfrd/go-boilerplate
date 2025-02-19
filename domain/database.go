package domain

import (
	"context"
	"database/sql"
)

// Transaction interface (wraps DBTX)
type Transaction interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Commit() error
	Rollback() error
}

// Database interface
type Database interface {
	BeginTx(ctx context.Context) (Transaction, error)
	Close() error
}
