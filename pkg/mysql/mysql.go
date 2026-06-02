// Package mysql implements MySQL connection.
package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
)

const (
	_defaultMaxPoolSize  = 10
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// MySQL -.
type MySQL struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *Pool
}

// New opens a MySQL connection pool.
func New(dsn string, opts ...Option) (*MySQL, error) {
	m := &MySQL{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(m)
	}

	m.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

	var (
		db  *sql.DB
		err error
	)

	for m.connAttempts > 0 {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.PingContext(context.Background())
		}

		if err == nil {
			break
		}

		log.Printf("MySQL is trying to connect, attempts left: %d", m.connAttempts)
		time.Sleep(m.connTimeout)
		m.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("mysql - New - connAttempts exhausted: %w", err)
	}

	db.SetMaxOpenConns(m.maxPoolSize)
	db.SetMaxIdleConns(m.maxPoolSize)

	m.Pool = newPool(db)

	return m, nil
}

// Close -.
func (m *MySQL) Close() {
	if m.Pool != nil && m.Pool.db != nil {
		_ = m.Pool.db.Close()
	}
}

// LastInsertID returns the auto-increment id after an INSERT.
func LastInsertID(ctx context.Context, pool *Pool, query string, args ...any) (int64, error) {
	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}
