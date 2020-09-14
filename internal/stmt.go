package internal

import (
	"context"
	"database/sql/driver"
)

type Stmt struct {
	driver.Stmt
}

func (s *Stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	return s.exec(ctx, args)
}

func (s *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	return s.query(ctx, args)
}

func (s *Stmt) exec(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if s, ok := s.Stmt.(driver.StmtExecContext); ok {
		return s.ExecContext(ctx, args)
	}

	values := make([]driver.Value, len(args))

	for _, arg := range args {
		values[arg.Ordinal-1] = arg.Value
	}

	return s.Exec(values)
}

func (s *Stmt) query(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if s, ok := s.Stmt.(driver.StmtQueryContext); ok {
		return s.QueryContext(ctx, args)
	}

	values := make([]driver.Value, len(args))

	for _, arg := range args {
		values[arg.Ordinal-1] = arg.Value
	}

	return s.Query(values)
}