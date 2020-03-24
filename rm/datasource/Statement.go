package datasource

import (
	"context"
	"database/sql"
)

type Stmt struct {
	stmt *sql.Stmt
}

func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	// todo

	return s.stmt.ExecContext(ctx, args)
}

func (s *Stmt) Exec(args ...interface{}) (sql.Result, error) {
	return s.ExecContext(context.Background(), args...)
}

func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
}

func (s *Stmt) Query(args ...interface{}) (*sql.Rows, error) {
	return s.QueryContext(context.Background(), args...)
}

func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	rows, err := s.QueryContext(ctx, args...)
	if err != nil {
		return &Row{err: err}
	}
	return &Row{rows: rows}
}

func (s *Stmt) QueryRow(args ...interface{}) *Row {
	return s.QueryRowContext(context.Background(), args...)
}

func (s *Stmt) Close() error {

}
