package mysql

import (
	"context"
	"database/sql"
)

// Pool wraps database/sql with helper methods used by repositories.
type Pool struct {
	db *sql.DB
}

func newPool(db *sql.DB) *Pool {
	return &Pool{db: db}
}

// QueryRow -.
func (p *Pool) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return p.db.QueryRowContext(ctx, query, args...)
}

// Query -.
func (p *Pool) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return p.db.QueryContext(ctx, query, args...)
}

// Exec -.
func (p *Pool) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return p.db.ExecContext(ctx, query, args...)
}

// DB exposes the underlying connection for health checks.
func (p *Pool) DB() *sql.DB {
	return p.db
}

// Tx wraps sql.Tx with context-aware helpers used by repositories.
type Tx struct {
	tx *sql.Tx
}

// Begin starts a transaction.
func (p *Pool) Begin(ctx context.Context) (*Tx, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Tx{tx: tx}, nil
}

// Exec -.
func (t *Tx) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

// Commit -.
func (t *Tx) Commit(ctx context.Context) error {
	return t.tx.Commit()
}

// Rollback -.
func (t *Tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback()
}
