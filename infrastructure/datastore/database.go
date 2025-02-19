package datastore

import (
	"context"
	"database/sql"
	"net/url"
	"time"

	_ "github.com/lib/pq"
	"github.com/syahidfrd/go-boilerplate/domain"
)

type Database struct {
	db *sql.DB
}

// NewDatabase will create new database instance
func NewDatabase(databaseURL string) (domain.Database, error) {
	parseDBUrl, _ := url.Parse(databaseURL)
	db, err := sql.Open(parseDBUrl.Scheme, databaseURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(5)

	return &Database{db: db}, nil
}

func (p *Database) BeginTx(ctx context.Context) (domain.Transaction, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	return &Transaction{tx: tx}, nil
}

func (p *Database) Close() error {
	return p.db.Close()
}

// Transaction implements Transaction interface
type Transaction struct {
	tx *sql.Tx
}

func (t *Transaction) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *Transaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

func (t *Transaction) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *Transaction) Commit() error {
	return t.tx.Commit()
}

func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}
