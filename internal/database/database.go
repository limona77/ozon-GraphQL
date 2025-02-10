package database

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type Database interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
}

type Row interface {
	Scan(dest ...interface{}) error
}
